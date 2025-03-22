package am

type Keys struct {
	ServerWebHost      string
	ServerWebPort      string
	ServerWebEnabled   string
	ServerAPIHost      string
	ServerAPIPort      string
	ServerAPIEnabled   string
	ServerFeatPath     string
	ServerIndexEnabled string
	
	DBAuthSQLiteDSN string

	SecCSRFKey      string
	SecCSRFRedirect string

	ButtonStyleGray   string
	ButtonStyleBlue   string
	ButtonStyleRed    string
	ButtonStyleGreen  string
	ButtonStyleYellow string

	NotificationSuccessStyle string
	NotificationInfoStyle    string
	NotificationWarnStyle    string
	NotificationErrorStyle   string
	NotificationDebugStyle   string

	RenderWebErrors string
	RenderAPIErrors string
}

var Key = Keys{
	ServerWebHost:      "server.web.host",
	ServerWebPort:      "server.web.port",
	ServerWebEnabled:   "server.web.enabled",
	ServerAPIHost:      "server.api.host",
	ServerAPIPort:      "server.api.port",
	ServerAPIEnabled:   "server.api.enabled",
	ServerFeatPath:     "server.feat.path",
	ServerIndexEnabled: "server.index.enabled",

	DBAuthSQLiteDSN: "db.auth.sqlite.dsn",

	SecCSRFKey:      "sec.csrf.key",
	SecCSRFRedirect: "sec.csrf.redirect",

	ButtonStyleGray:   "button.style.gray",
	ButtonStyleBlue:   "button.style.blue",
	ButtonStyleRed:    "button.style.red",
	ButtonStyleGreen:  "button.style.green",
	ButtonStyleYellow: "button.style.yellow",

	NotificationSuccessStyle: "notification.success.style",
	NotificationInfoStyle:    "notification.info.style",
	NotificationWarnStyle:    "notification.warn.style",
	NotificationErrorStyle:   "notification.error.style",
	NotificationDebugStyle:   "notification.debug.style",

	RenderWebErrors: "render.web.errors",
	RenderAPIErrors: "render.api.errors",
}
