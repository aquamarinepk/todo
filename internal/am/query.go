package am

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

type QueryManager struct {
	Core
	queries  sync.Map
	assetsFS embed.FS
	engine   string
}

func NewQueryManager(assetsFS embed.FS, engine string, opts ...Option) *QueryManager {
	core := NewCore("query-manager", opts...)
	qm := &QueryManager{
		Core:     core,
		assetsFS: assetsFS,
		engine:   engine,
	}
	return qm
}

func (qm *QueryManager) Load() {
	err := fs.WalkDir(qm.assetsFS, "assets/query", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			qm.loadQueries(path)
		}
		return nil
	})
	if err != nil {
		qm.Log().Error("Failed to load queries: ", err)
	}
}

func (qm *QueryManager) loadQueries(path string) {
	content, err := qm.assetsFS.ReadFile(path)
	if err != nil {
		qm.Log().Error("Failed to read query file: ", err)
		return
	}

	parts := strings.Split(path, string(filepath.Separator))
	if len(parts) < 4 {
		qm.Log().Error("Invalid query file path: ", path)
		return
	}
	engine := parts[2]
	feat := parts[3]
	resource := strings.TrimSuffix(parts[4], ".sql")

	queries := strings.Split(string(content), "\n-- ")
	for _, query := range queries {
		lines := strings.Split(query, "\n")
		if len(lines) > 0 {
			queryName := strings.TrimSpace(lines[0])
			if !isValidQueryName(queryName) {
				continue
			}
			key := engine + ":" + feat + ":" + resource + ":" + queryName
			value := strings.Join(lines[1:], "\n")
			qm.queries.Store(key, strings.TrimSpace(value))
		}
	}
}

func isValidQueryName(queryName string) bool {
	return queryName != "" && !strings.HasPrefix(queryName, "res:") && !strings.HasPrefix(queryName, "Table:")
}

func (qm *QueryManager) Get(feat, resource, queryName string) (string, error) {
	key := qm.engine + ":" + feat + ":" + resource + ":" + queryName
	if query, ok := qm.queries.Load(key); ok {
		return query.(string), nil
	}
	return "", errors.New("query not found")
}

func (qm *QueryManager) Debug() {
	qm.queries.Range(func(key, value interface{}) bool {
		query := value.(string)
		qm.Log().Debugf("%s\n%s\n---", key, query)
		return true
	})
}

func (qm *QueryManager) Setup(ctx context.Context) error {
	qm.Load()
	return nil
}
