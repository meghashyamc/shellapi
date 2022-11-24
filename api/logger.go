package api

import (
	"context"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func LogSetup() {
	log.SetFormatter(&logrus.JSONFormatter{})
}

func routeLogger(routeName string, ctx context.Context) *log.Entry {
	entry := log.WithFields(log.Fields{serviceKey: serviceValue, routeKey: routeName})
	if traceKeyVal := ctx.Value(traceKey); traceKeyVal != nil {
		entry = entry.WithField(traceKey, traceKeyVal)

	}
	return entry
}

func defaultLogger() *log.Entry {
	entry := log.WithFields(log.Fields{serviceKey: serviceValue})
	return entry
}

// used for panic recovery
type recoveryLogger struct {
}

func (l recoveryLogger) Println(v ...interface{}) {
	entry := log.WithFields(log.Fields{serviceKey: serviceValue})
	entry.Errorln(v...)
}
