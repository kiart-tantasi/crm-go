-- +goose Up
-- +goose StatementBegin
-- added_by and modified_by cannot use NOT NULL because the first user will not have it
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    firstname VARCHAR(255),
    lastname VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    is_published TINYINT(1) NOT NULL,
    date_added DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by INT,
    date_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
DROP TABLE IF EXISTS contacts;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
