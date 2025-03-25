-- +migrate Up
CREATE TABLE roles (
                       id TEXT PRIMARY KEY,
                       slug TEXT,
                       name TEXT,
                       description TEXT,
                       status TEXT,
                       created_by TEXT,
                       updated_by TEXT,
                       created_at TIMESTAMP,
                       updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE roles;
