-- Res: ResourcePermission
-- Table: resource_permissions

-- AddPermissionToResource
INSERT INTO resource_permissions (resource_id, permission_id) VALUES (?, ?);

-- RemovePermissionFromResource
DELETE FROM resource_permissions WHERE resource_id = ? AND permission_id = ?;

-- GetResourcePermissions
SELECT p.id, p.name, p.description, p.slug, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permissions p
JOIN resource_permissions rp ON p.id = rp.permission_id
WHERE rp.resource_id = ?; 