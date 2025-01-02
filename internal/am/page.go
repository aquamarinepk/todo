package am

// Page functionally should be moved to `am` package
type Page struct {
	Data    interface{}
	Flash   Flash
	Form    Form
	Actions []Action
}

type Form struct {
	Action string
	Method string
}

func NewPage(data interface{}) *Page {
	return &Page{
		Data:  data,
		Flash: Flash{},
		Form:  Form{Action: ""},
	}
}

func (p *Page) SetFlash(flash Flash) {
	p.Flash = flash
}

func (p *Page) SetFormAction(action string) {
	p.Form.Action = action
}

func (p *Page) SetActions(actions []Action) {
	p.Actions = actions
}
