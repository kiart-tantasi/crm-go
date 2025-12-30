-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_contact_lists (
    email_id INT NOT NULL,
    contact_list_id INT NOT NULL,
    date_added DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by INT NOT NULL,
    PRIMARY KEY (email_id, contact_list_id),
    FOREIGN KEY (email_id) REFERENCES emails(id) ON DELETE CASCADE,
    FOREIGN KEY (contact_list_id) REFERENCES contact_lists(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS email_contact_lists;
-- +goose StatementEnd
