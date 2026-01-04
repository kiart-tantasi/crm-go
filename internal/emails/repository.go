package emails

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/kiart-tantasi/crm-go/internal/contacts"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(ctx context.Context, e *Email) error {
	query := `INSERT INTO emails (id, alias, subject, from_name, from_address, template, added_by, modified_by) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?) 
	          ON DUPLICATE KEY UPDATE 
	          alias = VALUES(alias), 
	          subject = VALUES(subject), 
	          from_name = VALUES(from_name), 
	          from_address = VALUES(from_address), 
	          template = VALUES(template), 
	          modified_by = VALUES(modified_by)`
	_, err := r.db.ExecContext(ctx, query,
		e.ID, e.Alias, e.Subject, e.FromName, e.FromAddress, e.Template, e.AddedBy, e.ModifiedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert email: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*Email, error) {
	query := `SELECT id, alias, subject, from_name, from_address, template, added_by, modified_by
	          FROM emails WHERE id = ?`
	var e Email
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.Alias, &e.Subject, &e.FromName, &e.FromAddress, &e.Template, &e.AddedBy, &e.ModifiedBy,
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
	query := `SELECT id, alias, subject, from_name, from_address, template, added_by, modified_by
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
			&e.ID, &e.Alias, &e.Subject, &e.FromName, &e.FromAddress, &e.Template, &e.AddedBy, &e.ModifiedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan a sql row: %w", err)
		}
		emails = append(emails, e)
	}
	return emails, nil
}
func (r *Repository) AddContactLists(ctx context.Context, emailID int, contactListIDs []int, addedBy int) error {
	if len(contactListIDs) == 0 {
		return nil
	}

	// Handle dynamic amount of contact lists
	queryValues := make([]string, 0, len(contactListIDs))
	queryArguments := make([]interface{}, 0, len(contactListIDs)*3)
	for _, contactListID := range contactListIDs {
		queryValues = append(queryValues, "(?, ?, ?)")
		queryArguments = append(queryArguments, emailID, contactListID, addedBy)
	}
	query := fmt.Sprintf("INSERT INTO email_contact_lists (email_id, contact_list_id, added_by) VALUES %s",
		strings.Join(queryValues, ","))

	_, err := r.db.ExecContext(ctx, query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to add email contact lists: %w", err)
	}
	return nil
}

func (r *Repository) RemoveContactLists(ctx context.Context, emailID int, contactListIDs []int) error {
	if len(contactListIDs) == 0 {
		return nil
	}

	// Handle dynamic amount of contact list IDs
	queryPlaceholders := make([]string, len(contactListIDs))
	queryArguments := make([]interface{}, len(contactListIDs)+1)
	queryArguments[0] = emailID
	for i, id := range contactListIDs {
		queryPlaceholders[i] = "?"
		queryArguments[i+1] = id
	}

	// Execute
	query := fmt.Sprintf("DELETE FROM email_contact_lists WHERE email_id = ? AND contact_list_id IN (%s)",
		strings.Join(queryPlaceholders, ","))

	_, err := r.db.ExecContext(ctx, query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to remove email contact lists: %w", err)
	}
	return nil
}

func (r *Repository) GetContactsByEmailID(ctx context.Context, emailID int) ([]contacts.Contact, error) {
	query := `
		SELECT c.id, c.firstname, c.lastname, c.email, c.is_published, c.added_by, c.modified_by
		FROM contacts c
		JOIN contact_list_contacts clc ON c.id = clc.contact_id
		JOIN email_contact_lists ecl ON clc.contact_list_id = ecl.contact_list_id
		WHERE ecl.email_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, emailID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts for email: %w", err)
	}
	defer rows.Close()

	var result []contacts.Contact
	for rows.Next() {
		var c contacts.Contact
		var isPublished bool
		err := rows.Scan(
			&c.ID, &c.Firstname, &c.Lastname, &c.Email, &isPublished, &c.AddedBy, &c.ModifiedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contact: %w", err)
		}
		c.IsPublished = &isPublished
		result = append(result, c)
	}
	return result, nil
}
