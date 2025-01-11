package am

type Status string

const (
	Disabled    Status = "disabled"
	Initialized Status = "initialized"
	Started     Status = "started"
	Stopped     Status = "stopped"
)

type Dep struct {
	core   Core
	Status Status
}

func (d *Dep) Core() Core {
	return d.core
}
