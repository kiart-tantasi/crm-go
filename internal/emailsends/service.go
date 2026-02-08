package emailsends

import "context"

type Service struct {
	repo *Repository
}

type EmailSendServiceI interface {
	Insert(ctx context.Context, emailID int, contactID int, status string) error
}

func NewService(repo *Repository) EmailSendServiceI {
	return &Service{repo: repo}
}

func (s *Service) Insert(ctx context.Context, emailID int, contactID int, status string) error {
	return s.repo.Insert(ctx, emailID, contactID, status)
}
