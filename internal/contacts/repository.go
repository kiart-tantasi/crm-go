package contacts

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

func (r *Repository) Upsert(ctx context.Context, c *Contact) error {
	query := `INSERT INTO contacts (id, firstname, lastname, email, is_published, added_by, modified_by) 
	          VALUES (?, ?, ?, ?, ?, ?, ?) 
	          ON DUPLICATE KEY UPDATE 
	          firstname = VALUES(firstname), 
	          lastname = VALUES(lastname), 
	          email = VALUES(email), 
	          is_published = VALUES(is_published), 
	          modified_by = VALUES(modified_by)`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.Firstname, c.Lastname, c.Email, c.IsPublished, c.AddedBy, c.ModifiedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert contact: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*Contact, error) {
	query := `SELECT id, firstname, lastname, email, is_published, added_by, modified_by
	          FROM contacts WHERE id = ?`
	var c Contact
	var isPublished bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Firstname, &c.Lastname, &c.Email, &isPublished, &c.AddedBy, &c.ModifiedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get contact by id: %w", err)
	}
	c.IsPublished = &isPublished
	return &c, nil
}

func (r *Repository) List(ctx context.Context, limit int, offset int) ([]Contact, error) {
	query := `SELECT id, firstname, lastname, email, is_published, added_by, modified_by
	          FROM contacts ORDER BY id ASC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list contacts: %w", err)
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var c Contact
		var isPublished bool
		err := rows.Scan(
			&c.ID, &c.Firstname, &c.Lastname, &c.Email, &isPublished, &c.AddedBy, &c.ModifiedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan a sql row: %w", err)
		}
		c.IsPublished = &isPublished
		contacts = append(contacts, c)
	}
	return contacts, nil
}
