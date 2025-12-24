package contactlists

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

func (s *Service) Insert(ctx context.Context, cl *ContactList) error {
	return s.repo.Insert(ctx, cl)
}

func (s *Service) GetByID(ctx context.Context, id string) (*ContactList, error) {
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert id to int: %w", err)
	}
	return s.repo.GetByID(ctx, idAsInt)
}

func (s *Service) List(ctx context.Context, limit string, offset string) ([]ContactList, error) {
	var err error
	limitAsInt := 100
	if limit != "" {
		if limitAsInt, err = strconv.Atoi(limit); err != nil {
			return nil, fmt.Errorf("failed to convert limit to int: %w", err)
		}
	}
	offsetAsInt := 0
	if offset != "" {
		if offsetAsInt, err = strconv.Atoi(offset); err != nil {
			return nil, fmt.Errorf("failed to convert offset to int: %w", err)
		}
	}
	return s.repo.List(ctx, limitAsInt, offsetAsInt)
}

func (s *Service) AddContacts(ctx context.Context, contactListID string, contacts []AddContactItem) error {
	clID, err := strconv.Atoi(contactListID)
	if err != nil {
		return fmt.Errorf("failed to convert contact_list_id to int: %w", err)
	}
	return s.repo.AddContacts(ctx, clID, contacts)
}

func (s *Service) RemoveContacts(ctx context.Context, contactListID string, contactIDs []int) error {
	clID, err := strconv.Atoi(contactListID)
	if err != nil {
		return fmt.Errorf("failed to convert contact_list_id to int: %w", err)
	}
	return s.repo.RemoveContacts(ctx, clID, contactIDs)
}
