-- Res: UserPermissions
-- Table: user_permissions

-- GetUserAssignedPermissions
SELECT p.id, p.name, p.description, p.slug, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permissions p
WHERE p.id IN (
    SELECT rp.permission_id
    FROM role_permissions rp
             JOIN user_roles ur ON rp.role_id = ur.role_id
    WHERE ur.user_id = ?
    UNION
    SELECT up.permission_id
    FROM user_permissions up
    WHERE up.user_id = ?
);

-- GetUserIndirectPermissions
SELECT p.id, p.name, p.description, p.slug, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permissions p
WHERE p.id IN (
    SELECT rp.permission_id
    FROM role_permissions rp
             JOIN user_roles ur ON rp.role_id = ur.role_id
    WHERE ur.user_id = ?
);

-- GetUserDirectPermissions
SELECT p.id, p.name, p.description, p.slug, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permissions p
WHERE p.id IN (
    SELECT up.permission_id
    FROM user_permissions up
    WHERE up.user_id = ?
);

-- GetUserUnassignedPermissions
SELECT p.id, p.name, p.description, p.slug, p.created_by, p.updated_by, p.created_at, p.updated_at
FROM permissions p
WHERE p.id NOT IN (
    SELECT rp.permission_id
    FROM role_permissions rp
             JOIN user_roles ur ON rp.role_id = ur.role_id
    WHERE ur.user_id = ?
    UNION
    SELECT up.permission_id
    FROM user_permissions up
    WHERE up.user_id = ?
);

-- AddPermissionToUser
INSERT INTO user_permissions (user_id, permission_id) VALUES (?, ?);

-- RemovePermissionFromUser
DELETE FROM user_permissions WHERE user_id = ? AND permission_id = ?;
