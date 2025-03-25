-- +migrate Up
CREATE TABLE user_permissions (
                                  user_id TEXT,
                                  permission_id TEXT,
                                  PRIMARY KEY (user_id, permission_id),
                                  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                                  FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_permissions_user_id ON user_permissions(user_id);
CREATE INDEX idx_user_permissions_permission_id ON user_permissions(permission_id);

-- +migrate Down
DROP TABLE user_permissions;
