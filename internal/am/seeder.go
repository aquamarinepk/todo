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
	SeedPath = "assets/seed/%s"
)

type Seeder struct {
	Core
	db       *sql.DB
	assetsFS embed.FS
	engine   string
	seeds    sync.Map
}

type Seed struct {
	Datetime string
	Name     string
	Content  string
}

func NewSeeder(assetsFS embed.FS, engine string, opts ...Option) *Seeder {
	name := fmt.Sprintf("%s-seeder", engine)
	core := NewCore(name, opts...)
	return &Seeder{
		Core:     core,
		assetsFS: assetsFS,
		engine:   engine,
	}
}

func (s *Seeder) Setup(ctx context.Context) error {
	switch s.engine {
	case EngSQLite:
		return s.setupSQLite()
	case EngPostgres:
		// NOTE:: Will be implemented later
		return fmt.Errorf("unsupported engine: %s", s.engine)
	default:
		return fmt.Errorf("unsupported engine: %s", s.engine)
	}
}

func (s *Seeder) setupSQLite() error {
	dsn, ok := s.Cfg().StrVal(Key.DBSQLiteDSN)
	if !ok {
		return errors.New("database DSN not found in configuration")
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("cannot connect to SQLite: %w", err)
	}

	s.db = db
	return s.createSeedsTable()
}

func (s *Seeder) createSeedsTable() error {
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

func (s *Seeder) loadFileSeeds() ([]Seed, error) {
	var seeds []Seed
	seedPath := fmt.Sprintf(SeedPath, s.engine)
	err := fs.WalkDir(s.assetsFS, seedPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			filename := filepath.Base(path)
			parts := strings.SplitN(filename, "-", 2)
			if len(parts) < 2 {
				return fmt.Errorf("invalid seed filename: %s", filename)
			}

			content, err := s.assetsFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("cannot read seed file %s: %w", path, err)
			}

			seeds = append(seeds, Seed{
				Datetime: parts[0],
				Name:     strings.TrimSuffix(parts[1], ".sql"),
				Content:  string(content),
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot load file seeds: %w", err)
	}
	return seeds, nil
}

func (s *Seeder) loadDBSeeds() ([]Seed, error) {
	rows, err := s.db.Query("SELECT datetime, name FROM seeds ORDER BY datetime")
	if err != nil {
		return nil, fmt.Errorf("cannot load database seeds: %w", err)
	}
	defer rows.Close()

	var seeds []Seed
	for rows.Next() {
		var seed Seed
		if err := rows.Scan(&seed.Datetime, &seed.Name); err != nil {
			return nil, fmt.Errorf("cannot scan seed row: %w", err)
		}
		seeds = append(seeds, seed)
	}
	return seeds, nil
}

func (s *Seeder) findPendingSeeds(fileSeeds []Seed, dbSeeds []Seed) []Seed {
	dbSeedsMap := make(map[string]struct{})
	for _, dbSeed := range dbSeeds {
		dbSeedsMap[dbSeed.Datetime+dbSeed.Name] = struct{}{}
	}

	var pendingSeeds []Seed
	for _, fileSeed := range fileSeeds {
		if _, exists := dbSeedsMap[fileSeed.Datetime+fileSeed.Name]; !exists {
			pendingSeeds = append(pendingSeeds, fileSeed)
		}
	}
	return pendingSeeds
}

func (s *Seeder) logSeeds(fileSeeds []Seed, dbSeeds []Seed, pendingSeeds []Seed) {
	s.Log().Info("File-based seeds:")
	for _, seed := range fileSeeds {
		s.Log().Info(fmt.Sprintf("  %s-%s", seed.Datetime, seed.Name))
	}

	s.Log().Info("Database seeds:")
	for _, seed := range dbSeeds {
		s.Log().Info(fmt.Sprintf("  %s-%s", seed.Datetime, seed.Name))
	}

	s.Log().Info("Pending seeds:")
	for _, seed := range pendingSeeds {
		s.Log().Info(fmt.Sprintf("  %s-%s", seed.Datetime, seed.Name))
	}
}

func (s *Seeder) Start(ctx context.Context) error {
	fileSeeds, err := s.loadFileSeeds()
	if err != nil {
		return err
	}

	dbSeeds, err := s.loadDBSeeds()
	if err != nil {
		return err
	}

	pendingSeeds := s.findPendingSeeds(fileSeeds, dbSeeds)
	s.logSeeds(fileSeeds, dbSeeds, pendingSeeds)

	return s.Seed(pendingSeeds)
}

func (s *Seeder) Seed(pendingSeeds []Seed) error {
	if s.db == nil {
		return errors.New("database connection is not initialized")
	}

	for _, seed := range pendingSeeds {
		err := s.applySeed(seed)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) applySeed(seed Seed) error {
	if seed.Content == "" {
		return fmt.Errorf("no content found in seed %s-%s", seed.Datetime, seed.Name)
	}

	_, err := s.db.Exec(seed.Content)
	if err != nil {
		return fmt.Errorf("cannot execute seed %s-%s: %w", seed.Datetime, seed.Name, err)
	}

	return s.recordSeed(seed)
}

func (s *Seeder) recordSeed(seed Seed) error {
	id := uuid.New().String()
	appliedAt := time.Now().Format(time.RFC3339)

	query := `
	INSERT INTO seeds (id, datetime, name, type, context, created_at)
	VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, id, seed.Datetime, seed.Name, "sql", "", appliedAt)
	if err != nil {
		return fmt.Errorf("cannot record seed: %w", err)
	}
	return nil
}
