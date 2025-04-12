-- Res: RolePermission
-- Table: role_permissions

-- AddPermissionToRole
INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?);

-- RemovePermissionFromRole
DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?;

-- GetRolePermissions
SELECT p.* FROM permissions p
INNER JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = ?; 