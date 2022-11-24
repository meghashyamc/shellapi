package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type HTTPListener struct {
	server   *http.Server
	validate *validator.Validate
}

func NewHTTPListener() (*HTTPListener, error) {

	listener := &HTTPListener{validate: newValidator()}
	server := &http.Server{
		Handler:      addRecoveryOptions(addCORSOptions(listener.newRouter())),
		Addr:         fmt.Sprintf(":%d", servicePort),
		WriteTimeout: serverWriteTimeout,
		ReadTimeout:  serverReadTimeout,
	}
	listener.server = server
	return listener, nil

}

func (l *HTTPListener) Listen() {
	lg := defaultLogger()
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
