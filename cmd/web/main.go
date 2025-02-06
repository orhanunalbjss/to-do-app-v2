package main

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"os"
	"to-do-app-v2/internal/store"
	"to-do-app-v2/internal/web"
)

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceId, ok := ctx.Value(TraceIDHeader).(string); ok {
		r.AddAttrs(slog.String(TraceIDHeader, traceId))
	}
	return h.Handler.Handle(ctx, r)
}

const TraceIDHeader = "TraceID"

func main() {
	ctx := context.WithValue(context.Background(), TraceIDHeader, uuid.NewString())

	setNewDefaultLogger()

	itemStore := store.NewStore()
	itemHandler := web.NewWeb(itemStore)

	router := http.NewServeMux()

	router.HandleFunc("POST /items", itemHandler.HandleHTTPPost)
	router.HandleFunc("GET /items", itemHandler.HandleHTTPGet)
	router.HandleFunc("GET /items/{id}", itemHandler.HandleHTTPGetWithId)
	router.HandleFunc("PUT /items/{id}", itemHandler.HandleHTTPPut)
	router.HandleFunc("DELETE /items/{id}", itemHandler.HandleHTTPDelete)

	if err := http.ListenAndServe(":8080", traceIdMiddleware(router)); err != nil {
		slog.ErrorContext(ctx, err.Error())
	}
}

func setNewDefaultLogger() {
	var handler slog.Handler
	handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})
	handler = &ContextHandler{handler}

	slog.SetDefault(slog.New(handler))
}

func traceIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId := r.Header.Get(TraceIDHeader)
		if _, err := uuid.Parse(traceId); err != nil {
			traceId = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), "TraceID", traceId)
		w.Header().Set(TraceIDHeader, traceId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
