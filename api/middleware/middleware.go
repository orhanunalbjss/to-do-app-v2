package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const TraceIDHeader = "TraceID"

func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(TraceIDHeader)
		if _, err := uuid.Parse(traceID); err != nil {
			traceID = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), "TraceID", traceID)
		w.Header().Set(TraceIDHeader, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
