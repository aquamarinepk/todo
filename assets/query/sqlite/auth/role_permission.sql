-- Res: RolePermission
-- Table: role_permissions

-- AddPermissionToRole
INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?);

-- RemovePermissionFromRole
DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?; 