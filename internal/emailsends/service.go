package emailsends

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Upsert(ctx context.Context, emailID int, contactID int) error {
	return s.repo.Upsert(ctx, emailID, contactID)
}
