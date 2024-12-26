package am

type ConfigKeys struct {
	WebHostKey string
	WebPortKey string
	WebEnabled string
	APIHostKey string
	APIPortKey string
	APIEnabled string
}

var Keys = ConfigKeys{
	WebHostKey: "server.web.host",
	WebPortKey: "server.web.port",
	WebEnabled: "server.web.enabled",
	APIHostKey: "server.api.host",
	APIPortKey: "server.api.port",
	APIEnabled: "server.api.enabled",
}
