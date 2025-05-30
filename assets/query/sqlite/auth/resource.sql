-- Res: Resource
-- Table: resource

-- GetAll
SELECT id, name, description, short_id, created_by, updated_by, created_at, updated_at FROM resource;

-- Get
SELECT id, name, description, short_id, created_by, updated_by, created_at, updated_at
FROM resource
WHERE id = ?;

-- GetPreload
SELECT DISTINCT
    r.id, r.name, r.description, r.short_id, r.created_by, r.updated_by, r.created_at, r.updated_at,
    p.id AS permission_id, p.name AS permission_name, p.short_id AS permission_short_id
FROM resource r
    LEFT JOIN resource_permission rp ON r.id = rp.resource_id
    LEFT JOIN permission p ON rp.permission_id = p.id
WHERE r.id = ?;

-- Create
INSERT INTO resource (id, name, description, short_id, created_by, updated_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- Update
UPDATE resource SET name = ?, description = ?, short_id = ?, updated_by = ?, updated_at = ? WHERE id = ?;

-- Delete
DELETE FROM resource WHERE id = ?;
