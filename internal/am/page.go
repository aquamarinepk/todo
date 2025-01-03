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
	Button Button
}

type Button struct {
	Text  string
	Style string
}

func NewPage(data interface{}) *Page {
	return &Page{
		Data:  data,
		Flash: Flash{},
		Form: Form{
			Action: "",
			Button: Button{
				Text:  "Submit",
				Style: "",
			},
		},
	}
}

func (p *Page) SetFlash(flash Flash) {
	p.Flash = flash
}

func (p *Page) SetFormAction(action string) {
	p.Form.Action = action
}

func (p *Page) SetFormMethod(method string) {
	p.Form.Method = method
}

func (p *Page) SetFormButton(button Button) {
	p.Form.Button = button
}

func (p *Page) SetFormButtonText(text string) {
	p.Form.Button.Text = text
}

func (p *Page) SetFormButtonStyle(style string) {
	p.Form.Button.Style = style
}

func (p *Page) SetActions(actions []Action) {
	p.Actions = actions
}
