package contactlists

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(ctx context.Context, cl *ContactList) error {
	query := `INSERT INTO contact_lists (name, added_by, modified_by)
	          VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		cl.Name, cl.AddedBy, cl.ModifiedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to insert contact list: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*ContactList, error) {
	query := `SELECT id, name, added_by, modified_by
	          FROM contact_lists WHERE id = ?`
	var cl ContactList
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cl.ID, &cl.Name, &cl.AddedBy, &cl.ModifiedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get contact list by id: %w", err)
	}
	return &cl, nil
}

func (r *Repository) List(ctx context.Context, limit int, offset int) ([]ContactList, error) {
	query := `SELECT id, name, added_by, modified_by
	          FROM contact_lists ORDER BY id ASC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list contact lists: %w", err)
	}
	defer rows.Close()

	var lists []ContactList
	for rows.Next() {
		var cl ContactList
		err := rows.Scan(
			&cl.ID, &cl.Name, &cl.AddedBy, &cl.ModifiedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan a contact list row: %w", err)
		}
		lists = append(lists, cl)
	}
	return lists, nil
}

func (r *Repository) AddContacts(ctx context.Context, contactListID int, contacts []AddContactItem) error {
	if len(contacts) == 0 {
		return nil
	}

	// Handle dynamic amount of contacts
	queryValues := make([]string, 0, len(contacts))
	queryArguments := make([]interface{}, 0, len(contacts)*3)
	for _, c := range contacts {
		queryValues = append(queryValues, "(?, ?, ?)")
		queryArguments = append(queryArguments, contactListID, c.ContactID, c.AddedBy)
	}
	query := fmt.Sprintf("INSERT INTO contact_list_members (contact_list_id, contact_id, added_by) VALUES %s ON DUPLICATE KEY UPDATE added_by = VALUES(added_by)",
		strings.Join(queryValues, ","))

	_, err := r.db.ExecContext(ctx, query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to add contact list contacts: %w", err)
	}
	return nil
}

func (r *Repository) RemoveContacts(ctx context.Context, contactListID int, contactIDs []int) error {
	if len(contactIDs) == 0 {
		return nil
	}

	// Handle dynamic amount of contact IDs
	queryPlaceholders := make([]string, len(contactIDs))
	queryArguments := make([]interface{}, len(contactIDs)+1)
	queryArguments[0] = contactListID
	for i, id := range contactIDs {
		queryPlaceholders[i] = "?"
		queryArguments[i+1] = id
	}

	// Execute
	query := fmt.Sprintf("DELETE FROM contact_list_members WHERE contact_list_id = ? AND contact_id IN (%s)",
		strings.Join(queryPlaceholders, ","))

	_, err := r.db.ExecContext(ctx, query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to remove contact list contacts: %w", err)
	}
	return nil
}
