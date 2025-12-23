package templates

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

func (s *Service) Upsert(ctx context.Context, t *Template) error {
	if t.Alias == "" {
		return fmt.Errorf("alias is required")
	}
	if t.Content == "" {
		return fmt.Errorf("template is required")
	}
	if t.AddedBy == 0 {
		return fmt.Errorf("added_by is required")
	}
	if t.ModifiedBy == 0 {
		return fmt.Errorf("modified_by is required")
	}
	// date_added and date_modified will be handled by the database (ON UPDATE statement)
	return s.repo.Upsert(ctx, t)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func Render(bodyTemplate string, data any) (string, error) {
	// Create template
	tmpl, err := template.New("").Parse(bodyTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with map-data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}
