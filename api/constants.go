package api

const (
	// general
	servicePort = 8000

	// middleware
	traceKey     = "trace-id"
	serviceKey   = "service"
	serviceValue = "shellapi"
	routeKey     = "route"

	// route names
	RouteCommand = "command"
	RouteHome    = "home"
	RouteAPIHome = "api-home"

	// validation
	validateRequired = "required"
	validateMax      = "max"
	validateMin      = "min"
	validateAlphabet = "alpha"
)
