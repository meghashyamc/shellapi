package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	shutdownTime       = 5 * time.Second
	serverWriteTimeout = 60 * time.Second
	serverReadTimeout  = 60 * time.Second
)

type HTTPListener struct {
	server   *http.Server
	validate *validator.Validate
}

func NewHTTPListener() (*HTTPListener, error) {

	listener := &HTTPListener{validate: newValidator()}
	server := &http.Server{
		Handler:      addCORSOptions(listener.newRouter()),
		Addr:         fmt.Sprintf(":%d", servicePort),
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	listener.server = server
	return listener, nil

}

func (l *HTTPListener) Listen() {
	lg := logger("none", context.Background())
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := l.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.WithFields(log.Fields{"err": err.Error()}).Error("HTTP listener exited with an error")
		}
	}()

	lg.WithFields(log.Fields{"address": l.server.Addr}).Info("server started, listening successfully")
	signalReceived := <-done
	lg.WithFields(log.Fields{"signal": signalReceived.String()}).Info("server stopped because of signal")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer func() {
		cancel()
	}()

	if err := l.server.Shutdown(ctx); err != nil {
		lg.WithFields(log.Fields{"err": err.Error()}).Error("server shutdown failed")
		return
	}
	lg.Info("server exited gracefully")

}

func addCORSOptions(r *mux.Router) http.Handler {

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "X-Forwarded-Authorization", "Content-Type", "Access-Control-Allow-Origin", "Authorization", "X-API-Key", "Accept", "Accept-Encoding", "X-Request-Id", "Content-Length", "User-Agent"})
	originsOk := handlers.AllowedOrigins(strings.Split(os.Getenv("ORIGIN_ALLOWED"), ","))
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodOptions, http.MethodDelete, http.MethodPut, http.MethodPatch})
	return handlers.CORS(originsOk, headersOk, methodsOk)(r)
}
