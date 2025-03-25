-- Res: User
-- Table: users

-- GetAll
SELECT id, username, email, slug FROM users;

-- Get
SELECT DISTINCT
    u.id, u.username, u.email, u.slug,
    r.id AS role_id, r.name AS role_name,
    p.id AS permission_id, p.name AS permission_name
FROM users u
         LEFT JOIN user_roles ur ON u.id = ur.user_id
         LEFT JOIN roles r ON ur.role_id = r.id
         LEFT JOIN role_permissions rp ON r.id = rp.role_id
         LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE u.id = ?;

-- Create
INSERT INTO users (id, username, email, slug) VALUES (?, ?, ?, ?);

-- Update
UPDATE users SET username = ?, email = ?, slug = ? WHERE id = ?;

-- Delete
DELETE FROM users WHERE id = ?;
