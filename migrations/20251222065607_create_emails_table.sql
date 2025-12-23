-- +goose Up
-- +goose StatementBegin
CREATE TABLE emails (
    id INT AUTO_INCREMENT PRIMARY KEY,
    alias VARCHAR(255) NOT NULL,
    template TEXT NOT NULL,
    date_added DATETIME NOT NULL,
    added_by INT NOT NULL,
    date_modified DATETIME ON UPDATE CURRENT_TIMESTAMP,
    modified_by INT,
    FOREIGN KEY (added_by) REFERENCES users(id),
    FOREIGN KEY (modified_by) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS emails;
-- +goose StatementEnd
