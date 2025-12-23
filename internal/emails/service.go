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
	if e.Alias == "" {
		return fmt.Errorf("alias is required")
	}
	if e.Content == "" {
		return fmt.Errorf("email content is required")
	}
	if e.AddedBy == 0 {
		return fmt.Errorf("added_by is required")
	}
	if e.ModifiedBy == 0 {
		return fmt.Errorf("modified_by is required")
	}
	// date_added and date_modified will be handled by the database (ON UPDATE statement)
	return s.repo.Upsert(ctx, e)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
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
