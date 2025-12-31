package users

import (
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Upsert(ctx context.Context, u *User) error {
	// TODO: hash password
	return s.repo.Upsert(ctx, u)
}

func (s *Service) GetByID(ctx context.Context, id int) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, limit int, offset int) ([]User, error) {
	return s.repo.List(ctx, limit, offset)
}
