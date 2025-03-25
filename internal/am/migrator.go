package am

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	EngSQLite     = "sqlite"
	EngPostgres   = "postgres"
	MigrationPath = "assets/migration/%s"
)

type Migrator struct {
	Core
	db         *sql.DB
	assetsFS   embed.FS
	engine     string
	migrations sync.Map
}

type FileMigration struct {
	Datetime string
	Name     string
	Up       string
	Down     string
}

type DBMigration struct {
	Datetime string
	Name     string
}

func NewMigrator(assetsFS embed.FS, engine string, opts ...Option) *Migrator {
	name := fmt.Sprintf("%s-migrator", engine)
	core := NewCore(name, opts...)
	return &Migrator{
		Core:     core,
		assetsFS: assetsFS,
		engine:   engine,
	}
}

func (m *Migrator) Setup(ctx context.Context) error {
	switch m.engine {
	case EngSQLite:
		return m.setupSQLite()
	case EngPostgres:
		// NOTE:: Will be implemented later
		return fmt.Errorf("unsupported engine: %s", m.engine)
	default:
		return fmt.Errorf("unsupported engine: %s", m.engine)
	}
}

func (m *Migrator) setupSQLite() error {
	dsn, ok := m.Cfg().StrVal(Key.DBSQLiteDSN)
	if !ok {
		return errors.New("database DSN not found in configuration")
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("cannot connect to SQLite: %w", err)
	}

	m.db = db
	return m.createMigrationsTable()
}

func (m *Migrator) createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id TEXT PRIMARY KEY,
		datetime TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`
	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}
	return nil
}

func (m *Migrator) loadFileMigrations() ([]FileMigration, error) {
	var migrations []FileMigration
	migrationPath := fmt.Sprintf(MigrationPath, m.engine)
	err := fs.WalkDir(m.assetsFS, migrationPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			filename := filepath.Base(path)
			parts := strings.SplitN(filename, "-", 2)
			if len(parts) < 2 {
				return fmt.Errorf("invalid migration filename: %s", filename)
			}

			content, err := m.assetsFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("cannot read migration file %s: %w", path, err)
			}

			sections := strings.Split(string(content), "-- +migrate ")
			var upSection, downSection string
			for _, section := range sections {
				if strings.HasPrefix(section, "Up") {
					upSection = strings.TrimPrefix(section, "Up\n")
				} else if strings.HasPrefix(section, "Down") {
					downSection = strings.TrimPrefix(section, "Down\n")
				}
			}

			migrations = append(migrations, FileMigration{
				Datetime: parts[0],
				Name:     strings.TrimSuffix(parts[1], ".sql"),
				Up:       upSection,
				Down:     downSection,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot load file migrations: %w", err)
	}
	return migrations, nil
}

func (m *Migrator) loadDBMigrations() ([]DBMigration, error) {
	rows, err := m.db.Query("SELECT datetime, name FROM migrations ORDER BY datetime")
	if err != nil {
		return nil, fmt.Errorf("cannot load database migrations: %w", err)
	}
	defer rows.Close()

	var migrations []DBMigration
	for rows.Next() {
		var migration DBMigration
		if err := rows.Scan(&migration.Datetime, &migration.Name); err != nil {
			return nil, fmt.Errorf("cannot scan migration row: %w", err)
		}
		migrations = append(migrations, migration)
	}
	return migrations, nil
}

// Note: We could optimize by only checking the latest migration to determine pending migrations.
// However, for now, we are verifying all migrations to ensure completeness in certain scenarios.
func (m *Migrator) findPendingMigrations(fileMigrations []FileMigration, dbMigrations []DBMigration) []FileMigration {
	dbMigrationsMap := make(map[string]struct{})
	for _, dbMigration := range dbMigrations {
		dbMigrationsMap[dbMigration.Datetime+dbMigration.Name] = struct{}{}
	}

	var pendingMigrations []FileMigration
	for _, fileMigration := range fileMigrations {
		if _, exists := dbMigrationsMap[fileMigration.Datetime+fileMigration.Name]; !exists {
			pendingMigrations = append(pendingMigrations, fileMigration)
		}
	}
	return pendingMigrations
}

func (m *Migrator) logMigrations(fileMigrations []FileMigration, dbMigrations []DBMigration, pendingMigrations []FileMigration) {
	m.Log().Info("File-based migrations:")
	for _, migration := range fileMigrations {
		m.Log().Info(fmt.Sprintf("  %s-%s", migration.Datetime, migration.Name))
	}

	m.Log().Info("Database migrations:")
	for _, migration := range dbMigrations {
		m.Log().Info(fmt.Sprintf("  %s-%s", migration.Datetime, migration.Name))
	}

	m.Log().Info("Pending migrations:")
	for _, migration := range pendingMigrations {
		m.Log().Info(fmt.Sprintf("  %s-%s", migration.Datetime, migration.Name))
	}
}

func (m *Migrator) Start(ctx context.Context) error {
	fileMigrations, err := m.loadFileMigrations()
	if err != nil {
		return err
	}

	dbMigrations, err := m.loadDBMigrations()
	if err != nil {
		return err
	}

	pendingMigrations := m.findPendingMigrations(fileMigrations, dbMigrations)
	m.logMigrations(fileMigrations, dbMigrations, pendingMigrations)

	return m.Migrate(pendingMigrations)
}

func (m *Migrator) Migrate(pendingMigrations []FileMigration) error {
	if m.db == nil {
		return errors.New("database connection is not initialized")
	}

	for _, migration := range pendingMigrations {
		err := m.applyMigration(migration)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) applyMigration(migration FileMigration) error {
	if migration.Up == "" {
		return fmt.Errorf("no Up section found in migration %s-%s", migration.Datetime, migration.Name)
	}

	_, err := m.db.Exec(migration.Up)
	if err != nil {
		return fmt.Errorf("cannot execute migration %s-%s: %w", migration.Datetime, migration.Name, err)
	}

	return m.recordMigration(migration)
}

func (m *Migrator) recordMigration(migration FileMigration) error {
	id := uuid.New().String()
	appliedAt := time.Now().Format(time.RFC3339)

	query := `
	INSERT INTO migrations (id, datetime, name, created_at)
	VALUES (?, ?, ?, ?)`
	_, err := m.db.Exec(query, id, migration.Datetime, migration.Name, appliedAt)
	if err != nil {
		return fmt.Errorf("cannot record migration: %w", err)
	}
	return nil
}
