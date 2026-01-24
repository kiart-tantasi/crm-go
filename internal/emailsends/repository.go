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

func (r *Repository) Upsert(ctx context.Context, emailID int, contactID int) error {
	query := `INSERT INTO email_sends (email_id, contact_id) VALUES (?, ?)
	          ON DUPLICATE KEY UPDATE date_sent = VALUES(date_sent)`
	_, err := r.db.ExecContext(ctx, query, emailID, contactID)
	if err != nil {
		return fmt.Errorf("failed to upsert email send: %w", err)
	}
	return nil
}
