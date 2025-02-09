package am

type Service struct {
	Core
}

func NewService(name string, opts ...Option) *Service {
	core := NewCore(name, opts...)
	return &Service{
		Core: core,
	}
}
