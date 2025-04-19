package am

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/securecookie"
)

type NotificationTypes struct {
	Success string
	Info    string
	Warn    string
	Error   string
	Debug   string
}

var NotificationType = NotificationTypes{
	Success: "success",
	Info:    "info",
	Warn:    "warning",
	Error:   "danger",
	Debug:   "debug",
}

type Notification struct {
	Type string
	Msg  string
}

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

// FlashCookieName is the name of the cookie used to store flash messages
const FlashCookieName = "aqm_flash"

// FlashMiddleware handles flash messages using secure cookies
type FlashMiddleware struct {
	sc *securecookie.SecureCookie
}

// NewFlashMiddleware creates a new FlashMiddleware instance
func NewFlashMiddleware(encKey []byte) *FlashMiddleware {
	return &FlashMiddleware{
		sc: securecookie.New(encKey, nil),
	}
}

// Middleware returns the flash middleware handler
func (fm *FlashMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get flash messages from cookie
		flash := Flash{}
		if cookie, err := r.Cookie(FlashCookieName); err == nil {
			var value string
			if err := fm.sc.Decode(FlashCookieName, cookie.Value, &value); err == nil {
				json.Unmarshal([]byte(value), &flash)
			}
		}

		// Add flash to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, FlashCookieName, &flash)
		r = r.WithContext(ctx)

		// Clear flash cookie
		http.SetCookie(w, &http.Cookie{
			Name:     FlashCookieName,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
		})

		next.ServeHTTP(w, r)
	})
}

// SetFlash sets a flash message in the response
func (fm *FlashMiddleware) SetFlash(w http.ResponseWriter, r *http.Request, typ, msg string) {
	// Get existing flash messages
	flash := Flash{}
	if cookie, err := r.Cookie(FlashCookieName); err == nil {
		var value string
		if err := fm.sc.Decode(FlashCookieName, cookie.Value, &value); err == nil {
			json.Unmarshal([]byte(value), &flash)
		}
	}

	// Add new message
	flash.Add(typ, msg)

	// Encode and set cookie
	if encoded, err := json.Marshal(flash); err == nil {
		if value, err := fm.sc.Encode(FlashCookieName, string(encoded)); err == nil {
			http.SetCookie(w, &http.Cookie{
				Name:     FlashCookieName,
				Value:    value,
				Path:     "/",
				MaxAge:   60, // 1 minute
				HttpOnly: true,
				Secure:   true,
			})
		}
	}
}

// GetFlash gets flash messages from the request context
func (fm *FlashMiddleware) GetFlash(r *http.Request) *Flash {
	if flash, ok := r.Context().Value(FlashCookieName).(*Flash); ok {
		return flash
	}
	return &Flash{}
}
