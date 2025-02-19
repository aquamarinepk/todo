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
}

func NewQueryManager(assetsFS embed.FS, opts ...Option) *QueryManager {
	core := NewCore("query-manager", opts...)
	qm := &QueryManager{
		Core:     core,
		assetsFS: assetsFS,
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
	resource := parts[3]
	res := strings.TrimSuffix(parts[4], ".sql")

	queries := strings.Split(string(content), "-- ")
	for _, query := range queries {
		lines := strings.Split(query, "\n")
		if len(lines) > 1 {
			queryName := strings.TrimSpace(lines[0])
			if !isValidQueryName(queryName) {
				continue
			}
			key := engine + ":" + resource + ":" + res + ":" + queryName
			value := strings.Join(lines[1:], "\n")
			qm.queries.Store(key, strings.TrimSpace(value))
		}
	}
}

func isValidQueryName(queryName string) bool {
	return queryName != "" && !strings.HasPrefix(queryName, "res:") && !strings.HasPrefix(queryName, "Table:")
}

func (qm *QueryManager) Get(engine, resource, res, queryName string) (string, error) {
	key := engine + ":" + resource + ":" + res + ":" + queryName
	if query, ok := qm.queries.Load(key); ok {
		return query.(string), nil
	}
	return "", errors.New("query not found")
}

func (qm *QueryManager) Debug() {
	qm.queries.Range(func(key, value interface{}) bool {
		query := value.(string)
		qm.Log().Infof("Query key: %s, Query: %s", key, query)
		return true
	})
}

func (qm *QueryManager) Setup(ctx context.Context) error {
	qm.Load()
	return nil
}
