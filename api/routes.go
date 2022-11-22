package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (l *HTTPListener) newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", l.homeHandler)
	r.HandleFunc("/api", l.homeHandler)

	r.HandleFunc("/api/cmd", l.cmdHandler).Methods(http.MethodPost)

	return r
}
