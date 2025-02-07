package main

import (
	"context"
	"flag"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"to-do-app-v2/internal/app"
	"to-do-app-v2/internal/store"
)

type ContextHandler struct {
	slog.Handler
}

func (handler *ContextHandler) Handle(context context.Context, record slog.Record) error {
	if traceID, ok := context.Value("TraceID").(string); ok {
		record.AddAttrs(slog.String("TraceID", traceID))
	}
	return handler.Handler.Handle(context, record)
}

func main() {
	ctx := context.WithValue(context.Background(), "TraceID", uuid.NewString())

	setNewDefaultLogger()

	itemStore := store.NewStore()
	newCli := app.NewCli(itemStore)

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		slog.ErrorContext(ctx, "expected one of the following commands: add, list, update, delete")
		return
	}

	cmd, cmdArgs := args[0], args[1:]
	switch cmd {
	case "add":
		if err := newCli.AddCommand(cmdArgs); err != nil {
			slog.ErrorContext(ctx, err.Error())
			return
		}
		slog.InfoContext(ctx, "item added")
	case "list":
		if err := newCli.ListCommand(); err != nil {
			slog.ErrorContext(ctx, err.Error())
			return
		}
	case "update":
		if err := newCli.UpdateCommand(cmdArgs); err != nil {
			slog.ErrorContext(ctx, err.Error())
			return
		}
		slog.InfoContext(ctx, "item updated")
	case "delete":
		if err := newCli.DeleteCommand(cmdArgs); err != nil {
			slog.ErrorContext(ctx, err.Error())
			return
		}
		slog.InfoContext(ctx, "item deleted")
	default:
		slog.ErrorContext(ctx, "expected one of the following commands: add, list, update, delete")
		return
	}

	// Keep app running until SIGINT (CTRL+C) signal is sent
	/*
		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, syscall.SIGINT)
		<-quitChannel
	*/
}

func setNewDefaultLogger() {
	var handler slog.Handler
	handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})
	handler = &ContextHandler{handler}

	slog.SetDefault(slog.New(handler))
}
