-- +migrate Up
CREATE TABLE orgs (
    id TEXT PRIMARY KEY,
    slug TEXT,
    name TEXT,
    short_description TEXT,
    description TEXT,
    created_by TEXT,
    updated_by TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE orgs;
