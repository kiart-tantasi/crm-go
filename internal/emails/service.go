package emails

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strconv"
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

func (s *Service) GetByID(ctx context.Context, id string) (*Email, error) {
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert id to int: %w", err)
	}
	return s.repo.GetByID(ctx, idAsInt)
}

func (s *Service) List(ctx context.Context, limit string, offset string) ([]Email, error) {
	var err error
	// Limit
	limitAsInt := 100
	if limit != "" {
		if limitAsInt, err = strconv.Atoi(limit); err != nil {
			return nil, fmt.Errorf("failed to convert limit to int: %w", err)
		}

	}
	// Offset
	offsetAsInt := 0
	if offset != "" {
		if offsetAsInt, err = strconv.Atoi(offset); err != nil {
			return nil, fmt.Errorf("failed to convert offset to int: %w", err)
		}
	}
	return s.repo.List(ctx, limitAsInt, offsetAsInt)
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
