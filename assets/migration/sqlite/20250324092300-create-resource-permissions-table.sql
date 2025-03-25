-- +migrate Up
CREATE TABLE resource_permissions (
                                      resource_id TEXT,
                                      permission_id TEXT,
                                      PRIMARY KEY (resource_id, permission_id),
                                      FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
                                      FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE INDEX idx_resource_permissions_resource_id ON resource_permissions(resource_id);
CREATE INDEX idx_resource_permissions_permission_id ON resource_permissions(permission_id);

-- +migrate Down
DROP TABLE resource_permissions;
