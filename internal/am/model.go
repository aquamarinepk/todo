package am

import (
	"strings"
	"time"
	"unicode" // Added for normalize function

	"github.com/google/uuid"
)

// Model interface composes Identifiable, Auditable, and Seedable interfaces.
type Model interface {
	Identifiable
	Auditable
	Stampable
	Seedable
}

// Identifiable interface represents an entity with an ID, Slug, and TypeID.
type Identifiable interface {
	// Type returns the type of the entity.
	Type() string
	// ID returns the unique identifier of the entity.
	ID() uuid.UUID
	// GenID generates and sets the unique identifier of the entity if it is not set yet.
	GenID()
	// SetID sets the unique identifier of the entity.
	SetID(id uuid.UUID, force ...bool)
	// ShortID returns the short ID portion of the slug.
	ShortID() string
	// GenShortID generates and sets the short ID if it is not set yet.
	GenShortID()
	// SetShortID sets the short ID of the entity.
	SetShortID(shortID string, force ...bool)
	// TypeID returns a universal identifier for a specific model instance.
	TypeID() string
	// Slug returns the slug of the entity,
	Slug() string
}

// Auditable interface represents an entity with audit information.
type Auditable interface {
	// CreatedBy returns the UUID of the user who created the entity.
	CreatedBy() uuid.UUID
	// UpdatedBy returns the UUID of the user who last updated the entity.
	UpdatedBy() uuid.UUID
	// CreatedAt returns the creation time of the entity.
	CreatedAt() time.Time
	// UpdatedAt returns the last update time of the entity.
	UpdatedAt() time.Time
}

type Stampable interface {
	GenCreateValues(userID ...uuid.UUID) // Modified
	GenUpdateValues(userID ...uuid.UUID) // Modified
}

// Seedable is an interface for entities that need special logic before being inserted into the database during seeding.
type Seedable interface {
	// Ref returns the reference string for this entity (used in seed data for relationships).
	Ref() string
	// SetRef sets the reference string for this entity.
	SetRef(ref string)
}

type BaseModel struct {
	modelType string
	id        uuid.UUID
	shortID   string
	createdBy uuid.UUID
	updatedBy uuid.UUID
	createdAt time.Time
	updatedAt time.Time
	RefValue  string `json:"ref"`
}

// ModelOption defines a functional option for configuring a BaseModel.
type ModelOption func(*BaseModel)

// WithType sets the type of the BaseModel.
func WithType(t string) ModelOption {
	return func(m *BaseModel) {
		m.modelType = t
	}
}

// WithID sets the id of the BaseModel.
func WithID(id uuid.UUID) ModelOption {
	return func(m *BaseModel) {
		m.id = id
	}
}

// WithShortID sets the shortID of the BaseModel.
func WithShortID(shortID string) ModelOption {
	return func(m *BaseModel) {
		m.shortID = shortID
	}
}

// WithCreatedBy sets the createdBy field of the BaseModel.
func WithCreatedBy(createdBy uuid.UUID) ModelOption {
	return func(m *BaseModel) {
		m.createdBy = createdBy
	}
}

// WithUpdatedBy sets the updatedBy field of the BaseModel.
func WithUpdatedBy(updatedBy uuid.UUID) ModelOption {
	return func(m *BaseModel) {
		m.updatedBy = updatedBy
	}
}

// WithCreatedAt sets the createdAt field of the BaseModel.
func WithCreatedAt(createdAt time.Time) ModelOption {
	return func(m *BaseModel) {
		m.createdAt = createdAt
	}
}

// WithUpdatedAt sets the updatedAt field of the BaseModel.
func WithUpdatedAt(updatedAt time.Time) ModelOption {
	return func(m *BaseModel) {
		m.updatedAt = updatedAt
	}
}

// NewModel creates a new BaseModel with the provided options.
func NewModel(options ...ModelOption) *BaseModel {
	m := &BaseModel{}
	for _, opt := range options {
		opt(m)
	}

	return m
}

// Type returns the type of the entity.
func (m *BaseModel) Type() string {
	if m.modelType == "" {
		return "model"
	}

	return m.modelType
}

// SetType sets the type of the entity.
func (m *BaseModel) SetType(t string) {
	m.modelType = t
}

// ID returns the unique identifier of the entity.
func (m *BaseModel) ID() uuid.UUID {
	return m.id
}

// GenID generates and sets the unique identifier of the entity if it is not set
// yet.
func (m *BaseModel) GenID() {
	if m.id == uuid.Nil {
		m.id = uuid.New()
	}
}

// SetID sets the unique identifier of the entity.
func (m *BaseModel) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if m.id == uuid.Nil || (shouldForce && id != uuid.Nil) {
		m.id = id
	}
}

// ShortID returns the short ID portion of the slug.
func (m *BaseModel) ShortID() string {
	return m.shortID
}

// GenShortID generates and sets the short ID if it is not set yet.
func (m *BaseModel) GenShortID() {
	if m.shortID == "" {
		newUUID := uuid.New()
		segments := strings.Split(newUUID.String(), "-")
		m.shortID = segments[len(segments)-1]
	}
}

// SetShortID sets the short ID of the entity.
func (m *BaseModel) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if m.shortID == "" || shouldForce {
		m.shortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
// It combines the normalized model type with its ShortID.
// This is useful for uniquely referencing an entity across different model types.
func (m *BaseModel) TypeID() string {
	return Normalize(m.Type()) + "-" + m.ShortID()
}

// Slug returns a human-readable, URL-friendly string identifier for the entity.
// It typically combines a normalized, recognizable short text (often derived from a title or name)
// with a unique, more random component (like a short ID) to ensure uniqueness
// while maintaining readability. This makes it suitable for use in URLs or as an external reference.
func (m *BaseModel) Slug() string {
	return Normalize(m.Type()) + "-" + m.ShortID()
}

// GenCreateValues sets the values for creation.
// If a userID is provided, it sets the CreatedBy field.
func (m *BaseModel) GenCreateValues(userID ...uuid.UUID) { // Modified
	m.GenID()
	m.GenShortID()
	m.createdAt = time.Now()
	m.updatedAt = m.createdAt
	if len(userID) > 0 {
		m.createdBy = userID[0]
		m.updatedBy = userID[0]
	}
}

// GenUpdateValues sets the values for an update.
// If a userID is provided, it sets the UpdatedBy field.
func (m *BaseModel) GenUpdateValues(userID ...uuid.UUID) { // Modified
	m.updatedAt = time.Now()
	if len(userID) > 0 {
		m.updatedBy = userID[0]
	}
}

// CreatedBy returns the UUID of the user who created the entity.
func (m *BaseModel) CreatedBy() uuid.UUID {
	return m.createdBy
}

// UpdatedBy returns the UUID of the user who last updated the entity.
func (m *BaseModel) UpdatedBy() uuid.UUID {
	return m.updatedBy
}

// CreatedAt returns the creation time of the entity.
func (m *BaseModel) CreatedAt() time.Time {
	return m.createdAt
}

// UpdatedAt returns the last update time of the entity.
func (m *BaseModel) UpdatedAt() time.Time {
	return m.updatedAt
}

func (m *BaseModel) Ref() string {
	return m.RefValue
}

func (m *BaseModel) SetRef(ref string) {
	m.RefValue = ref
}

// normalize replaces space-like characters or non-ASCII characters with '-'
// and converts the string to lowercase.
func Normalize(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsSpace(r) {
			b.WriteRune('-')
		} else if r > unicode.MaxASCII {
			b.WriteRune('-')
		} else {
			b.WriteRune(r)
		}
	}
	return strings.ToLower(b.String())
}
