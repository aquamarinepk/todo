-- +migrate Up
CREATE TABLE user (
    id TEXT PRIMARY KEY,
    short_id TEXT,
    name TEXT,
    username TEXT,
    email_enc BLOB,
    password_enc BLOB,
    created_by TEXT,
    updated_by TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    last_login_at TIMESTAMP,
    last_login_ip TEXT,
    is_active BOOLEAN DEFAULT 1
);

CREATE INDEX idx_user_email_enc ON user(email_enc);

-- +migrate Down
DROP TABLE user;
