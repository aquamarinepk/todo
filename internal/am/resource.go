package am

import "github.com/google/uuid"

// Resource interface provides functions that let UI elements render associated elements to it, such as buttons, links,
// etc.
type Resource interface {
	ID() uuid.UUID
	Type() string
}
