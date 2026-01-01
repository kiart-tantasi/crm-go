package emails

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"

	"github.com/kiart-tantasi/crm-go/internal/contacts"
	"github.com/kiart-tantasi/crm-go/internal/httpclient"
	"github.com/kiart-tantasi/crm-go/internal/mailing"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Upsert(ctx context.Context, e *Email) error {
	return s.repo.Upsert(ctx, e)
}

func (s *Service) GetByID(ctx context.Context, id int) (*Email, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, limit int, offset int) ([]Email, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) AddContactLists(ctx context.Context, emailID int, contactListIDs []int, addedBy int) error {
	return s.repo.AddContactLists(ctx, emailID, contactListIDs, addedBy)
}

func (s *Service) RemoveContactLists(ctx context.Context, emailID int, contactListIDs []int) error {
	return s.repo.RemoveContactLists(ctx, emailID, contactListIDs)
}

func (s *Service) Send(ctx context.Context, id int) error {
	// Fetch email
	email, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get email: %w", err)
	}
	if email == nil {
		return fmt.Errorf("email not found")
	}

	// Fetch contacts
	contacts, err := s.repo.GetContactsByEmailID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get contacts: %w", err)
	}

	mailer := mailing.NewMailer("localhost", 25, "", "")

	// Render and send email
	for _, contact := range contacts {
		rendered, err := RenderWithContact(email.Template, contact)
		if err != nil {
			fmt.Printf("failed to render email for %s: %v\n", contact.Email, err)
			continue
		}

		// DEBUG
		// TODO: remove hardcoded fields
		// TODO: add subject column to emails
		// TODO: add from_name and from_email column to emails
		// TODO: add from_name and from_email to global config table
		// TODO: create smtp pool to send emails
		// TODO: return response without waiting for emails to be sent
		params := mailing.EmailParams{
			FromName: "from",
			FromAddr: "from@test.com",
			ToName:   fmt.Sprintf("%s %s", contact.Firstname, contact.Lastname),
			ToAddr:   contact.Email,
			Subject:  "TEST SUBJECT",
			Body:     rendered,
		}
		err = mailer.Send(params)
		if err != nil {
			log.Fatalf("Error sending email: %v", err)
		}

		log.Printf("Rendered %s", rendered)
		log.Printf("TODO: implement function to send email to smtp server (%s)\n", contact.Email)
		// END OF DEBUG
	}
	return nil
}

func Render(bodyEmail string) (string, error) {
	return render(bodyEmail, nil)
}

func RenderWithContact(bodyEmail string, contact contacts.Contact) (string, error) {
	data := map[string]any{
		"contact": contact,
	}
	return render(bodyEmail, data)
}

func render(bodyEmail string, data any) (string, error) {
	// FuncMap
	funcMap := template.FuncMap{
		"fetch": func(url string) map[string]any {
			httpclient := httpclient.NewClient()
			// TODO: make user able to choose between making email fail on error or ignoring error
			data, _ := httpclient.FetchDataAndMap(url)
			return data
		},
	}

	// Create template
	tmpl, err := template.New("").Funcs(funcMap).Parse(bodyEmail)
	if err != nil {
		return "", fmt.Errorf("failed to parse email: %w", err)
	}

	// Execute and return as string
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute email: %w", err)
	}
	return buf.String(), nil
}
