package am

import (
	"fmt"
	"net/url"
	"path"

	"github.com/google/uuid"
)

// MenuItemStyle defines the styling for a menu item.
type MenuItemStyle string

// Define constants for button styles.
// These styles are expected to be configurable by editing some Sass/CSS when the assets pipeline is in place.
const (
	BtnPrimaryStyle   MenuItemStyle = "bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"     // main action (e.g., "Save", "Submit")
	BtnSecondaryStyle MenuItemStyle = "bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded"     // neutral action (e.g., "Back", "Cancel")
	BtnDangerStyle    MenuItemStyle = "bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"       // destructive action (e.g., "Delete")
	BtnWarningStyle   MenuItemStyle = "bg-yellow-500 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded" // risky but not destructive (e.g., "Override", "Reset")
	BtnInfoStyle      MenuItemStyle = "bg-teal-500 hover:bg-teal-700 text-white font-bold py-2 px-4 rounded"     // informative or contextual actions (optional, e.g., "More info")
	BtnGenericStyle   MenuItemStyle = "bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"   // generic action
)

// MenuItem represents a single item in the menu.
// NOTE: Maybe some of these field can be removed later.
type MenuItem struct {
	Feat        Feat              // The feature this menu item calls
	Text        string            // The text to display for the menu item
	Method      string            // "GET" or "POST"
	IsForm      bool              // Indicates if the action should be triggered via a form submission (POST)
	Style       MenuItemStyle     // The style of the menu item
	QueryParams map[string]string // Query parameters for the URL
	CsrfToken   string            // Only applicable for POST requests
}

// Menu represents the entire menu structure (optional, depending on your needs).
// NOTE: BasePath, FeatName and CsrfToken are stored in order to be able to provide them to the MenuItem.
// but maybe a better approach could be implemented later. This is a WIP
type Menu struct {
	Path      string
	Items     []MenuItem
	CsrfToken string
}

// GenHref constructs the Href from the MenuItem data.
func (i *MenuItem) GenHref() string {
	basePath := path.Join(i.Feat.FeatPath, i.Feat.Action)

	if i.Feat.PathSuffix != "" {
		basePath = path.Join(basePath, i.Feat.PathSuffix)
	}

	if len(i.QueryParams) == 0 {
		return basePath
	}

	query := url.Values{}
	for key, value := range i.QueryParams {
		query.Add(key, value)
	}

	return basePath + "?" + query.Encode()
}

// GenLinkButton generates an HTML link button.
func (i *MenuItem) GenLinkButton() string {
	href := i.GenHref()
	return fmt.Sprintf(`<a href="%s" class="%s">%s</a>`, href, i.Style, i.Text)
}

// NewMenu creates a new Menu with the given parameters.
func NewMenu(path string) *Menu {
	return &Menu{
		Path:  path,
		Items: []MenuItem{},
	}
}

// SetCSRFToken sets the CSRF token for the menu.
func (m *Menu) SetCSRFToken(csrfToken string) {
	m.CsrfToken = csrfToken
}

// AddListItem adds a new MenuItem for listing resources.
func (m *Menu) AddListItem(action string) {
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			FeatPath:   m.FeatName,
			Action:     action,
			PathSuffix: action,
		},
		Text:  "Back",
		Style: BtnSecondaryStyle,
	})
}

// AddEditItem adds a new MenuItem for editing a resource.
func (m *Menu) AddEditItem(action string, id uuid.UUID) {
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			FeatPath:   m.FeatName,
			Action:     action,
			PathSuffix: action,
		},
		Text:  "Edit",
		Style: BtnPrimaryStyle,
		QueryParams: map[string]string{
			"id": id.String(),
		},
	})
}

// AddDeleteItem adds a new MenuItem for deleting a resource.
func (m *Menu) AddDeleteItem(action string, id uuid.UUID) {
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			FeatPath:   m.FeatName,
			Action:     action,
			PathSuffix: action,
		},
		Text:      "Delete",
		Style:     BtnDangerStyle,
		IsForm:    true,
		CsrfToken: m.CsrfToken,
		QueryParams: map[string]string{
			"id": id.String(),
		},
	})
}

// AddGenericItem adds a new generic MenuItem.
func (m *Menu) AddGenericItem(action, url, text string) {
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			FeatPath:   m.FeatName,
			Action:     action,
			PathSuffix: action,
		},
		Text:  text,
		Style: BtnGenericStyle,
		QueryParams: map[string]string{
			"id": url,
		},
	})
}
