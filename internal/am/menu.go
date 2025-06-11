package am

import (
	"fmt"
	"net/url"
	"path"
)

// MenuItemStyle defines the styling for a menu item.
type MenuItemStyle string

// Define constants for button styles.
// These styles are expected to be configurable by editing some SCSS when the assets pipeline is in place.
const (
	// BtnPrimaryStyle is the main action style (e.g., "Save", "Submit").
	BtnPrimaryStyle MenuItemStyle = "bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"

	// BtnSecondaryStyle is the neutral action style (e.g., "Back", "Cancel").
	BtnSecondaryStyle MenuItemStyle = "bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded"

	// BtnDangerStyle is the destructive action style (e.g., "Delete").
	BtnDangerStyle MenuItemStyle = "bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"

	// BtnWarningStyle is the risky but not destructive action style (e.g., "Override", "Reset").
	BtnWarningStyle MenuItemStyle = "bg-yellow-500 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded"

	// BtnInfoStyle is the informative or contextual action style (optional, e.g., "More info").
	BtnInfoStyle MenuItemStyle = "bg-teal-500 hover:bg-teal-700 text-white font-bold py-2 px-4 rounded"

	// BtnGenericStyle is the generic action style.
	BtnGenericStyle MenuItemStyle = "bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
)

// MenuItem represents a single item in the menu.
// NOTE: Maybe some of these field can be removed later.
type MenuItem struct {
	Feat        Feat              // The feature this menu item calls
	Text        string            // The text to display for the menu item
	Method      string            // "GET" or "POST"
	IsForm      bool              // Indicates if the action should be triggered via a form submission (POST)
	Style       MenuItemStyle     // The style of the menu item
	QueryParams map[string]string // Query parameters for the Path
	CSRFToken   string            // Only applicable for POST requests
}

// Menu represents the entire menu structure (optional, depending on your needs).
// NOTE: Path and CSRFToken are stored in order to be able to provide them to the MenuItem.
// but maybe a better approach could be implemented later. This is a WIP
type Menu struct {
	Path      string
	Items     []MenuItem
	CSRFToken string
}

// GenPath constructs the href path from the MenuItem data.
func (i *MenuItem) GenPath() string {
	basePath := path.Join(i.Feat.Path, i.Feat.Action)

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
	href := i.GenPath()
	return fmt.Sprintf(`<a href="%s" class="%s">%s</a>`, href, i.Style, i.Text)
}

// Path generates the href path for a menu item based on the feature and menu item data.
func (i *MenuItem) Path() string {
	basePath := path.Join(i.Feat.Path, i.Feat.Action)

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

// NewMenu creates a new Menu with the given parameters.
func NewMenu(path string) *Menu {
	return &Menu{
		Path:  path,
		Items: []MenuItem{},
	}
}

// AddListItem adds a new MenuItem for listing resources.
func (m *Menu) AddListItem(resource Resource, text ...string) {
	// TODO: Use a pluralization library to get the plural form of the resource type.
	action := fmt.Sprintf("list-%ss", resource.Type())
	btnText := "Back"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path:   m.Path,
			Action: action,
		},
		Text:  btnText,
		Style: BtnSecondaryStyle,
	})
}

// AddNewItem adds a new MenuItem for creating a new resource.
func (m *Menu) AddNewItem(resourceType string, text ...string) {
	action := fmt.Sprintf("new-%s", resourceType)
	btnText := "New"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path:   m.Path,
			Action: action,
		},
		Text:  btnText,
		Style: BtnPrimaryStyle,
	})
}

// AddShowItem adds a new MenuItem for showing a resource.
func (m *Menu) AddShowItem(resource Resource, text ...string) {
	action := fmt.Sprintf("show-%s", resource.Type())
	btnText := "Show"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path:   m.Path,
			Action: action,
		},
		Text:  btnText,
		Style: BtnPrimaryStyle,
		QueryParams: map[string]string{
			"id": resource.ID().String(),
		},
	})
}

// AddEditItem adds a new MenuItem for editing a resource.
func (m *Menu) AddEditItem(resource Resource, text ...string) {
	action := fmt.Sprintf("edit-%s", resource.Type())
	btnText := "Edit"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path:   m.Path,
			Action: action,
		},
		Text:  btnText,
		Style: BtnPrimaryStyle,
		QueryParams: map[string]string{
			"id": resource.ID().String(),
		},
	})
}

// AddDeleteItem adds a new MenuItem for deleting a resource.
func (m *Menu) AddDeleteItem(resource Resource, text ...string) {
	action := fmt.Sprintf("delete-%s", resource.Type())
	btnText := "Delete"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path:   m.Path,
			Action: action,
		},
		Text:      btnText,
		Style:     BtnDangerStyle,
		IsForm:    true,
		CSRFToken: m.CSRFToken,
		QueryParams: map[string]string{
			"id": resource.ID().String(),
		},
	})
}

// AddGenericItem adds a new generic MenuItem.
func (m *Menu) AddGenericItem(action, url string, text ...string) {
	btnText := "Generic"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path:   m.Path,
			Action: action,
		},
		Text:  btnText,
		Style: BtnGenericStyle,
		QueryParams: map[string]string{
			"id": url,
		},
	})
}
