package auth

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type Seeder struct {
	*am.JSONSeeder
	repo Repo
}

type SeedData struct {
	Users               []User              `json:"users"`
	Orgs                []Org               `json:"orgs"`
	Teams               []Team              `json:"teams"`
	Roles               []Role              `json:"roles"`
	Permissions         []Permission        `json:"permissions"`
	Resources           []Resource          `json:"resources"`
	UserRoles           []map[string]string `json:"user_roles"`
	RolePermissions     []map[string]string `json:"role_permission"`
	UserPermissions     []map[string]string `json:"user_permissions"`
	ResourcePermissions []map[string]string `json:"resource_permissions"`
	OrgOwners           []map[string]string `json:"org_owners"`
}

func NewSeeder(assetsFS embed.FS, engine string, repo Repo) *Seeder {
	return &Seeder{
		JSONSeeder: am.NewJSONSeeder(assetsFS, engine),
		repo:       repo,
	}
}

func (s *Seeder) Setup(ctx context.Context) error {
	if err := s.JSONSeeder.Setup(ctx); err != nil {
		return err
	}
	return s.SeedAll(ctx)
}

// SeedAll loads and applies all auth seeds in a single transaction.
func (s *Seeder) SeedAll(ctx context.Context) error {
	byFeature, err := s.JSONSeeder.LoadJSONSeeds()
	if err != nil {
		return fmt.Errorf("failed to load JSON seeds: %w", err)
	}
	for feature, seeds := range byFeature {
		for _, seed := range seeds {
			applied, err := s.JSONSeeder.SeedApplied(seed.Datetime, seed.Name, feature)
			if err != nil {
				return fmt.Errorf("failed to check if seed was applied: %w", err)
			}
			if applied {
				s.Log().Debugf("Seed already applied: %s-%s [%s]", seed.Datetime, seed.Name, feature)
				continue
			}

			var data SeedData
			err = json.Unmarshal([]byte(seed.Content), &data)
			if err != nil {
				return fmt.Errorf("failed to unmarshal %s seed: %w", feature, err)
			}

			err = s.seedData(ctx, &data)
			if err != nil {
				return err
			}

			err = s.JSONSeeder.ApplyJSONSeed(seed.Datetime, seed.Name, feature, seed.Content)
			if err != nil {
				s.Log().Errorf("error recording JSON seed: %v", err)
			}
		}
	}
	return nil
}

// seedData applies a single SeedData in a transaction.
func (s *Seeder) seedData(ctx context.Context, data *SeedData) error {
	userRefMap := make(map[string]uuid.UUID)
	roleRefMap := make(map[string]uuid.UUID)
	permRefMap := make(map[string]uuid.UUID)
	orgRefMap := make(map[string]uuid.UUID)
	teamRefMap := make(map[string]uuid.UUID)
	resourceRefMap := make(map[string]uuid.UUID)

	err := s.seedUsers(ctx, data, userRefMap)
	if err != nil {
		return err
	}
	err = s.seedOrgs(ctx, data, orgRefMap)
	if err != nil {
		return err
	}
	err = s.seedOrgOwners(ctx, data, orgRefMap, userRefMap)
	if err != nil {
		return err
	}
	err = s.seedTeams(ctx, data, teamRefMap, orgRefMap)
	if err != nil {
		return err
	}
	err = s.seedRoles(ctx, data, roleRefMap)
	if err != nil {
		return err
	}
	err = s.seedPermissions(ctx, data, permRefMap)
	if err != nil {
		return err
	}
	err = s.seedResources(ctx, data, resourceRefMap)
	if err != nil {
		return err
	}
	err = s.seedUserRoles(ctx, data, userRefMap, roleRefMap)
	if err != nil {
		return err
	}
	err = s.seedRolePermissions(ctx, data, roleRefMap, permRefMap)
	if err != nil {
		return err
	}
	err = s.seedUserPermissions(ctx, data, userRefMap, permRefMap)
	if err != nil {
		return err
	}
	err = s.seedResourcePermissions(ctx, data, resourceRefMap, permRefMap)
	if err != nil {
		return err
	}

	return nil
}

// --- Helper functions for each entity type ---
func (s *Seeder) withEncryptionKey(ctx context.Context) context.Context {
	key := s.Cfg().ByteSliceVal("sec.encryption.key")
	return context.WithValue(ctx, "encryptionKey", key)
}

func (s *Seeder) seedUsers(ctx context.Context, data *SeedData, userRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedUsers: %w", err)
	}
	defer tx.Rollback()
	s.Log().Debug("Seeding users: start")
	defer s.Log().Debug("Seeding users: end")
	for i := range data.Users {
		u := &data.Users[i]
		u.GenCreateValues()
		userCtx := s.withEncryptionKey(ctx)
		err := u.PrePersist(userCtx)
		if err != nil {
			return fmt.Errorf("error preparing user for insert: %w", err)
		}
		err = s.repo.CreateUser(ctx, *u)
		if err != nil {
			return fmt.Errorf("error inserting user: %w", err)
		}
		userRefMap[u.Ref()] = u.ID()
	}
	return tx.Commit()
}

func (s *Seeder) seedRoles(ctx context.Context, data *SeedData, roleRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedRoles: %w", err)
	}
	defer tx.Rollback()
	s.Log().Debug("Seeding roles: start")
	defer s.Log().Debug("Seeding roles: end")
	for i := range data.Roles {
		r := &data.Roles[i]
		r.GenCreateValues()
		err := s.repo.CreateRole(ctx, *r)
		if err != nil {
			return fmt.Errorf("error inserting role: %w", err)
		}
		roleRefMap[r.Ref()] = r.ID()
	}
	return tx.Commit()
}

func (s *Seeder) seedPermissions(ctx context.Context, data *SeedData, permRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedPermissions: %w", err)
	}
	defer tx.Rollback()
	s.Log().Debug("Seeding permissions: start")
	defer s.Log().Debug("Seeding permissions: end")
	for i := range data.Permissions {
		p := &data.Permissions[i]
		p.GenCreateValues()
		err := s.repo.CreatePermission(ctx, *p)
		if err != nil {
			return fmt.Errorf("error inserting permission: %w", err)
		}
		permRefMap[p.Ref()] = p.ID()
	}
	return tx.Commit()
}

func (s *Seeder) seedOrgs(ctx context.Context, data *SeedData, orgRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedOrgs: %w", err)
	}
	defer tx.Rollback()
	s.Log().Debug("Seeding orgs: start")
	defer s.Log().Debug("Seeding orgs: end")
	for i := range data.Orgs {
		o := &data.Orgs[i]
		o.GenCreateValues()
		err := s.repo.CreateOrg(ctx, *o)
		if err != nil {
			return fmt.Errorf("error inserting org: %w", err)
		}
		orgRefMap[o.Ref()] = o.ID()
	}
	b, _ := json.MarshalIndent(orgRefMap, "", "  ")
	s.Log().Debug("=== [SEED] SeedData state before commit ===")
	s.Log().Debug(string(b))
	return tx.Commit()
}

func (s *Seeder) seedTeams(ctx context.Context, data *SeedData, teamRefMap map[string]uuid.UUID, orgRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedTeams: %w", err)
	}
	defer tx.Rollback()
	s.Log().Debug("Seeding teams: start")
	defer s.Log().Debug("Seeding teams: end")
	for i := range data.Teams {
		t := &data.Teams[i]
		orgRef := t.OrgRef
		orgID, ok := orgRefMap[orgRef]
		if !ok {
			return fmt.Errorf("error finding org ref for team: %s", orgRef)
		}
		t.OrgID = orgID
		t.GenCreateValues()
		err := s.repo.CreateTeam(ctx, *t)
		if err != nil {
			return fmt.Errorf("error inserting team: %w", err)
		}
		teamRefMap[t.Ref()] = t.ID()
		b, _ := json.MarshalIndent(teamRefMap, "", "  ")
		s.Log().Debug("=== [SEED] SeedData state before commit ===")
		s.Log().Debug(string(b))
	}
	return tx.Commit()
}

func (s *Seeder) seedResources(ctx context.Context, data *SeedData, resourceRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedResources: %w", err)
	}
	defer tx.Rollback()
	s.Log().Debug("Seeding resources: start")
	defer s.Log().Debug("Seeding resources: end")
	for i := range data.Resources {
		r := &data.Resources[i]
		r.GenCreateValues()
		err := s.repo.CreateResource(ctx, *r)
		if err != nil {
			return fmt.Errorf("error inserting resource: %w", err)
		}
		resourceRefMap[r.Ref()] = r.ID()
	}
	return tx.Commit()
}

func (s *Seeder) seedUserRoles(ctx context.Context, data *SeedData, userRefMap, roleRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedUserRoles: %w", err)
	}
	defer tx.Rollback()
	s.Log().Debug("Seeding user roles: start")
	defer s.Log().Debug("Seeding user roles: end")
	for _, ur := range data.UserRoles {
		userID, ok1 := userRefMap[ur["user_ref"]]
		roleID, ok2 := roleRefMap[ur["role_ref"]]
		if !ok1 || !ok2 {
			return fmt.Errorf("error finding user or role ref for user_role: %+v", ur)
		}
		err := s.repo.AddRole(ctx, userID, roleID, "", "")
		if err != nil {
			return fmt.Errorf("error adding role to user: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Seeder) seedRolePermissions(ctx context.Context, data *SeedData, roleRefMap, permRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedRolePermissions: %w", err)
	}
	defer tx.Rollback()
	s.Log().Debug("Seeding role permissions: start")
	defer s.Log().Debug("Seeding role permissions: end")
	for _, rp := range data.RolePermissions {
		roleID, ok1 := roleRefMap[rp["role_ref"]]
		permID, ok2 := permRefMap[rp["permission_ref"]]
		if !ok1 || !ok2 {
			return fmt.Errorf("error finding role or permission ref for role_permission: %+v", rp)
		}
		perm, err := s.repo.GetPermission(ctx, permID)
		if err != nil {
			return fmt.Errorf("error getting permission: %w", err)
		}
		err = s.repo.AddPermissionToRole(ctx, roleID, perm)
		if err != nil {
			return fmt.Errorf("error adding permission to role: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Seeder) seedUserPermissions(ctx context.Context, data *SeedData, userRefMap, permRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedUserPermissions: %w", err)
	}
	defer tx.Rollback()
	s.Log().Debug("Seeding user permissions: start")
	defer s.Log().Debug("Seeding user permissions: end")
	for _, up := range data.UserPermissions {
		userID, ok1 := userRefMap[up["user_ref"]]
		permID, ok2 := permRefMap[up["permission_ref"]]
		if !ok1 || !ok2 {
			return fmt.Errorf("error finding user or permission ref for user_permission: %+v", up)
		}
		perm, err := s.repo.GetPermission(ctx, permID)
		if err != nil {
			return fmt.Errorf("error getting permission: %w", err)
		}
		err = s.repo.AddPermissionToUser(ctx, userID, perm)
		if err != nil {
			return fmt.Errorf("error adding permission to user: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Seeder) seedResourcePermissions(ctx context.Context, data *SeedData, resourceRefMap, permRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedResourcePermissions: %w", err)
	}
	defer tx.Rollback()
	s.Log().Debug("Seeding resource permissions: start")
	defer s.Log().Debug("Seeding resource permissions: end")
	for _, rp := range data.ResourcePermissions {
		resourceID, ok1 := resourceRefMap[rp["resource_ref"]]
		permID, ok2 := permRefMap[rp["permission_ref"]]
		if !ok1 || !ok2 {
			return fmt.Errorf("error finding resource or permission ref for resource_permission: %+v", rp)
		}
		perm, err := s.repo.GetPermission(ctx, permID)
		if err != nil {
			return fmt.Errorf("error getting permission: %w", err)
		}
		err = s.repo.AddPermissionToResource(ctx, resourceID, perm)
		if err != nil {
			return fmt.Errorf("error adding permission to resource: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Seeder) seedOrgOwners(ctx context.Context, data *SeedData, orgRefMap, userRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedOrgOwners: %w", err)
	}
	defer tx.Rollback()
	s.Log().Debug("Seeding org owners: start")
	defer s.Log().Debug("Seeding org owners: end")
	for _, oo := range data.OrgOwners {
		orgID, ok1 := orgRefMap[oo["org_ref"]]
		userID, ok2 := userRefMap[oo["user_ref"]]
		if !ok1 || !ok2 {
			return fmt.Errorf("error finding org or user ref for org_owner")
		}
		err := s.repo.AddOrgOwner(ctx, orgID, userID)
		if err != nil {
			return fmt.Errorf("error adding org owner: %w", err)
		}
	}
	return tx.Commit()
}
