package emailsends

import "context"

type Service struct {
	repo *Repository
}

type EmailSendServiceI interface {
	Upsert(ctx context.Context, emailID int, contactID int) error
}

func NewService(repo *Repository) EmailSendServiceI {
	return &Service{repo: repo}
}

func (s *Service) Upsert(ctx context.Context, emailID int, contactID int) error {
	return s.repo.Upsert(ctx, emailID, contactID)
}
