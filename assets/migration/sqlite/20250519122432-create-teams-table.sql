-- +migrate Up
CREATE TABLE teams (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL,
    slug TEXT NOT NULL,
    name TEXT NOT NULL,
    short_description TEXT,
    description TEXT,
    created_by TEXT,
    updated_by TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (org_id) REFERENCES orgs(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_teams_slug ON teams(slug);

-- +migrate Down
DROP TABLE teams;
