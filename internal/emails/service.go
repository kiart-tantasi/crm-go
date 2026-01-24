package emails

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"

	"github.com/kiart-tantasi/crm-go/internal/contacts"
	"github.com/kiart-tantasi/crm-go/internal/emailsends"
	"github.com/kiart-tantasi/crm-go/internal/httpclient"
	"github.com/kiart-tantasi/crm-go/internal/smtppool"
)

type Service struct {
	repo             *Repository
	emailSendService *emailsends.Service
}

func NewService(repo *Repository, emailSendService *emailsends.Service) *Service {
	return &Service{
		repo:             repo,
		emailSendService: emailSendService,
	}
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

	// Init new smtp pool for each request
	pool, err := smtppool.New()
	if err != nil {
		return fmt.Errorf("failed to create smtp pool: %w", err)
	}

	// Render and send email
	taskId := rand.Intn(100_000)
	log.Printf("Sending email to %d contacts (task id: %d)", len(contacts), taskId)

	// ==================
	// TODO: remove after creating global config table
	fromName := email.FromName.String
	if !email.FromName.Valid {
		fromName = "Demo"
	}
	fromAddress := email.FromAddress.String
	if !email.FromAddress.Valid {
		fromAddress = "no-reply@petchblog.net"
	}
	// ==================

	for _, contact := range contacts {
		rendered, err := RenderWithContact(email.Template, contact)
		if err != nil {
			log.Printf("failed to render email for %s: %v\n", contact.Email, err)
			continue
		}
		header := fmt.Sprintf("From: %s <%s>\r\nTo: %s <%s>\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n", fromName, fromAddress, fmt.Sprintf("%s %s", contact.Firstname, contact.Lastname), contact.Email, email.Subject)

		// Send async
		pool.SendMail(fromAddress, []string{contact.Email}, []byte(fmt.Sprintf("%s%s", header, rendered)))

		// Record sent status immediately after queuing
		// Note: Since SendMail is async, we record it here to prevent duplicate queuing if Send is called again quickly.
		if err := s.emailSendService.Upsert(ctx, id, contact.ID); err != nil {
			log.Printf("failed to mark email as sent for contact %d: %v\n", contact.ID, err)
		}
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
