-- Res: Permission
-- Table: permission

-- GetAll
SELECT id, short_id, name, description, created_by, updated_by, created_at, updated_at FROM permission;

-- Get
SELECT id, short_id, name, description, created_by, updated_by, created_at, updated_at
FROM permission
WHERE id = ?;

-- Create
INSERT INTO permission (id, short_id, name, description, created_by, updated_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- Update
UPDATE permission SET short_id = ?, name = ?, description = ?, updated_by = ?, updated_at = ? WHERE id = ?;

-- Delete
DELETE FROM permission WHERE id = ?;
