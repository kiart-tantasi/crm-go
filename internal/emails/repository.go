package emails

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(ctx context.Context, e *Email) error {
	query := `INSERT INTO emails (id, alias, template, date_added, added_by, date_modified, modified_by) 
	          VALUES (?, ?, ?, ?, ?, ?, ?) 
	          ON DUPLICATE KEY UPDATE 
	          alias = VALUES(alias), 
	          template = VALUES(template), 
	          date_modified = VALUES(date_modified),
	          modified_by = VALUES(modified_by)`
	_, err := r.db.ExecContext(ctx, query,
		e.ID, e.Alias, e.Content, e.DateAdded, e.AddedBy, e.DateModified, e.ModifiedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert email: %w", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM emails WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete email: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*Email, error) {
	query := `SELECT id, alias, template, date_added, added_by, date_modified, modified_by 
	          FROM emails WHERE id = ?`
	var e Email
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.Alias, &e.Content, &e.DateAdded, &e.AddedBy, &e.DateModified, &e.ModifiedBy,
	)
	if err == sql.ErrNoRows {
		// debug
		// TODO: remove
		log.Println("[DEBUG] There are no emails now....")
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get email by id: %w", err)
	}
	return &e, nil
}
