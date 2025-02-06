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

func (handler *ContextHandler) Handle(context context.Context, record slog.Record) error {
	if traceId, ok := context.Value("TraceID").(string); ok {
		record.AddAttrs(slog.String("TraceID", traceId))
	}
	return handler.Handler.Handle(context, record)
}

func main() {
	ctx := context.WithValue(context.Background(), "TraceID", uuid.NewString())

	setNewDefaultLogger()

	itemStore := store.NewStore()
	itemHandler := web.NewWeb(itemStore)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /items", itemHandler.HandleHTTPPost)
	mux.HandleFunc("GET /items", itemHandler.HandleHTTPGet)
	mux.HandleFunc("GET /items/{id}", itemHandler.HandleHTTPGetWithId)
	mux.HandleFunc("PUT /items/{id}", itemHandler.HandleHTTPPut)
	mux.HandleFunc("DELETE /items/{id}", itemHandler.HandleHTTPDelete)

	if err := http.ListenAndServe(":8080", mux); err != nil {
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
