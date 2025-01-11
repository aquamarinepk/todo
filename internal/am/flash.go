package am

type NotificationTypes struct {
	Succes string
	Info   string
	Warn   string
	Error  string
	Debug  string
}

var NotificationType = NotificationTypes{
	Succes: "success",
	Info:   "info",
	Warn:   "warning",
	Error:  "danger",
	Debug:  "debug",
}

type Notification struct {
	Type string
	Msg  string
}

type Flash struct {
	Notifications []Notification
}
