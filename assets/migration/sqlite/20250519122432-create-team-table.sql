-- +migrate Up
CREATE TABLE team (
    id TEXT PRIMARY KEY,
    short_id TEXT,
    org_id TEXT NOT NULL,
    name TEXT NOT NULL,
    short_description TEXT,
    description TEXT,
    created_by TEXT,
    updated_by TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (org_id) REFERENCES org(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE team;
