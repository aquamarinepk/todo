package am

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type JSONSeed struct {
	Datetime string
	Name     string
	Content  string
}

type JSONSeeder struct {
	Core
	assetsFS embed.FS
	engine   string
	db       *sql.DB
}

func NewJSONSeeder(assetsFS embed.FS, engine string, opts ...Option) *JSONSeeder {
	name := fmt.Sprintf("%s-json-seeder", engine)
	core := NewCore(name, opts...)
	return &JSONSeeder{
		Core:     core,
		assetsFS: assetsFS,
		engine:   engine,
	}
}

// LoadJSONSeeds loads all JSON seed files by feature (e.g. "auth")
func (s *JSONSeeder) LoadJSONSeeds() (map[string][]JSONSeed, error) {
	seedsByFeature := make(map[string][]JSONSeed)
	seedPath := filepath.Join("assets", "seed", s.engine)
	err := fs.WalkDir(s.assetsFS, seedPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			filename := filepath.Base(path)
			parts := strings.SplitN(filename, "-", 2)
			if len(parts) < 2 {
				return fmt.Errorf("invalid seed filename: %s", filename)
			}
			content, err := s.assetsFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("cannot read seed file %s: %w", path, err)
			}
			feature := strings.SplitN(parts[1], "-", 2)[0]
			seedsByFeature[feature] = append(seedsByFeature[feature], JSONSeed{
				Datetime: parts[0],
				Name:     strings.TrimSuffix(parts[1], ".json"),
				Content:  string(content),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return seedsByFeature, nil
}

func (s *JSONSeeder) Setup(ctx context.Context) error {
	dsn, ok := s.Cfg().StrVal(Key.DBSQLiteDSN)
	if !ok {
		return fmt.Errorf("database DSN not found in configuration")
	}
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("cannot connect to SQLite: %w", err)
	}
	s.db = db
	return s.createSeedsTable()
}

func (s *JSONSeeder) createSeedsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS seeds (
		id TEXT PRIMARY KEY,
		datetime TEXT NOT NULL,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		context TEXT,
		created_at TIMESTAMP NOT NULL
	)`
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("cannot create seeds table: %w", err)
	}
	return nil
}

func (s *JSONSeeder) recordJSONSeed(datetime, name, context string) error {
	id := uuid.New().String()
	appliedAt := time.Now().Format(time.RFC3339)
	query := `
	INSERT INTO seeds (id, datetime, name, type, context, created_at)
	VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, id, datetime, name, "json", context, appliedAt)
	if err != nil {
		return fmt.Errorf("cannot record json seed: %w", err)
	}
	return nil
}

func (s *JSONSeeder) SeedApplied(datetime, name, context string) (bool, error) {
	if s.db == nil {
		return false, fmt.Errorf("database connection is not initialized")
	}
	row := s.db.QueryRow("SELECT COUNT(1) FROM seeds WHERE datetime = ? AND name = ? AND type = ? AND context = ?", datetime, name, "json", context)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("cannot check if seed is applied: %w", err)
	}
	return count > 0, nil
}

func (s *JSONSeeder) ApplyJSONSeed(datetime, name, context, content string) error {
	if s.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	applied, err := s.SeedApplied(datetime, name, context)
	if err != nil {
		return err
	}
	if applied {
		s.Log().Debugf("Seed already applied: %s-%s [%s]", datetime, name, context)
		return nil
	}
	// For now, just log for demonstration
	s.Log().Debugf("Applying JSON seed: %s-%s [%s]", datetime, name, context)
	// ...apply the seed content...
	return s.recordJSONSeed(datetime, name, context)
}
