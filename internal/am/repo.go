package am

import "context"

type Repo struct {
	Core  Core
	Query *QueryManager
}

func NewRepo(name string, qm *QueryManager, opts ...Option) *Repo {
	core := NewCore(name, opts...)
	return &Repo{
		Core:  core,
		Query: qm,
	}
}

func (r *Repo) Name() string {
	return r.Core.Name()
}

func (r *Repo) SetName(name string) {
	r.Core.SetName(name)
}

func (r *Repo) Log() Logger {
	return r.Core.Log()
}

func (r *Repo) SetLog(log Logger) {
	r.Core.SetLog(log)
}

func (r *Repo) Cfg() *Config {
	return r.Core.Cfg()
}

func (r *Repo) SetCfg(cfg *Config) {
	r.Core.SetCfg(cfg)
}

func (r *Repo) Setup(ctx context.Context) error {
	return r.Core.Setup(ctx)
}

func (r *Repo) Start(ctx context.Context) error {
	return r.Core.Start(ctx)
}

func (r *Repo) Stop(ctx context.Context) error {
	return r.Core.Stop(ctx)
}
