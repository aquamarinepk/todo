-- Res: ResourcePermission
-- Table: resource_permission

-- AddPermissionToResource
INSERT INTO resource_permission (resource_id, permission_id) VALUES (?, ?);

-- RemovePermissionFromResource
DELETE FROM resource_permission WHERE resource_id = ? AND permission_id = ?;

-- GetResourcePermissions
SELECT p.id, p.name, p.description, p.short_id, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permission p
JOIN resource_permission rp ON p.id = rp.permission_id
WHERE rp.resource_id = ?;

-- GetResourceUnassignedPermissions
SELECT p.id, p.name, p.description, p.short_id, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permission p
WHERE p.id NOT IN (
    SELECT rp.permission_id
    FROM resource_permission rp
    WHERE rp.resource_id = ?
);