-- +migrate Up
CREATE TABLE user_role (
    user_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE user_role;
