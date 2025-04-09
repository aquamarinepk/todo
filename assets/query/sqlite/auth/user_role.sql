-- Res: UserRole
-- Table: user_role

-- AddRole
INSERT INTO user_roles (user_id, role_id)
SELECT ?, ?
WHERE EXISTS (
    SELECT 1 FROM roles WHERE id = ?
);

-- RemoveRole
DELETE FROM user_roles WHERE user_id = ? AND role_id = ?;

-- GetUserRoles
SELECT r.id, r.name, r.description, r.slug, r.status, r.created_by, r.updated_by, r.created_at, r.updated_at
FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = ?;

-- GetUserUnassignedRoles
SELECT r.id, r.name, r.description, r.slug, r.status, r.created_by, r.updated_by, r.created_at, r.updated_at
FROM roles r
WHERE r.id NOT IN (
    SELECT role_id FROM user_roles WHERE user_id = ?
);