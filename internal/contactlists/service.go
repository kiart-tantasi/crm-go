package contactlists

import (
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Upsert(ctx context.Context, cl *ContactList) error {
	return s.repo.Upsert(ctx, cl)
}

func (s *Service) GetByID(ctx context.Context, id int) (*ContactList, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, limit int, offset int) ([]ContactList, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) AddContacts(ctx context.Context, contactListID int, contactIDs []int, addedBy int) error {
	return s.repo.AddContacts(ctx, contactListID, contactIDs, addedBy)
}

func (s *Service) RemoveContacts(ctx context.Context, contactListID int, contactIDs []int) error {
	return s.repo.RemoveContacts(ctx, contactListID, contactIDs)
}
