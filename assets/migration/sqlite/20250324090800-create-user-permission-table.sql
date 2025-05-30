-- +migrate Up
CREATE TABLE user_permission (
    user_id TEXT NOT NULL,
    permission_id TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (user_id, permission_id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permission(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE user_permission;
