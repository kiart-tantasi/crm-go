package emailsends

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(ctx context.Context, emailID int, contactID int, status string) error {
	query := `INSERT INTO email_sends (email_id, contact_id, status) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, emailID, contactID, status)
	if err != nil {
		return fmt.Errorf("failed to insert email send: %w", err)
	}
	return nil
}
