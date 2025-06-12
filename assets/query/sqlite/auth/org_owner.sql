-- Res: OrgOwner
-- Table: org_owner

-- Add
INSERT INTO org_owner (id, org_id, user_id) VALUES (?, ?, ?);

-- Remove
DELETE FROM org_owner WHERE org_id = ? AND user_id = ?;

-- GetOrgOwners
SELECT u.* FROM "user" u
INNER JOIN org_owner o ON u.id = o.user_id
WHERE o.org_id = ?;

-- GetOrgUnassignedOwners
SELECT u.* FROM "user" u
WHERE u.id NOT IN (
    SELECT user_id FROM org_owner WHERE org_id = ?
);
