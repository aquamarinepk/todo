-- Res: List
-- Table: list

-- GetAll
SELECT id, name, short_id FROM list;

-- Get
SELECT id, name, short_id FROM list WHERE id = ?;

-- Create
INSERT INTO list (id, name, short_id) VALUES (?, ?, ?);

-- Update
UPDATE list SET name = ?, short_id = ? WHERE id = ?;

-- Delete
DELETE FROM list WHERE id = ?;
