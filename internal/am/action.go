package am

import (
	"github.com/google/uuid"
)

type Action struct {
	Path   string
	Text   string
	Style  string
	IsForm bool
}

func ListAction(basePath, resName, style string) Action {
	return Action{
		Path:  ListPath(basePath, resName),
		Text:  "Back",
		Style: style,
	}
}

func EditAction(basePath, resName string, id uuid.UUID, style string) Action {
	return Action{
		Path:  EditPath(basePath, resName, id),
		Text:  "Edit",
		Style: style,
	}
}

func DeleteAction(basePath, resName string, id uuid.UUID, style string) Action {
	return Action{
		Path:   DeletePath(basePath, resName),
		Text:   "Delete",
		Style:  style,
		IsForm: true,
	}
}

func NewAction(url, text, style string) Action {
	return Action{
		Path:  url,
		Text:  text,
		Style: style,
	}
}
