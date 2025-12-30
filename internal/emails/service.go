package emails

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

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

func Render(bodyEmail string) (string, error) {

	// FuncMap
	funcMap := template.FuncMap{
		// Function to fetch api data into map
		"fetch": func(url string) map[string]any {
			httpclient := httpclient.NewClient()
			data, _ := httpclient.FetchDataAndMap(url)
			// NOTE: ignore error for now
			// TODO: make user able to choose between making email fail on error or ignoring error
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
	if err := tmpl.Execute(&buf, nil); err != nil {
		return "", fmt.Errorf("failed to execute email: %w", err)
	}
	return buf.String(), nil
}
