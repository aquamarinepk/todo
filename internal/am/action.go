package am

import (
	"fmt"

	"github.com/google/uuid"
)

type Action struct {
	URL    string
	Text   string
	Style  string
	IsForm bool
}

func NewListAction(basePath, resName, style string) Action {
	return Action{
		URL:   fmt.Sprintf("%s/list-%ss", basePath, resName),
		Text:  "Back",
		Style: style,
	}
}

func NewEditAction(basePath, resName string, id uuid.UUID, style string) Action {
	return Action{
		URL:   fmt.Sprintf("%s/edit-%s?id=%s", basePath, resName, id),
		Text:  "Edit",
		Style: style,
	}
}

func NewDeleteAction(basePath, resName string, id uuid.UUID, style string) Action {
	return Action{
		URL:    fmt.Sprintf("%s/delete-%s?id=%s", basePath, resName, id),
		Text:   "Delete",
		Style:  style,
		IsForm: true,
	}
}

func NewAction(url, text, style string) Action {
	return Action{
		URL:   url,
		Text:  text,
		Style: style,
	}
}
