-- +migrate Up
CREATE TABLE role_permission (
    role_id TEXT NOT NULL,
    permission_id TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permission(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE role_permission;
