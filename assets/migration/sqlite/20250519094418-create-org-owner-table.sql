-- +migrate Up
CREATE TABLE org_owner (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (org_id) REFERENCES org(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    UNIQUE (org_id, user_id)
);

-- +migrate Down
DROP TABLE org_owner;
