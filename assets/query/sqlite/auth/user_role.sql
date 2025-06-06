-- Res: UserRole
-- Table: user_role

-- AddRole
INSERT INTO user_role (user_id, role_id)
SELECT ?, ?
WHERE EXISTS (
    SELECT 1 FROM role WHERE id = ?
);

-- RemoveRole
DELETE FROM user_role WHERE user_id = ? AND role_id = ?;

-- GetUserAssignedRoles
SELECT r.id, r.name, r.description, r.short_id, r.status, r.created_by, r.updated_by, r.created_at, r.updated_at
FROM role r
JOIN user_role ur ON r.id = ur.role_id
WHERE ur.user_id = ?;

-- GetUserUnassignedRoles
SELECT r.id, r.name, r.description, r.short_id, r.status, r.created_by, r.updated_by, r.created_at, r.updated_at
FROM role r
WHERE r.id NOT IN (
    SELECT role_id FROM user_role WHERE user_id = ?
);