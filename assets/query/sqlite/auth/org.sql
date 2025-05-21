-- Res: Org
-- Table: orgs

-- Create
INSERT INTO orgs (id, slug, name, short_description, description, created_by, updated_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- GetDefault
SELECT id, slug, name, short_description, description, created_by, updated_by, created_at, updated_at FROM orgs ORDER BY created_at ASC LIMIT 1;
