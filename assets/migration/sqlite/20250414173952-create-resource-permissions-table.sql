-- +migrate Up
CREATE TABLE resource_permissions (
    resource_id TEXT NOT NULL,
    permission_id TEXT NOT NULL,
    created_by TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (resource_id, permission_id),
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE resource_permissions; 