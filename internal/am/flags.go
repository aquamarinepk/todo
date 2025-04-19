package am

var Flags = map[string]interface{}{
	Key.ServerWebHost:    "localhost",
	Key.ServerWebPort:    "8080",
	Key.ServerWebEnabled: true,
	Key.ServerAPIHost:    "localhost",
	Key.ServerAPIPort:    "8081",
	Key.ServerAPIEnabled: true,
	Key.ServerResPath:    "/res",

	Key.SecHashKey:  "0123456789abcdef0123456789abcdef",
	Key.SecBlockKey: "0123456789abcdef0123456789abcdef",
}
