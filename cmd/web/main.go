package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	js "to-do-app-v2/pkg/jsonstore"
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

	router := http.NewServeMux()

	router.HandleFunc("POST /create", addItem(ctx))
	router.HandleFunc("GET /get", getItems(ctx))
	router.HandleFunc("PUT /update/{id}", updateItem(ctx))
	router.HandleFunc("DELETE /delete/{id}", deleteItem(ctx))

	if err := http.ListenAndServe(":8080", router); err != nil {
		slog.ErrorContext(ctx, "Error starting server")
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

func addItem(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var item js.Item

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&item); err != nil {
			err = errors.Wrap(err, "decode item")
			http.Error(w, err.Error(), http.StatusBadRequest)
			slog.ErrorContext(ctx, err.Error())
			return
		}

		if err := js.LoadItems(); err != nil {
			err = errors.Wrap(err, "load items")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.ErrorContext(ctx, err.Error())
			return
		}

		js.AddItem(item)

		if err := js.SaveItems(); err != nil {
			err = errors.Wrap(err, "save items")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.ErrorContext(ctx, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func getItems(ctx context.Context) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := js.LoadItems(); err != nil {
			err = errors.Wrap(err, "load items")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.ErrorContext(ctx, err.Error())
			return
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(js.Items); err != nil {
			err = errors.Wrap(err, "encode items")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.ErrorContext(ctx, err.Error())
			return
		}
	}
}

func updateItem(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			err = errors.Wrap(err, "parse id")
			http.Error(w, err.Error(), http.StatusBadRequest)
			slog.ErrorContext(ctx, err.Error())
			return
		}

		var item js.Item
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&item); err != nil {
			err = errors.Wrap(err, "decode item")
			http.Error(w, err.Error(), http.StatusBadRequest)
			slog.ErrorContext(ctx, err.Error())
			return
		}

		if err = js.LoadItems(); err != nil {
			err = errors.Wrap(err, "load items")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.ErrorContext(ctx, err.Error())
			return
		}

		err = js.UpdateItem(id, item)
		if err != nil {
			err = errors.Wrap(err, "update item")
			http.Error(w, err.Error(), http.StatusBadRequest)
			slog.ErrorContext(ctx, err.Error())
			return
		}

		if err = js.SaveItems(); err != nil {
			err = errors.Wrap(err, "save items")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.ErrorContext(ctx, err.Error())
			return
		}
	}
}

func deleteItem(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			err = errors.Wrap(err, "parse id")
			http.Error(w, err.Error(), http.StatusBadRequest)
			slog.ErrorContext(ctx, err.Error())
			return
		}

		if err = js.LoadItems(); err != nil {
			err = errors.Wrap(err, "load items")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.ErrorContext(ctx, err.Error())
			return
		}

		err = js.DeleteItem(id)
		if err != nil {
			err = errors.Wrap(err, "delete item")
			http.Error(w, err.Error(), http.StatusBadRequest)
			slog.ErrorContext(ctx, err.Error())
			return
		}

		if err = js.SaveItems(); err != nil {
			err = errors.Wrap(err, "save items")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.ErrorContext(ctx, err.Error())
			return
		}
	}
}
