package main

import (
	"context"
	"flag"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log/slog"
	"os"
	"to-do-app-v2/pkg/store"
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
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		slog.ErrorContext(ctx, "expected one of the following commands: add, list, update, delete")
		return
	}

	cmd, cmdArgs := args[0], args[1:]
	switch cmd {
	case "add":
		if err := addCommand(cmdArgs); err != nil {
			err = errors.Wrap(err, "add command")
			slog.ErrorContext(ctx, err.Error())
			return
		}

		slog.InfoContext(ctx, "item added")
	case "list":
		if err := listCommand(); err != nil {
			err = errors.Wrap(err, "list command")
			slog.ErrorContext(ctx, err.Error())
			return
		}
	case "update":
		if err := updateCommand(cmdArgs); err != nil {
			err = errors.Wrap(err, "update command")
			slog.ErrorContext(ctx, err.Error())
			return
		}

		slog.InfoContext(ctx, "item updated")
	case "delete":
		if err := deleteCommand(cmdArgs); err != nil {
			err = errors.Wrap(err, "add command")
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

func addCommand(args []string) error {
	var err error

	cmd := flag.NewFlagSet("add", flag.ExitOnError)
	var (
		name, desc, status string
	)

	cmd.StringVar(&name, "name", "", "item name")
	cmd.StringVar(&desc, "description", "", "item description")
	cmd.StringVar(&status, "status", "", "item status")

	if err = cmd.Parse(args); err != nil {
		return errors.Wrapf(err, "parse arguments: %v", args)
	}

	if err = store.LoadItems(); err != nil {
		return errors.Wrap(err, "load items")
	}

	item := store.Item{Name: name, Desc: desc, Status: status}

	store.AddItem(item)
	store.ListItems()
	if err = store.SaveItems(); err != nil {
		return errors.Wrap(err, "save items")
	}

	return nil
}

func listCommand() error {
	var err error

	if err = store.LoadItems(); err != nil {
		return errors.Wrap(err, "load items")
	}

	store.ListItems()

	return nil
}

func updateCommand(args []string) error {
	var err error

	cmd := flag.NewFlagSet("update", flag.ExitOnError)
	var (
		id                 int
		name, desc, status string
	)

	cmd.IntVar(&id, "id", 0, "item id")
	cmd.StringVar(&name, "name", "", "item name")
	cmd.StringVar(&desc, "description", "", "item description")
	cmd.StringVar(&status, "status", "", "item status")

	if err = cmd.Parse(args); err != nil {
		return errors.Wrapf(err, "parse arguments: %v", args)
	}

	if err = store.LoadItems(); err != nil {
		return errors.Wrap(err, "load items")
	}

	item := store.Item{Name: name, Desc: desc, Status: status}

	if err = store.UpdateItem(id, item); err != nil {
		return errors.Wrap(err, "update item")
	}
	store.ListItems()
	err = store.SaveItems()

	return errors.Wrap(err, "save items")
}

func deleteCommand(args []string) error {
	var err error

	cmd := flag.NewFlagSet("delete", flag.ExitOnError)

	var id int
	cmd.IntVar(&id, "id", 0, "item id")

	if err = cmd.Parse(args); err != nil {
		return errors.Wrapf(err, "parse arguments: %v", args)
	}

	if err = store.LoadItems(); err != nil {
		return errors.Wrap(err, "load items")
	}

	if err = store.DeleteItem(id); err != nil {
		return errors.Wrap(err, "delete item")
	}
	store.ListItems()
	err = store.SaveItems()

	return errors.Wrap(err, "save items")
}
