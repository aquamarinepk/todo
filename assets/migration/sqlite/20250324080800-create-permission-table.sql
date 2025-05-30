-- +migrate Up
CREATE TABLE permission (
                             id TEXT PRIMARY KEY,
                             short_id TEXT,
                             name TEXT,
                             description TEXT,
                             created_by TEXT,
                             updated_by TEXT,
                             created_at TIMESTAMP,
                             updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE permission;
