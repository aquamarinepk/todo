package am

import (
	"context"
)

type Repo interface {
	Core
	Query() *QueryManager
	BeginTx(ctx context.Context) (context.Context, Tx, error)
}

type txContextKey struct{}

// WithTx returns a new context with the transaction stored.
func WithTx(ctx context.Context, tx interface{}) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

// TxFromContext retrieves the transaction from the context, if present.
func TxFromContext(ctx context.Context) (interface{}, bool) {
	tx := ctx.Value(txContextKey{})
	if tx == nil {
		return nil, false
	}
	return tx, true
}

type Tx interface {
	Commit() error
	Rollback() error
}

func (r BaseRepo) BeginTx(ctx context.Context) (updatedCtx context.Context, tx Tx, err error) {
	return updatedCtx, nil, err
}

type BaseRepo struct {
	*BaseCore
	query *QueryManager
}

func NewRepo(name string, qm *QueryManager, opts ...Option) *BaseRepo {
	core := NewCore(name, opts...)
	return &BaseRepo{
		BaseCore: core,
		query:    qm,
	}
}

func (r *BaseRepo) Query() *QueryManager {
	return r.query
}
