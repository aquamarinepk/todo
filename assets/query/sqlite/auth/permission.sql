-- Res: Permission
-- Table: permissions

-- GetAll
SELECT id, name, description, slug, created_by, updated_by, created_at, updated_at FROM permissions;

-- Get
SELECT id, name, description, slug, created_by, updated_by, created_at, updated_at
FROM permissions
WHERE id = ?;

-- Create
INSERT INTO permissions (id, name, description, slug, created_by, updated_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- Update
UPDATE permissions SET name = ?, description = ?, slug = ?, updated_by = ?, updated_at = ? WHERE id = ?;

-- Delete
DELETE FROM permissions WHERE id = ?;
