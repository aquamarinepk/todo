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
	RolePermissions     []map[string]string `json:"role_permissions"`
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

// SeedAll loads and applies all auth seeds in a single transaction.
func (s *Seeder) SeedAll(ctx context.Context) error {
	byFeature, err := s.JSONSeeder.LoadJSONSeeds()
	if err != nil {
		return fmt.Errorf("failed to load JSON seeds: %w", err)
	}
	seeds, ok := byFeature["auth"]
	if !ok || len(seeds) == 0 {
		return fmt.Errorf("no auth seeds found")
	}
	for _, seed := range seeds {
		var data SeedData
		if err := json.Unmarshal([]byte(seed.Content), &data); err != nil {
			return fmt.Errorf("failed to unmarshal auth seed: %w", err)
		}
		if err := s.seedData(ctx, &data); err != nil {
			return err
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
func (s *Seeder) seedUsers(ctx context.Context, data *SeedData, userRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx for seedUsers: %w", err)
	}
	defer tx.Rollback()
	fmt.Println("=== [SEED] BEGIN seedUsers ===")
	defer fmt.Println("=== [SEED] END seedUsers ===")
	for i := range data.Users {
		u := &data.Users[i]
		u.GenCreationValues()
		err := s.repo.CreateUser(ctx, *u)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error inserting user: %w", err)
		}
		userRefMap[u.Ref()] = u.ID()
	}
	return tx.Commit()
}

func (s *Seeder) seedRoles(ctx context.Context, data *SeedData, roleRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx for seedRoles: %w", err)
	}
	defer tx.Rollback()
	fmt.Println("=== [SEED] BEGIN seedRoles ===")
	defer fmt.Println("=== [SEED] END seedRoles ===")
	for i := range data.Roles {
		r := &data.Roles[i]
		r.GenCreationValues()
		err := s.repo.CreateRole(ctx, *r)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error inserting role: %w", err)
		}
		roleRefMap[r.Ref()] = r.ID()
	}
	return tx.Commit()
}

func (s *Seeder) seedPermissions(ctx context.Context, data *SeedData, permRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx for seedPermissions: %w", err)
	}
	defer tx.Rollback()
	fmt.Println("=== [SEED] BEGIN seedPermissions ===")
	defer fmt.Println("=== [SEED] END seedPermissions ===")
	for i := range data.Permissions {
		p := &data.Permissions[i]
		p.GenCreationValues()
		err := s.repo.CreatePermission(ctx, *p)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error inserting permission: %w", err)
		}
		permRefMap[p.Ref()] = p.ID()
	}
	return tx.Commit()
}

func (s *Seeder) seedOrgs(ctx context.Context, data *SeedData, orgRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx for seedOrgs: %w", err)
	}
	defer tx.Rollback()
	fmt.Println("=== [SEED] BEGIN seedOrgs ===")
	defer fmt.Println("=== [SEED] END seedOrgs ===")
	for i := range data.Orgs {
		o := &data.Orgs[i]
		o.GenCreationValues()
		err := s.repo.CreateOrg(ctx, *o)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error inserting org: %w", err)
		}
		orgRefMap[o.Ref()] = o.ID()
	}
	b, _ := json.MarshalIndent(orgRefMap, "", "  ")
	fmt.Println("=== [SEED] SeedData state before commit ===")
	fmt.Println(string(b))
	return tx.Commit()
}

func (s *Seeder) seedTeams(ctx context.Context, data *SeedData, teamRefMap map[string]uuid.UUID, orgRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx for seedTeams: %w", err)
	}
	defer tx.Rollback()
	fmt.Println("=== [SEED] BEGIN seedTeams ===")
	defer fmt.Println("=== [SEED] END seedTeams ===")
	for i := range data.Teams {
		t := &data.Teams[i]
		orgRef := t.OrgRef
		orgID, ok := orgRefMap[orgRef]
		if !ok {
			return fmt.Errorf("###DEBUG###: error finding org ref for team: %s", orgRef)
		}
		t.OrgID = orgID
		t.GenCreationValues()
		err := s.repo.CreateTeam(ctx, *t)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error inserting team: %w", err)
		}
		teamRefMap[t.Ref()] = t.ID()
		b, _ := json.MarshalIndent(teamRefMap, "", "  ")
		fmt.Println("=== [SEED] SeedData state before commit ===")
		fmt.Println(string(b))
	}
	return tx.Commit()
}

func (s *Seeder) seedResources(ctx context.Context, data *SeedData, resourceRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx for seedResources: %w", err)
	}
	defer tx.Rollback()
	fmt.Println("=== [SEED] BEGIN seedResources ===")
	defer fmt.Println("=== [SEED] END seedResources ===")
	for i := range data.Resources {
		r := &data.Resources[i]
		r.GenCreationValues()
		err := s.repo.CreateResource(ctx, *r)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error inserting resource: %w", err)
		}
		resourceRefMap[r.Ref()] = r.ID()
	}
	return tx.Commit()
}

func (s *Seeder) seedUserRoles(ctx context.Context, data *SeedData, userRefMap, roleRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx for seedUserRoles: %w", err)
	}
	defer tx.Rollback()
	fmt.Println("=== [SEED] BEGIN seedUserRoles ===")
	defer fmt.Println("=== [SEED] END seedUserRoles ===")
	for _, ur := range data.UserRoles {
		userID, ok1 := userRefMap[ur["user_ref"]]
		roleID, ok2 := roleRefMap[ur["role_ref"]]
		if !ok1 || !ok2 {
			return fmt.Errorf("###DEBUG###: error finding user or role ref for user_role: %+v", ur)
		}
		err := s.repo.AddRole(ctx, userID, roleID)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error adding role to user: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Seeder) seedRolePermissions(ctx context.Context, data *SeedData, roleRefMap, permRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx for seedRolePermissions: %w", err)
	}
	defer tx.Rollback()
	fmt.Println("=== [SEED] BEGIN seedRolePermissions ===")
	defer fmt.Println("=== [SEED] END seedRolePermissions ===")
	for _, rp := range data.RolePermissions {
		roleID, ok1 := roleRefMap[rp["role_ref"]]
		permID, ok2 := permRefMap[rp["permission_ref"]]
		if !ok1 || !ok2 {
			return fmt.Errorf("###DEBUG###: error finding role or permission ref for role_permission: %+v", rp)
		}
		perm, err := s.repo.GetPermission(ctx, permID)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error getting permission: %w", err)
		}
		err = s.repo.AddPermissionToRole(ctx, roleID, perm)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error adding permission to role: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Seeder) seedUserPermissions(ctx context.Context, data *SeedData, userRefMap, permRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx for seedUserPermissions: %w", err)
	}
	defer tx.Rollback()
	fmt.Println("=== [SEED] BEGIN seedUserPermissions ===")
	defer fmt.Println("=== [SEED] END seedUserPermissions ===")
	for _, up := range data.UserPermissions {
		userID, ok1 := userRefMap[up["user_ref"]]
		permID, ok2 := permRefMap[up["permission_ref"]]
		if !ok1 || !ok2 {
			return fmt.Errorf("###DEBUG###: error finding user or permission ref for user_permission: %+v", up)
		}
		perm, err := s.repo.GetPermission(ctx, permID)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error getting permission: %w", err)
		}
		err = s.repo.AddPermissionToUser(ctx, userID, perm)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error adding permission to user: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Seeder) seedResourcePermissions(ctx context.Context, data *SeedData, resourceRefMap, permRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx for seedResourcePermissions: %w", err)
	}
	defer tx.Rollback()
	fmt.Println("=== [SEED] BEGIN seedResourcePermissions ===")
	defer fmt.Println("=== [SEED] END seedResourcePermissions ===")
	for _, rp := range data.ResourcePermissions {
		resourceID, ok1 := resourceRefMap[rp["resource_ref"]]
		permID, ok2 := permRefMap[rp["permission_ref"]]
		if !ok1 || !ok2 {
			return fmt.Errorf("###DEBUG###: error finding resource or permission ref for resource_permission: %+v", rp)
		}
		perm, err := s.repo.GetPermission(ctx, permID)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error getting permission: %w", err)
		}
		err = s.repo.AddPermissionToResource(ctx, resourceID, perm)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error adding permission to resource: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Seeder) seedOrgOwners(ctx context.Context, data *SeedData, orgRefMap, userRefMap map[string]uuid.UUID) error {
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx for seedOrgOwners: %w", err)
	}
	defer tx.Rollback()
	fmt.Println("=== [SEED] BEGIN seedOrgOwners ===")
	defer fmt.Println("=== [SEED] END seedOrgOwners ===")
	for _, oo := range data.OrgOwners {
		orgID, ok1 := orgRefMap[oo["org_ref"]]
		userID, ok2 := userRefMap[oo["user_ref"]]
		if !ok1 || !ok2 {
			return fmt.Errorf("###DEBUG###: error finding org or user ref for org_owner")
		}
		err := s.repo.AddOrgOwner(ctx, orgID, userID)
		if err != nil {
			return fmt.Errorf("###DEBUG###: error adding org owner: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Seeder) Start(ctx context.Context) error {
	return s.SeedAll(ctx)
}
