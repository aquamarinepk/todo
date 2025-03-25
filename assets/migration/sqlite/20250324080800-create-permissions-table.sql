-- +migrate Up
CREATE TABLE permissions (
                             id TEXT PRIMARY KEY,
                             name TEXT,
                             description TEXT,
                             created_by TEXT,
                             updated_by TEXT,
                             created_at TIMESTAMP,
                             updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE permissions;
