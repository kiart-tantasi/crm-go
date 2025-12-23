-- +goose Up
-- +goose StatementBegin
CREATE TABLE contacts_emails (
    contact_id INT NOT NULL,
    email_id INT NOT NULL,
    PRIMARY KEY (contact_id, email_id),
    FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE,
    FOREIGN KEY (email_id) REFERENCES emails(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contacts_emails;
-- +goose StatementEnd
