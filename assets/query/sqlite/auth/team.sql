-- Res: Team
-- Table: teams

-- Create
INSERT INTO teams (id, org_id, slug, name, short_description, description, created_by, updated_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- GetAll
SELECT id, org_id, slug, name, short_description, description, created_by, updated_by, created_at, updated_at FROM teams WHERE org_id = ?;

-- Get
SELECT id, org_id, slug, name, short_description, description, created_by, updated_by, created_at, updated_at FROM teams WHERE id = ?;

-- Update
UPDATE teams SET slug = ?, name = ?, short_description = ?, description = ?, updated_by = ?, updated_at = ? WHERE id = ?;

-- Delete
DELETE FROM teams WHERE id = ?;
