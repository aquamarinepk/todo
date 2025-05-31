-- +migrate Up
CREATE TABLE user_role (
    user_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    context_type TEXT,
    context_id   TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (user_id, role_id, context_type, context_id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_role_context ON user_role(context_type, context_id);

-- +migrate Down
DROP TABLE user_role;
