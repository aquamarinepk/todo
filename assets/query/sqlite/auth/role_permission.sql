-- Res: RolePermission
-- Table: role_permissions

-- AddPermissionToRole
INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?);

-- RemovePermissionFromRole
DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?;

-- GetRolePermissions
SELECT p.id, p.slug, p.name, p.description, p.created_by, p.updated_by, p.created_at, p.updated_at  FROM permissions p
INNER JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = ?;

-- GetRoleUnassignedPermissions
SELECT p.id, p.slug, p.name, p.description, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permissions p
WHERE p.id NOT IN (
    SELECT rp.permission_id
    FROM role_permissions rp
    WHERE rp.role_id = ?
); 