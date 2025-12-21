-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    firstname VARCHAR(255),
    lastname VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    is_published TINYINT(1) NOT NULL,
    date_added DATETIME,
    added_by INT,
    date_modified DATETIME ON UPDATE CURRENT_TIMESTAMP,
    modified_by INT,
    FOREIGN KEY (added_by) REFERENCES users(id),
    FOREIGN KEY (modified_by) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE contacts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    firstname VARCHAR(255),
    lastname VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    is_published TINYINT(1) NOT NULL,
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
DROP TABLE IF EXISTS contacts;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
