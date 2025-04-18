-- Res: Permission
-- Table: permissions

-- GetAll
SELECT id, slug, name, description, created_by, updated_by, created_at, updated_at FROM permissions;

-- Get
SELECT id, slug, name, description, created_by, updated_by, created_at, updated_at
FROM permissions
WHERE id = ?;

-- Create
INSERT INTO permissions (id, slug, name, description, created_by, updated_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- Update
UPDATE permissions SET slug = ?, name = ?, description = ?, updated_by = ?, updated_at = ? WHERE id = ?;

-- Delete
DELETE FROM permissions WHERE id = ?;
