-- Res: Permission
-- Table: permissions

-- GetAll
SELECT id, name, description, slug FROM permissions;

-- Get
SELECT id, name, description, slug FROM permissions WHERE id = ?;

-- Create
INSERT INTO permissions (id, name, description, slug) VALUES (?, ?, ?, ?);

-- Update
UPDATE permissions SET name = ?, description = ?, slug = ? WHERE id = ?;

-- Delete
DELETE FROM permissions WHERE id = ?;

