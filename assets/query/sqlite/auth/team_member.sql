-- Res: TeamMember
-- AddUserToTeam
INSERT INTO team_member (id, team_id, user_id, relation_type, created_at, updated_at)
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- RemoveUserFromTeam
DELETE FROM team_member WHERE team_id = ? AND user_id = ?;

-- ListTeamMembers
SELECT u.*
FROM user u
JOIN team_member tm ON u.id = tm.user_id
WHERE tm.team_id = ?;

-- ListUsersNotInTeam
SELECT u.*
FROM user u
WHERE u.id NOT IN (
    SELECT user_id FROM team_member WHERE team_id = ?
);

-- ListUserTeams
SELECT t.*
FROM team t
JOIN team_member tm ON t.id = tm.team_id
WHERE tm.user_id = ?;

-- ListTeamsUserNotMember
SELECT t.*
FROM team t
WHERE t.id NOT IN (
    SELECT team_id FROM team_member WHERE user_id = ?
);

-- ListTeamMemberRoles
SELECT r.id, r.name, r.description, r.short_id, r.contextual, r.status, r.created_by, r.updated_by, r.created_at, r.updated_at
FROM role r
JOIN user_role ur ON r.id = ur.role_id
WHERE ur.user_id = ?
  AND ur.context_type = 'team'
  AND ur.context_id = ?
  AND r.contextual = true;

-- ListTeamUnassignedRoles
SELECT r.id, r.name, r.description, r.short_id, r.contextual, r.status, r.created_by, r.updated_by, r.created_at, r.updated_at
FROM role r
WHERE r.contextual = true
  AND r.id NOT IN (
    SELECT role_id
    FROM user_role
    WHERE user_id = ?
      AND context_type = 'team'
      AND context_id = ?
);
