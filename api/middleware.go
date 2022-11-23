package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
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
