package api

import "time"

const (
	// general
	servicePort        = 8000
	shutdownTime       = 5 * time.Second
	serverWriteTimeout = 60 * time.Second
	serverReadTimeout  = 60 * time.Second

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
