package am

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// FormConfig holds configuration for form mapping
type FormConfig struct {
	FieldTag    string // Form field tag name
	RequiredTag string // Required field tag name
}

// DefaultFormConfig returns the default form configuration
func DefaultFormConfig() FormConfig {
	return FormConfig{
		FieldTag:    "form",
		RequiredTag: "required",
	}
}

// TypeConverter converts a string form value to a specific type
type TypeConverter func(string) (interface{}, error)

// FormField represents a form field mapping
type FormField struct {
	Name     string
	Required bool
	Type     reflect.Type
	Convert  TypeConverter
}

// FormMapper maps form values to struct fields
type FormMapper struct {
	config FormConfig
	fields map[string]FormField
}

// NewFormMapper creates a new FormMapper for the given struct type
func NewFormMapper(v interface{}, config ...FormConfig) (*FormMapper, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil, errors.New("v must be a pointer to a struct")
	}

	cfg := DefaultFormConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	mapper := &FormMapper{
		config: cfg,
		fields: make(map[string]FormField),
	}

	typ := val.Elem().Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(cfg.FieldTag)
		if tag == "" {
			continue
		}

		converter := getTypeConverter(field.Type)
		if converter == nil {
			continue
		}

		mapper.fields[tag] = FormField{
			Name:     tag,
			Required: field.Tag.Get(cfg.RequiredTag) == "true",
			Type:     field.Type,
			Convert:  converter,
		}
	}

	return mapper, nil
}

// getTypeConverter returns a TypeConverter for the given type
func getTypeConverter(typ reflect.Type) TypeConverter {
	switch typ.Kind() {
	case reflect.String:
		return func(s string) (interface{}, error) { return s, nil }
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(s string) (interface{}, error) {
			return strconv.ParseInt(s, 10, 64)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(s string) (interface{}, error) {
			return strconv.ParseUint(s, 10, 64)
		}
	case reflect.Bool:
		return func(s string) (interface{}, error) {
			return strconv.ParseBool(s)
		}
	case reflect.Float32, reflect.Float64:
		return func(s string) (interface{}, error) {
			return strconv.ParseFloat(s, 64)
		}
	case reflect.Struct:
		switch typ {
		case reflect.TypeOf(time.Time{}):
			return func(s string) (interface{}, error) {
				return time.Parse(time.RFC3339, s)
			}
		case reflect.TypeOf(uuid.UUID{}):
			return func(s string) (interface{}, error) {
				return uuid.Parse(s)
			}
		}
	}
	return nil
}

// MapForm maps form values to the struct
func (m *FormMapper) MapForm(r *http.Request, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("v must be a pointer to a struct")
	}

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(m.config.FieldTag)
		if tag == "" {
			continue
		}

		formValue := r.FormValue(tag)
		if formValue == "" {
			if field.Tag.Get(m.config.RequiredTag) == "true" {
				return errors.New("required field " + tag + " is missing")
			}
			continue
		}

		fieldValue := val.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		formField, ok := m.fields[tag]
		if !ok || formField.Convert == nil {
			continue
		}

		convertedValue, err := formField.Convert(formValue)
		if err != nil {
			return err
		}

		fieldValue.Set(reflect.ValueOf(convertedValue))
	}

	return nil
}

// ToForm is a convenience function that creates a FormMapper and maps form values
// from an HTTP request to a struct.
func ToForm(r *http.Request, v interface{}, config ...FormConfig) error {
	mapper, err := NewFormMapper(v, config...)
	if err != nil {
		return err
	}
	return mapper.MapForm(r, v)
}
