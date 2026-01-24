-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_sends (
    email_id INT NOT NULL,
    contact_id INT NOT NULL,
    date_sent DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (email_id, contact_id),
    FOREIGN KEY (email_id) REFERENCES emails(id),
    FOREIGN KEY (contact_id) REFERENCES contacts(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS email_sends;
-- +goose StatementEnd
