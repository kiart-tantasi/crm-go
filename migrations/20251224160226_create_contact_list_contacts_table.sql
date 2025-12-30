-- +goose Up
-- +goose StatementBegin
CREATE TABLE contact_list_contacts (
    contact_list_id INT NOT NULL,
    contact_id INT NOT NULL,
    date_added DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by INT NOT NULL,
    PRIMARY KEY (contact_list_id, contact_id),
    FOREIGN KEY (contact_list_id) REFERENCES contact_lists(id) ON DELETE CASCADE,
    FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contact_list_contacts;
-- +goose StatementEnd

