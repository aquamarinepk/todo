-- Res: RolePermission
-- Table: role_permission

-- AddPermissionToRole
INSERT INTO role_permission (role_id, permission_id) VALUES (?, ?);

-- RemovePermissionFromRole
DELETE FROM role_permission WHERE role_id = ? AND permission_id = ?;

-- GetRolePermissions
SELECT p.id, p.short_id, p.name, p.description, p.created_by, p.updated_by, p.created_at, p.updated_at  FROM permission p
INNER JOIN role_permission rp ON p.id = rp.permission_id
WHERE rp.role_id = ?;

-- GetRoleUnassignedPermissions
SELECT p.id, p.short_id, p.name, p.description, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permission p
WHERE p.id NOT IN (
    SELECT rp.permission_id
    FROM role_permission rp
    WHERE rp.role_id = ?
);