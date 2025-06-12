-- Res: Org
-- Table: org

-- Create
INSERT INTO org (id, short_id, name, short_description, description, created_by, updated_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- GetDefault
SELECT id, short_id, name, short_description, description, created_by, updated_by, created_at, updated_at FROM org ORDER BY created_at ASC LIMIT 1;

-- Delete
DELETE FROM org WHERE id = ?;
