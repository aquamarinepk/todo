package am

type HTTPMethods struct {
	GET    string
	POST   string
	PUT    string
	PATCH  string
	DELETE string
	HEAD   string
}

var HTTPMethod = HTTPMethods{
	GET:    "GET",
	POST:   "POST",
	PUT:    "PUT",
	PATCH:  "PATCH",
	DELETE: "DELETE",
	HEAD:   "HEAD",
}
