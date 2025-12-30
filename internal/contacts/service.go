package contacts

import (
	"context"
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

func (s *Service) GetByID(ctx context.Context, id int) (*Contact, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, limit int, offset int) ([]Contact, error) {
	return s.repo.List(ctx, limit, offset)
}
