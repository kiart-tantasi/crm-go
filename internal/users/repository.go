package users

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

func (r *Repository) Upsert(ctx context.Context, u *User) error {
	query := `INSERT INTO users (id, username, password, firstname, lastname, email, is_published, added_by, modified_by) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) 
	          ON DUPLICATE KEY UPDATE 
	          username = VALUES(username), 
	          password = VALUES(password), 
	          firstname = VALUES(firstname), 
	          lastname = VALUES(lastname), 
	          email = VALUES(email), 
	          is_published = VALUES(is_published), 
	          modified_by = VALUES(modified_by)`
	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.Username, u.Password, u.Firstname, u.Lastname, u.Email, u.IsPublished, u.AddedBy, u.ModifiedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*User, error) {
	query := `SELECT id, username, password, firstname, lastname, email, is_published, added_by, modified_by
	          FROM users WHERE id = ?`
	var u User
	var isPublished bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Password, &u.Firstname, &u.Lastname, &u.Email, &isPublished, &u.AddedBy, &u.ModifiedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	u.IsPublished = &isPublished
	return &u, nil
}

func (r *Repository) List(ctx context.Context, limit int, offset int) ([]User, error) {
	query := `SELECT id, username, password, firstname, lastname, email, is_published, added_by, modified_by
	          FROM users ORDER BY id ASC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		var isPublished bool
		err := rows.Scan(
			&u.ID, &u.Username, &u.Password, &u.Firstname, &u.Lastname, &u.Email, &isPublished, &u.AddedBy, &u.ModifiedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan a sql row: %w", err)
		}
		u.IsPublished = &isPublished
		users = append(users, u)
	}
	return users, nil
}
