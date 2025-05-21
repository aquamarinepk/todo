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

type Migration struct {
	Datetime string
	Name     string
	Up       string
	Down     string
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
	var err error
	switch m.engine {
	case EngSQLite:
		err = m.setupSQLite()
	case EngPostgres:
		// NOTE:: Will be implemented later
		err = fmt.Errorf("unsupported engine: %s", m.engine)
	default:
		err = fmt.Errorf("unsupported engine: %s", m.engine)
	}

	if err != nil {
		return err
	}

	return m.SetupMigrations()
}

func (m *Migrator) Start(ctx context.Context) error {
	// return m.SetupMigrations()
	return nil
}

func (m *Migrator) SetupMigrations() error {
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

func (m *Migrator) loadFileMigrations() ([]Migration, error) {
	var migrations []Migration
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

			migrations = append(migrations, Migration{
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

func (m *Migrator) loadDBMigrations() ([]Migration, error) {
	rows, err := m.db.Query("SELECT datetime, name FROM migrations ORDER BY datetime")
	if err != nil {
		return nil, fmt.Errorf("cannot load database migrations: %w", err)
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var migration Migration
		if err := rows.Scan(&migration.Datetime, &migration.Name); err != nil {
			return nil, fmt.Errorf("cannot scan migration row: %w", err)
		}
		migrations = append(migrations, migration)
	}
	return migrations, nil
}

// Note: We could optimize by only checking the latest migration to determine pending ones.
// However, for now, we are verifying all of them to ensure completeness in certain scenarios.
func (m *Migrator) findPendingMigrations(fileMigrations []Migration, dbMigrations []Migration) []Migration {
	dbMigrationsMap := make(map[string]struct{})
	for _, dbMigration := range dbMigrations {
		dbMigrationsMap[dbMigration.Datetime+dbMigration.Name] = struct{}{}
	}

	var pendingMigrations []Migration
	for _, fileMigration := range fileMigrations {
		if _, exists := dbMigrationsMap[fileMigration.Datetime+fileMigration.Name]; !exists {
			pendingMigrations = append(pendingMigrations, fileMigration)
		}
	}
	return pendingMigrations
}

func (m *Migrator) logMigrations(fileMigrations []Migration, dbMigrations []Migration, pendingMigrations []Migration) {
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

func (m *Migrator) Migrate(pendingMigrations []Migration) error {
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

func (m *Migrator) applyMigration(migration Migration) error {
	if migration.Up == "" {
		return fmt.Errorf("no Up section found in migration %s-%s", migration.Datetime, migration.Name)
	}

	_, err := m.db.Exec(migration.Up)
	if err != nil {
		return fmt.Errorf("cannot execute migration %s-%s: %w", migration.Datetime, migration.Name, err)
	}

	return m.recordMigration(migration)
}

func (m *Migrator) recordMigration(migration Migration) error {
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
