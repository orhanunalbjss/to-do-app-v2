package main

import (
	"context"
	"github.com/google/uuid"
	"html/template"
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

type TodoPageData struct {
	PageTitle string
	Items     []store.Item
}

const TraceIDHeader = "TraceID"

func main() {
	ctx := context.WithValue(context.Background(), TraceIDHeader, uuid.NewString())

	setNewDefaultLogger()
	tmpl, _ := template.ParseFiles("./web/templates/list.html")

	itemStore := store.NewStore()
	itemHandler := handler.NewHandler(itemStore)

	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("./web/static/"))
	router.Handle("/about/", http.StripPrefix("/about/", fs))

	router.HandleFunc("GET /list", func(w http.ResponseWriter, r *http.Request) {
		items, err := itemStore.ReadAll()
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return
		}

		data := TodoPageData{
			PageTitle: "My TODO list",
			Items:     items,
		}

		if err := tmpl.Execute(w, data); err != nil {
			slog.ErrorContext(ctx, err.Error())
			return
		}
	})

	router.HandleFunc("POST /items", itemHandler.HandleHTTPPost)
	router.HandleFunc("GET /items", itemHandler.HandleHTTPGet)
	router.HandleFunc("GET /items/{id}", itemHandler.HandleHTTPGetWithID)
	router.HandleFunc("PUT /items/{id}", itemHandler.HandleHTTPPut)
	router.HandleFunc("DELETE /items/{id}", itemHandler.HandleHTTPDelete)

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
