-- +migrate Up
CREATE TABLE users (
                       id TEXT PRIMARY KEY,
                       name TEXT,
                       username TEXT,
                       email TEXT UNIQUE,
                       password TEXT,
                       slug TEXT,
                       created_by TEXT,
                       updated_by TEXT,
                       created_at TIMESTAMP,
                       updated_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- +migrate Down
DROP TABLE users;
