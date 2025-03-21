-- Res: Resource
-- Table: resources

-- GetAll
SELECT id, name, description, slug FROM resources;

-- Get
SELECT id, name, description, slug FROM resources WHERE id = ?;

-- Create
INSERT INTO resources (id, name, description, slug) VALUES (?, ?, ?, ?);

-- Update
UPDATE resources SET name = ?, description = ?, slug = ? WHERE id = ?;

-- Delete
DELETE FROM resources WHERE id = ?;

