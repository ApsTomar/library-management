package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const RequestTracingID = "requestID"

func RequestTracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracingID := uuid.New().String()
		ctx := context.WithValue(r.Context(), RequestTracingID, tracingID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
