package main

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"os"
	"to-do-app-v2/api/handler"
	"to-do-app-v2/api/middleware"
	"to-do-app-v2/internal/store"
)

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID, ok := ctx.Value(TraceIDHeader).(string); ok {
		r.AddAttrs(slog.String(TraceIDHeader, traceID))
	}
	return h.Handler.Handle(ctx, r)
}

const TraceIDHeader = "TraceID"

func main() {
	ctx := context.WithValue(context.Background(), TraceIDHeader, uuid.NewString())

	setNewDefaultLogger()

	itemStore := store.NewStore()
	itemHandler := handler.NewHandler(itemStore)

	router := http.NewServeMux()

	router.HandleFunc("POST /items", itemHandler.HandleHTTPPost)
	router.HandleFunc("GET /items", itemHandler.HandleHTTPGet)
	router.HandleFunc("GET /items/{id}", itemHandler.HandleHTTPGetWithID)
	router.HandleFunc("PUT /items/{id}", itemHandler.HandleHTTPPut)
	router.HandleFunc("DELETE /items/{id}", itemHandler.HandleHTTPDelete)

	if err := http.ListenAndServe(":8080", middleware.TraceIDMiddleware(router)); err != nil {
		slog.ErrorContext(ctx, err.Error())
	}
}

func setNewDefaultLogger() {
	var h slog.Handler
	h = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})
	h = &ContextHandler{h}

	slog.SetDefault(slog.New(h))
}
