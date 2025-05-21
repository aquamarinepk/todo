-- Res: Org
-- Table: org_owners

-- Add
INSERT INTO org_owners (id, org_id, user_id, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- GetOrgOwners
SELECT u.id, u.slug, u.username, u.email_enc, u.name, u.password_enc, u.created_by, u.updated_by, u.created_at, u.updated_at, u.last_login_at, u.last_login_ip, u.is_active
FROM users u
JOIN org_owners oo ON u.id = oo.user_id
WHERE oo.org_id = ?;

-- GetOrgUnassignedOwners
SELECT u.id, u.slug, u.username, u.email_enc, u.name, u.password_enc, u.created_by, u.updated_by, u.created_at, u.updated_at, u.last_login_at, u.last_login_ip, u.is_active
FROM users u
WHERE u.id NOT IN (SELECT user_id FROM org_owners WHERE org_id = ?);
