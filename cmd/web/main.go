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

	fs := http.FileServer(http.Dir("./web/static/"))
	router.Handle("/web/static/", http.StripPrefix("/web/static/", fs))

	router.HandleFunc("/about/", itemHandler.HandleAboutPage)
	router.HandleFunc("/list/", itemHandler.HandleListItemsPage)

	router.HandleFunc("POST /items", itemHandler.HandleCreateItem)
	router.HandleFunc("GET /items", itemHandler.HandleGetItems)
	router.HandleFunc("GET /items/{id}", itemHandler.HandleGetItemWithID)
	router.HandleFunc("PUT /items/{id}", itemHandler.HandleUpdateItem)
	router.HandleFunc("DELETE /items/{id}", itemHandler.HandleDeleteItem)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: middleware.TraceIDMiddleware(router),
	}

	if err := srv.ListenAndServe(); err != nil {
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
