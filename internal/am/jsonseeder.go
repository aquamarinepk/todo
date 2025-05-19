package am

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
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
