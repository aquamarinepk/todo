-- Res: Resource
-- Table: resources

-- GetAll
SELECT id, name, description, slug, created_by, updated_by, created_at, updated_at FROM resources;

-- Get
SELECT id, name, description, slug, created_by, updated_by, created_at, updated_at
FROM resources
WHERE id = ?;

-- GetPreload
SELECT DISTINCT
    r.id, r.name, r.description, r.slug, r.created_by, r.updated_by, r.created_at, r.updated_at,
    p.id AS permission_id, p.name AS permission_name, p.slug AS permission_slug
FROM resources r
    LEFT JOIN resource_permissions rp ON r.id = rp.resource_id
    LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE r.id = ?;

-- Create
INSERT INTO resources (id, name, description, slug, created_by, updated_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- Update
UPDATE resources SET name = ?, description = ?, slug = ?, updated_by = ?, updated_at = ? WHERE id = ?;

-- Delete
DELETE FROM resources WHERE id = ?;
