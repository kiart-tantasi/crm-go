package emails

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

func (r *Repository) Upsert(ctx context.Context, e *Email) error {
	query := `INSERT INTO emails (id, alias, template, added_by, modified_by) 
	          VALUES (?, ?, ?, ?, ?) 
	          ON DUPLICATE KEY UPDATE 
	          alias = VALUES(alias), 
	          template = VALUES(template), 
	          modified_by = VALUES(modified_by)`
	_, err := r.db.ExecContext(ctx, query,
		e.ID, e.Alias, e.Template, e.AddedBy, e.ModifiedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert email: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*Email, error) {
	query := `SELECT id, alias, template, added_by, modified_by
	          FROM emails WHERE id = ?`
	var e Email
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.Alias, &e.Template, &e.AddedBy, &e.ModifiedBy,
	)
	// Not found
	if err == sql.ErrNoRows {
		return nil, nil
	}
	// Other errors
	if err != nil {
		return nil, fmt.Errorf("failed to get email by id: %w", err)
	}
	return &e, nil
}

func (r *Repository) List(ctx context.Context, limit int, offset int) ([]Email, error) {
	query := `SELECT id, alias, template, added_by, modified_by
	          FROM emails ORDER BY id ASC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list emails: %w", err)
	}
	defer rows.Close()

	var emails []Email
	for rows.Next() {
		var e Email
		err := rows.Scan(
			&e.ID, &e.Alias, &e.Template, &e.AddedBy, &e.ModifiedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan a sql row: %w", err)
		}
		emails = append(emails, e)
	}
	return emails, nil
}
