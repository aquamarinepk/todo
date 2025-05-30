-- +migrate Up
CREATE TABLE resource_permission (
    resource_id TEXT NOT NULL,
    permission_id TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (resource_id, permission_id),
    FOREIGN KEY (resource_id) REFERENCES resource(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permission(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE resource_permission;