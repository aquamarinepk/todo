-- Res: UserPermissions
-- Table: user_permission

-- GetUserAssignedPermissions
SELECT p.id, p.name, p.description, p.short_id, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permission p
WHERE p.id IN (
    SELECT rp.permission_id
    FROM role_permission rp
             JOIN user_role ur ON rp.role_id = ur.role_id
    WHERE ur.user_id = ?
    UNION
    SELECT up.permission_id
    FROM user_permission up
    WHERE up.user_id = ?
);

-- GetUserIndirectPermissions
SELECT p.id, p.name, p.description, p.short_id, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permission p
WHERE p.id IN (
    SELECT rp.permission_id
    FROM role_permission rp
             JOIN user_role ur ON rp.role_id = ur.role_id
    WHERE ur.user_id = ?
);

-- GetUserDirectPermissions
SELECT p.id, p.name, p.description, p.short_id, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permission p
WHERE p.id IN (
    SELECT up.permission_id
    FROM user_permission up
    WHERE up.user_id = ?
);

-- GetUserUnassignedPermissions
SELECT p.id, p.name, p.description, p.short_id, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permission p
WHERE p.id NOT IN (
    SELECT rp.permission_id
    FROM role_permission rp
             JOIN user_role ur ON rp.role_id = ur.role_id
    WHERE ur.user_id = ?
    UNION
    SELECT up.permission_id
    FROM user_permission up
    WHERE up.user_id = ?
);

-- AddPermissionToUser
INSERT INTO user_permission (user_id, permission_id) VALUES (?, ?);

-- RemovePermissionFromUser
DELETE FROM user_permission WHERE user_id = ? AND permission_id = ?;
