-- +migrate Up
CREATE TABLE role (
                       id TEXT PRIMARY KEY,
                       short_id TEXT,
                       name TEXT,
                       description TEXT,
                       contextual BOOLEAN DEFAULT 0 NOT NULL,
                       status TEXT,
                       created_by TEXT,
                       updated_by TEXT,
                       created_at TIMESTAMP,
                       updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE role;
