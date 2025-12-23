-- +goose Up
-- +goose StatementBegin
CREATE TABLE contact_list_emails (
    contact_list_id INT NOT NULL,
    email_id INT NOT NULL,
    PRIMARY KEY (contact_list_id, email_id),
    FOREIGN KEY (contact_list_id) REFERENCES contact_lists(id) ON DELETE CASCADE,
    FOREIGN KEY (email_id) REFERENCES emails(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contact_list_emails;
-- +goose StatementEnd
