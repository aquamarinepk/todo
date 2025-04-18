-- +migrate Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT,
    username TEXT,
    email_enc BLOB,
    password_enc BLOB,
    slug TEXT,
    created_by TEXT,
    updated_by TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    last_login_at TIMESTAMP,
    last_login_ip TEXT,
    is_active BOOLEAN DEFAULT 1
);

CREATE INDEX idx_users_email_enc ON users(email_enc);

-- +migrate Down
DROP TABLE users;
