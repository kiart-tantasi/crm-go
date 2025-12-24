-- +goose Up
-- +goose StatementBegin
CREATE TABLE contact_lists (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date_added DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by INT NOT NULL,
    date_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    modified_by INT NOT NULL,
    FOREIGN KEY (added_by) REFERENCES users(id),
    FOREIGN KEY (modified_by) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contact_lists;
-- +goose StatementEnd
