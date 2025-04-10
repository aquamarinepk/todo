-- Res: UserPermission
-- Table: user_permissions

-- AddPermissionToUser
INSERT INTO user_permissions (user_id, permission_id) VALUES (?, ?);

-- RemovePermissionFromUser
DELETE FROM user_permissions WHERE user_id = ? AND permission_id = ?; 