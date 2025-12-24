package contacts

import (
	"context"
	"fmt"
	"strconv"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Upsert(ctx context.Context, c *Contact) error {
	return s.repo.Upsert(ctx, c)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Contact, error) {
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert id to int: %w", err)
	}
	return s.repo.GetByID(ctx, idAsInt)
}

func (s *Service) List(ctx context.Context, limit string, offset string) ([]Contact, error) {
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
