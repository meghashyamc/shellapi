package api

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func tracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		gottenTraceID := r.Header.Get(traceKey)
		traceID, err := uuid.Parse(gottenTraceID)
		if err != nil {
			traceID = uuid.New()
		}

		ctx := context.WithValue(r.Context(), traceKey, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func addCORSOptions(r *mux.Router) http.Handler {

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "X-Forwarded-Authorization", "Content-Type", "Access-Control-Allow-Origin", "Authorization", "X-API-Key", "Accept", "Accept-Encoding", "X-Request-Id", "Content-Length", "User-Agent"})
	originsOk := handlers.AllowedOrigins(strings.Split(os.Getenv("ORIGIN_ALLOWED"), ","))
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodOptions, http.MethodDelete, http.MethodPut, http.MethodPatch})
	return handlers.CORS(originsOk, headersOk, methodsOk)(r)
}

func addRecoveryOptions(next http.Handler) http.Handler {

	return handlers.RecoveryHandler(handlers.PrintRecoveryStack(true), handlers.RecoveryLogger(recoveryLogger{}))(next)
}
