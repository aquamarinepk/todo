-- Res: List
-- Table: lists

-- GetAll
SELECT id, name, slug FROM lists;

-- Get
SELECT id, name, slug FROM lists WHERE id = ?;

-- Create
INSERT INTO lists (id, name, slug) VALUES (?, ?, ?);

-- Update
UPDATE lists SET name = ?, slug = ? WHERE id = ?;

-- Delete
DELETE FROM lists WHERE id = ?;
