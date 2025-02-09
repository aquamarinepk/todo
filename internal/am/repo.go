package am

type Repo struct {
	Core
	Query *QueryManager
}

func NewRepo(name string, qm *QueryManager, opts ...Option) *Repo {
	core := NewCore(name, opts...)
	return &Repo{
		Core:  core,
		Query: qm,
	}
}
