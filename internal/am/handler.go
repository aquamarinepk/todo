package am

type Handler struct {
	Core
}

func NewHandler(name string, opts ...Option) *Handler {
	core := NewCore(name, opts...)
	return &Handler{
		Core: core,
	}
}
