-- Res: Role
-- Table: role

-- GetAll
SELECT id, name, description, short_id FROM role;

-- Get
SELECT id, name, description, short_id, created_by, updated_by, created_at, updated_at
FROM role
WHERE id = ?;

-- GetPreload
SELECT DISTINCT
    r.id, r.name, r.description, r.short_id, r.created_by, r.updated_by, r.created_at, r.updated_at,
    p.id AS permission_id, p.name AS permission_name, p.short_id AS permission_short_id
FROM role r
    LEFT JOIN role_permission rp ON r.id = rp.role_id
    LEFT JOIN permission p ON rp.permission_id = p.id
WHERE r.id = ?;

-- Create
INSERT INTO role (id, name, description, short_id, created_by, updated_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- Update
UPDATE role SET name = ?, description = ?, short_id = ?, updated_by = ?, updated_at = ? WHERE id = ?;

-- Delete
DELETE FROM role WHERE id = ?;
