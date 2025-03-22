-- Res: Role
-- Table: roles

-- GetAll
SELECT id, name, description, slug FROM roles;

-- Get
SELECT id, name, description, slug FROM roles WHERE id = ?;

-- Create
INSERT INTO roles (id, name, description, slug) VALUES (?, ?, ?, ?);

-- Update
UPDATE roles SET name = ?, description = ?, slug = ? WHERE id = ?;

-- Delete
DELETE FROM roles WHERE id = ?;


