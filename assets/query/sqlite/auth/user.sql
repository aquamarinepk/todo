-- Res: User
-- Table: users

-- GetAll
SELECT id, username, email, slug, created_by, updated_by, created_at, updated_at FROM users;

-- Get
SELECT id, name, username, email, password, slug, created_by, updated_by, created_at, updated_at
FROM users
WHERE id = ?;

-- GetPreload
SELECT DISTINCT
    u.id, u.name, u.username, u.email, u.password, u.slug, u.created_by, u.updated_by, u.created_at, u.updated_at,
    r.id AS role_id, r.name AS role_name,
    p.id AS permission_id, p.name AS permission_name
FROM users u
       LEFT JOIN user_roles ur ON u.id = ur.user_id
       LEFT JOIN roles r ON ur.role_id = r.id
       LEFT JOIN role_permissions rp ON r.id = rp.role_id
       LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE u.id = ?;

-- Create
INSERT INTO users (id, username, email, name, password, slug, created_by, updated_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- Update
UPDATE users SET username = ?, email = ?, name = ?, password = ?, slug = ?, updated_by = ?, updated_at = ? WHERE id = ?;

-- Delete
DELETE FROM users WHERE id = ?;
