package emails

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
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

func Render(bodyEmail string, data any) (string, error) {
	// Create email
	tmpl, err := template.New("").Parse(bodyEmail)
	if err != nil {
		return "", fmt.Errorf("failed to parse email: %w", err)
	}

	// Execute email with map-data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute email: %w", err)
	}
	return buf.String(), nil
}
