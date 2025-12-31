package emails

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/kiart-tantasi/crm-go/internal/contacts"
	"github.com/kiart-tantasi/crm-go/internal/httpclient"
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

	// Render and send email
	for _, contact := range contacts {
		_, err := RenderWithContact(email.Template, contact)
		if err != nil {
			fmt.Printf("failed to render email for %s: %v\n", contact.Email, err)
			continue
		}
		fmt.Printf("TODO: implement function to send email to smtp server (%s)\n", contact.Email)
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
