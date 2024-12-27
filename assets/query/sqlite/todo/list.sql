-- Resource: List
-- Table: lists

-- GetAll
SELECT id, name, slug FROM lists;

-- GetByID
SELECT id, name, slug FROM lists WHERE id = ?;

-- GetBySlug
SELECT id, name, slug FROM lists WHERE slug = ?;

-- Create
INSERT INTO lists (id, name, slug) VALUES (?, ?, ?);

-- Update
UPDATE lists SET name = ?, slug = ? WHERE id = ?;

-- Delete
DELETE FROM lists WHERE id = ?;