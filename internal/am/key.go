package am

type Keys struct {
	ServerWebHost      string
	ServerWebPort      string
	ServerWebEnabled   string
	ServerAPIHost      string
	ServerAPIPort      string
	ServerAPIEnabled   string
	ServerResPath      string
	ServerIndexEnabled string

	DBSQLiteDSN string

	SecCSRFKey       string
	SecCSRFRedirect  string
	SecEncryptionKey string
	SecHashKey       string
	SecBlockKey      string

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
	ServerResPath:      "server.res.path",
	ServerIndexEnabled: "server.index.enabled",

	DBSQLiteDSN: "db.sqlite.dsn",

	SecCSRFKey:       "sec.csrf.key",
	SecCSRFRedirect:  "sec.csrf.redirect",
	SecEncryptionKey: "sec.encryption.key",
	SecHashKey:       "sec.hash.key",
	SecBlockKey:      "sec.block.key",

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
