-- Res: User
-- Table: users

-- GetAll
SELECT id, username, email_enc, password_enc, name, slug, created_by, updated_by, created_at, updated_at, last_login_at, last_login_ip, is_active FROM users;

-- Get
SELECT id, name, username, email_enc, password_enc, slug, created_by, updated_by, created_at, updated_at, last_login_at, last_login_ip, is_active
FROM users
WHERE id = ?;

-- GetPreload
SELECT DISTINCT
    u.id, u.name, u.username, u.email_enc, u.password_enc, u.slug, u.created_by, u.updated_by, u.created_at, u.updated_at, u.last_login_at, u.last_login_ip, u.is_active,
    r.id AS role_id, r.name AS role_name,
    p.id AS permission_id, p.name AS permission_name
FROM users u
       LEFT JOIN user_roles ur ON u.id = ur.user_id
       LEFT JOIN roles r ON ur.role_id = r.id
       LEFT JOIN role_permissions rp ON r.id = rp.role_id
       LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE u.id = ?;

-- Create
INSERT INTO users (id, username, email_enc, name, password_enc, slug, created_by, updated_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- Update
UPDATE users SET username = ?, email_enc = ?, name = ?, slug = ?, updated_by = ?, updated_at = ? WHERE id = ?;

-- Delete
DELETE FROM users WHERE id = ?;

-- UpdatePassword
UPDATE users
SET password_enc = ?, updated_by = ?, updated_at = ?
WHERE id = ?;
