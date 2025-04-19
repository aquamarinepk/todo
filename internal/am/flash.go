package am

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

// FlashCtxKey is a custom type for context keys
type FlashCtxKey string

const (
	// FlashKey is the key used to store flash messages in the context
	FlashKey FlashCtxKey = "aqmflash"

	// FlashCookieName is the name of the cookie used to store flash messages
	FlashCookieName string = "aqmflash"
)

// NotificationTypes defines the types of notifications
type NotificationTypes struct {
	Success string
	Info    string
	Warn    string
	Error   string
	Debug   string
}

// NotificationType holds the default notification types
var NotificationType = NotificationTypes{
	Success: "success",
	Info:    "info",
	Warn:    "warning",
	Error:   "danger",
	Debug:   "debug",
}

// Notification represents a single notification message
type Notification struct {
	Type string
	Msg  string
}

// Flash is a container for notifications
type Flash struct {
	Notifications []Notification
}

// Add adds a new notification to the flash messages
func (f *Flash) Add(typ, msg string) {
	f.Notifications = append(f.Notifications, Notification{
		Type: typ,
		Msg:  msg,
	})
}

// Clear clears all notifications
func (f *Flash) Clear() {
	f.Notifications = []Notification{}
}

// HasMessages returns true if there are any notifications
func (f *Flash) HasMessages() bool {
	return len(f.Notifications) > 0
}

// FlashManager handles the storage and retrieval of flash messages using secure cookies.
// WIP: This is a work in progress. The flash message system is still not available
// to deliver notifications. Some tweaking is still needed to properly handle
// flash messages across requests.
type FlashManager struct {
	Core
	encoder *securecookie.SecureCookie
}

// NewFlashManager creates a new FlashManager instance
func NewFlashManager(opts ...Option) *FlashManager {
	core := NewCore("flash-manager", opts...)
	return &FlashManager{
		Core: core,
	}
}

// Setup initializes the FlashManager with configuration values
func (fm *FlashManager) Setup(ctx context.Context) error {
	err := fm.Core.Setup(ctx)
	if err != nil {
		return err
	}

	cfg := fm.Cfg()

	hashKey := cfg.ByteSliceVal(Key.SecHashKey)
	blockKey := cfg.ByteSliceVal(Key.SecBlockKey)

	if len(hashKey) == 0 || len(blockKey) == 0 {
		return errors.New("missing hashKey or blockKey in configuration")
	}

	fm.encoder = securecookie.New(hashKey, blockKey)
	return nil
}

func (fm *FlashManager) GetFlash(r *http.Request) (f Flash, err error) {
	cookie, err := r.Cookie(FlashCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return f, nil // No flash messages
		}

		return f, err
	}

	err = fm.encoder.Decode(FlashCookieName, cookie.Value, &f)
	if err != nil {
		return f, fmt.Errorf("cannot decode flash messages: %w", err)
	}

	return f, nil
}

func (fm *FlashManager) SetFlash(w http.ResponseWriter, flash Flash) error {
	encoded, err := fm.encoder.Encode(FlashCookieName, flash)
	if err != nil {
		return fmt.Errorf("cannot encode flash messages: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     FlashCookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	return nil
}

func (fm *FlashManager) AddFlash(w http.ResponseWriter, r *http.Request, typ, msg string) error {
	flash, err := fm.GetFlash(r)
	if err != nil {
		return fmt.Errorf("cannot retrieve flash messages: %w", err)
	}

	flash.Add(typ, msg)

	return fm.SetFlash(w, flash)
}

func (fm *FlashManager) ClearFlash(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     FlashCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	})
}

// responseWriter is a wrapper to capture flash messages during the request lifecycle
type responseWriter struct {
	http.ResponseWriter
	flash         Flash
	flashModified bool
}

// WriteHeader overrides the default WriteHeader to track modifications
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.flashModified = true
	rw.ResponseWriter.WriteHeader(statusCode)
}

// AddFlash adds a flash message to the responseWriter
func (rw *responseWriter) AddFlash(typ, msg string) {
	rw.flash.Add(typ, msg)
	rw.flashModified = true
}

// GetFlash retrieves the flash messages
func (rw *responseWriter) GetFlash() Flash {
	return rw.flash
}

// Middleware is a middleware function that handles flash messages
func (fm *FlashManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flash, _ := fm.GetFlash(r)

		ctx := context.WithValue(r.Context(), FlashKey, flash)
		r = r.WithContext(ctx)

		ww := &responseWriter{ResponseWriter: w, flash: flash}

		next.ServeHTTP(ww, r)

		if !ww.flashModified {
			fm.ClearFlash(w)
		}
	})
}

// WithFlashMiddleware adds the FlashManager's middleware to the router.
func WithFlashMiddleware(fm *FlashManager) Option {
	return func(c Core) {
		if router, ok := c.(*Router); ok {
			router.Use(fm.Middleware)
		}
	}
}
