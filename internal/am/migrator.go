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
		// TODO: Setup SQLite
	case EngPostgres:
		return fmt.Errorf("unsupported engine: %s", m.engine)
	default:
		return fmt.Errorf("unsupported engine: %s", m.engine)
	}
	return m.load()
}

func (m *Migrator) Start(ctx context.Context) error {
	return nil
}

func (m *Migrator) Stop(ctx context.Context) error {
	return m.Close()
}

func (m *Migrator) ConnectSQLite(dataSourceName string) error {
	db, err := sql.Open("sqlite3", dataSourceName)
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

func (m *Migrator) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *Migrator) load() error {
	migrationPath := fmt.Sprintf(MigrationPath, m.engine)
	err := fs.WalkDir(m.assetsFS, migrationPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			return m.loadMigration(path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("cannot load migrations: %w", err)
	}
	return nil
}

func (m *Migrator) loadMigration(path string) error {
	content, err := m.assetsFS.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read migration file: %w", err)
	}

	parts := strings.Split(string(content), "-- +migrate ")
	for _, part := range parts {
		lines := strings.SplitN(part, "\n", 2)
		if len(lines) < 2 {
			continue
		}
		direction := strings.TrimSpace(lines[0])
		script := strings.TrimSpace(lines[1])
		m.migrations.Store(filepath.Base(path)+":"+direction, script)
	}
	return nil
}

func (m *Migrator) Migrate() error {
	if m.db == nil {
		return errors.New("database connection is not initialized")
	}
	return m.exec("Up")
}

func (m *Migrator) Rollback() error {
	if m.db == nil {
		return errors.New("database connection is not initialized")
	}
	return m.exec("Down")
}

func (m *Migrator) exec(direction string) error {
	var err error
	m.migrations.Range(func(key, value interface{}) bool {
		if strings.HasSuffix(key.(string), ":"+direction) {
			_, err = m.db.Exec(value.(string))
			if err != nil {
				err = fmt.Errorf("cannot execute migration %s: %w", key, err)
				return false
			}
			if direction == "Up" {
				err = m.rec(key.(string))
				if err != nil {
					return false
				}
			}
		}
		return true
	})
	return err
}

func (m *Migrator) rec(migration string) error {
	parts := strings.Split(migration, ":")
	if len(parts) < 2 {
		return fmt.Errorf("invalid migration format: %s", migration)
	}
	filename := parts[0]
	datetime := filename[:14]
	name := strings.TrimSuffix(filename[15:], ".sql")
	id := uuid.New().String()
	appliedAt := time.Now().Format(time.RFC3339)

	query := `
	INSERT INTO migrations (id, datetime, name, created_at)
	VALUES (?, ?, ?, ?)`
	_, err := m.db.Exec(query, id, datetime, name, appliedAt)
	if err != nil {
		return fmt.Errorf("cannot record migration: %w", err)
	}
	return nil
}
