package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (l *HTTPListener) newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", l.homeHandler).Methods(http.MethodGet).Name(RouteHome)
	r.HandleFunc("/api", l.homeHandler).Methods(http.MethodGet).Name(RouteAPIHome)
	r.HandleFunc("/api/cmd", l.cmdHandler).Methods(http.MethodPost).Name(RouteCommand)
	r.Use(tracingMiddleware)
	return r
}
