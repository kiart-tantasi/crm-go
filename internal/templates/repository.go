package templates

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

func (r *Repository) Upsert(ctx context.Context, t *Template) error {
	query := `INSERT INTO emails (id, alias, template, date_added, added_by, date_modified, modified_by) 
	          VALUES (?, ?, ?, ?, ?, ?, ?) 
	          ON DUPLICATE KEY UPDATE 
	          alias = VALUES(alias), 
	          template = VALUES(template), 
	          date_modified = VALUES(date_modified),
	          modified_by = VALUES(modified_by)`

	result, err := r.db.ExecContext(ctx, query,
		t.ID, t.Alias, t.Content, t.DateAdded, t.AddedBy, t.DateModified, t.ModifiedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert email: %w", err)
	}

	if t.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
		t.ID = int(id)
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

func (r *Repository) GetByID(ctx context.Context, id int) (*Template, error) {
	query := `SELECT id, alias, template, date_added, added_by, date_modified, modified_by 
	          FROM emails WHERE id = ?`
	var t Template
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Alias, &t.Content, &t.DateAdded, &t.AddedBy, &t.DateModified, &t.ModifiedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get email by id: %w", err)
	}
	return &t, nil
}
