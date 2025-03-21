package am

import (
	"fmt"

	"github.com/google/uuid"
)

type Action struct {
	URL        string
	Text       string
	StyleClass string // Add the StyleClass field
	IsForm     bool
}

func NewListAction(basePath, resName, styleClass string) Action {
	return Action{
		URL:        fmt.Sprintf("%s/list-%ss", basePath, resName),
		Text:       "Back",
		StyleClass: styleClass,
	}
}

func NewEditAction(basePath, resName string, id uuid.UUID, styleClass string) Action {
	return Action{
		URL:        fmt.Sprintf("%s/edit-%s?id=%s", basePath, resName, id),
		Text:       "Edit",
		StyleClass: styleClass,
	}
}

func NewDeleteAction(basePath, resName string, id uuid.UUID, styleClass string) Action {
	return Action{
		URL:        fmt.Sprintf("%s/delete-%s?id=%s", basePath, resName, id),
		Text:       "Delete",
		StyleClass: styleClass,
		IsForm:     true,
	}
}

func NewAction(url, text, styleClass string) Action {
	return Action{
		URL:        url,
		Text:       text,
		StyleClass: styleClass,
	}
}
