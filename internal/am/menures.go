// Package am provides core functionality for the application.
// This file contains RESTful menu item creation methods.
// WIP: This is a temporary solution to handle RESTful routes differently from command-query routes.
// In the future, we should consider a more unified approach that can handle both RESTful and command-query patterns.

package am

// AddResListItem adds a new MenuItem for listing resources in a RESTful way.
func (m *Menu) AddResListItem(resource Resource, text ...string) {
	btnText := "Back"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path: m.Path,
		},
		Text:  btnText,
		Style: BtnSecondaryStyle,
	})
}

// AddResNewItem adds a new MenuItem for creating a new resource in a RESTful way.
func (m *Menu) AddResNewItem(resourceType string, text ...string) {
	btnText := "New"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path:   m.Path,
			Action: "new",
		},
		Text:  btnText,
		Style: BtnPrimaryStyle,
	})
}

// AddResShowItem adds a new MenuItem for showing a resource in a RESTful way.
func (m *Menu) AddResShowItem(resource Resource, text ...string) {
	btnText := "Show"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path:   m.Path,
			Action: resource.ID().String(),
		},
		Text:  btnText,
		Style: BtnPrimaryStyle,
	})
}

// AddResEditItem adds a new MenuItem for editing a resource in a RESTful way.
func (m *Menu) AddResEditItem(resource Resource, text ...string) {
	btnText := "Edit"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path:   m.Path,
			Action: resource.ID().String() + "/edit",
		},
		Text:  btnText,
		Style: BtnPrimaryStyle,
	})
}

// AddResDeleteItem adds a new MenuItem for deleting a resource in a RESTful way.
func (m *Menu) AddResDeleteItem(resource Resource, text ...string) {
	btnText := "Delete"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path:   m.Path,
			Action: resource.ID().String(),
		},
		Text:      btnText,
		Style:     BtnDangerStyle,
		IsForm:    true,
		CSRFToken: m.CSRFToken,
	})
}

// AddResGenericItem adds a new MenuItem for a generic RESTful action.
func (m *Menu) AddResGenericItem(action, id string, text ...string) {
	btnText := "Generic"
	if len(text) > 0 {
		btnText = text[0]
	}
	m.Items = append(m.Items, MenuItem{
		Feat: Feat{
			Path:   m.Path,
			Action: id + "/" + action,
		},
		Text:  btnText,
		Style: BtnGenericStyle,
	})
}
