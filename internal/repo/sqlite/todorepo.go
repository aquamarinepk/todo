package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/res/todo"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	queryGetAll    = "SELECT id, name, slug FROM lists"
	queryGetByID   = "SELECT id, name, slug FROM lists WHERE id = ?"
	queryGetBySlug = "SELECT id, name, slug FROM lists WHERE slug = ?"
	queryCreate    = "INSERT INTO lists (id, name, slug) VALUES (?, ?, ?)"
	queryUpdate    = "UPDATE lists SET name = ?, slug = ? WHERE id = ?"
	queryDelete    = "DELETE FROM lists WHERE id = ?"
)

type TodoRepo struct {
	core am.Core
	db   *sql.DB
}

func NewRepo(dsn string, opts ...am.Option) (*TodoRepo, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	return &TodoRepo{
		core: am.NewCore("sqlite-repo", opts...),
		db:   db,
	}, nil
}

func (repo *TodoRepo) GetAll(ctx context.Context) ([]todo.List, error) {
	rows, err := repo.db.QueryContext(ctx, queryGetAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []todo.List
	for rows.Next() {
		var list todo.List
		if err := rows.Scan(&list.ID, &list.Name, &list.Slug); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, nil
}

func (repo *TodoRepo) GetByID(ctx context.Context, id uuid.UUID) (todo.List, error) {
	var list todo.List
	err := repo.db.QueryRowContext(ctx, queryGetByID, id).Scan(&list.ID, &list.Name, &list.Slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return list, errors.New("list not found")
		}
		return list, err
	}
	return list, nil
}

func (repo *TodoRepo) GetBySlug(ctx context.Context, slug string) (todo.List, error) {
	var list todo.List
	err := repo.db.QueryRowContext(ctx, queryGetBySlug, slug).Scan(&list.ID, &list.Name, &list.Slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return list, errors.New("list not found")
		}
		return list, err
	}
	return list, nil
}

func (repo *TodoRepo) Create(ctx context.Context, list todo.List) error {
	_, err := repo.db.ExecContext(ctx, queryCreate, list.ID, list.Name, list.Slug)
	return err
}

func (repo *TodoRepo) Update(ctx context.Context, list todo.List) error {
	_, err := repo.db.ExecContext(ctx, queryUpdate, list.Name, list.Slug, list.ID)
	return err
}

func (repo *TodoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.ExecContext(ctx, queryDelete, id)
	return err
}

// Implementing the am.Core interface methods by delegating to the core field
func (repo *TodoRepo) Name() string {
	return repo.core.Name()
}

func (repo *TodoRepo) SetName(name string) {
	repo.core.SetName(name)
}

func (repo *TodoRepo) Log() am.Logger {
	return repo.core.Log()
}

func (repo *TodoRepo) SetLog(log am.Logger) {
	repo.core.SetLog(log)
}

func (repo *TodoRepo) Cfg() *am.Config {
	return repo.core.Cfg()
}

func (repo *TodoRepo) SetCfg(cfg *am.Config) {
	repo.core.SetCfg(cfg)
}

func (repo *TodoRepo) Setup(ctx context.Context) error {
	return repo.core.Setup(ctx)
}

func (repo *TodoRepo) Start(ctx context.Context) error {
	return repo.core.Start(ctx)
}

func (repo *TodoRepo) Stop(ctx context.Context) error {
	return repo.core.Stop(ctx)
}
