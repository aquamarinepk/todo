package am

import (
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/csrf"
)

// Page struct represents a web page with data, flash messages, form, menu, and feature information.
type Page struct {
	Data  interface{}
	Flash Flash
	Form  Form
	Menu  Menu
	Feat  Feat
}

// Form struct represents a form with action, method, CSRF token, and a button.
type Form struct {
	Action string
	Method string
	CSRF   string
	Button Button
}

// Button struct represents a button with text and style.
type Button struct {
	Text  string
	Style string
}

// Feat struct represents a feature with base path, feature name, and action.
type Feat struct {
	BasePath   string // The root path for features
	FeatPath   string // The path of a particular feature
	Action     string // Action name (command or query)
	PathSuffix string // The path suffix for the feature
}

// NewFeat creates a new Feat with the given parameters.
func NewFeat(basePath, featName, action, pathSufix string) Feat {
	return Feat{
		BasePath:   basePath,
		FeatPath:   featName,
		Action:     action,
		PathSuffix: pathSufix,
	}
}

// NewPage creates a new Page with the given data.
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
		Menu: Menu{
			Items: []MenuItem{},
		},
	}
}

// SetFlash sets the flash message for the page.
func (p *Page) SetFlash(flash Flash) {
	p.Flash = flash
}

// SetFormAction sets the form action for the page.
func (p *Page) SetFormAction(action string) {
	p.Form.Action = action
}

// SetFormMethod sets the form method for the page.
func (p *Page) SetFormMethod(method string) {
	p.Form.Method = method
}

// SetFormButton sets the form button for the page.
func (p *Page) SetFormButton(button Button) {
	p.Form.Button = button
}

// SetFormButtonText sets the form button text for the page.
func (p *Page) SetFormButtonText(text string) {
	p.Form.Button.Text = text
}

// SetFormButtonStyle sets the form button style for the page.
func (p *Page) SetFormButtonStyle(style string) {
	p.Form.Button.Style = style
}

// SetFeat sets the Feat struct for the page.
func (p *Page) SetFeat(feat Feat) {
	p.Feat = feat
}

// SetMenuItems sets the menu items for the page.
func (p *Page) SetMenuItems(items []MenuItem) {
	p.Menu.Items = items
}

// GenCSRFToken generates a CSRF token and sets it in the form.
func (p *Page) GenCSRFToken(r *http.Request) {
	p.Form.CSRF = csrf.Token(r)
}

// GenerateHref generates the href for a menu item based on the feature and menu item data.
func (p *Page) GenerateHref(feat Feat, item MenuItem) string {
	basePath := path.Join(feat.BasePath, feat.FeatPath, feat.Action)

	if len(item.QueryParams) == 0 {
		return basePath
	}

	query := url.Values{}
	for key, value := range item.QueryParams {
		query.Add(key, value)
	}

	return basePath + "?" + query.Encode()
}
