-- +migrate Up
CREATE TABLE team_member (
    id TEXT PRIMARY KEY,
    team_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    relation_type TEXT NOT NULL DEFAULT 'direct',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES team(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    UNIQUE (team_id, user_id)
);

-- +migrate Down
DROP TABLE team_member;
