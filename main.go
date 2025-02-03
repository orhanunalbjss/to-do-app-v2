package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/fs"
	"log/slog"
	"os"
)

type Item struct {
	Name   string `json:"name"`
	Desc   string `json:"description"`
	Status string `json:"status"`
}

func (item Item) String() string {
	return fmt.Sprintf("Name: %s, Description: %s, Status: %s", item.Name, item.Desc, item.Status)
}

type ContextHandler struct {
	slog.Handler
}

func (handler *ContextHandler) Handle(context context.Context, record slog.Record) error {
	if traceId, ok := context.Value("TraceID").(string); ok {
		record.AddAttrs(slog.String("TraceID", traceId))
	}
	return handler.Handler.Handle(context, record)
}

const ItemsFilename = "items.json"

var Items []Item

func main() {
	var err error
	ctx := context.WithValue(context.Background(), "TraceID", uuid.NewString())

	setNewDefaultLogger()
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		slog.ErrorContext(ctx, "expected one of the following commands: add, list, update, delete")
		os.Exit(1)
	}

	cmd, cmdArgs := args[0], args[1:]
	switch cmd {
	case "add":
		if err = addCommand(cmdArgs); err != nil {
			err = errors.Wrap(err, "add command")
			slog.ErrorContext(ctx, err.Error())
			return
		}

		slog.InfoContext(ctx, "item added")
	case "list":
		if err = listCommand(); err != nil {
			err = errors.Wrap(err, "list command")
			slog.ErrorContext(ctx, err.Error())
			return
		}
	case "update":
		if err = updateCommand(cmdArgs); err != nil {
			err = errors.Wrap(err, "update command")
			slog.ErrorContext(ctx, err.Error())
			return
		}

		slog.InfoContext(ctx, "item updated")
	case "delete":
		if err = deleteCommand(cmdArgs); err != nil {
			err = errors.Wrap(err, "add command")
			slog.ErrorContext(ctx, err.Error())
			return
		}

		slog.InfoContext(ctx, "item deleted")
	default:
		slog.ErrorContext(ctx, "expected one of the following commands: add, list, update, delete")
		os.Exit(1)
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

	if err = loadItems(); err != nil {
		return errors.Wrap(err, "load items")
	}

	item := Item{Name: name, Desc: desc, Status: status}

	addItem(item)
	listItems()
	if err = saveItems(); err != nil {
		return errors.Wrap(err, "save items")
	}

	return nil
}

func listCommand() error {
	var err error

	if err = loadItems(); err != nil {
		return errors.Wrap(err, "load items")
	}

	listItems()

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

	if err = loadItems(); err != nil {
		return errors.Wrap(err, "load items")
	}

	item := Item{Name: name, Desc: desc, Status: status}

	if err = updateItem(id, item); err != nil {
		return errors.Wrap(err, "update item")
	}
	listItems()
	err = saveItems()

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

	if err = loadItems(); err != nil {
		return errors.Wrap(err, "load items")
	}

	if err = deleteItem(id); err != nil {
		return errors.Wrap(err, "delete item")
	}
	listItems()
	err = saveItems()

	return errors.Wrap(err, "save items")
}

func addItem(item Item) {
	Items = append(Items, item)
}

func listItems() {
	for index, item := range Items {
		fmt.Printf("%d: %s\n", index+1, item)
	}
}

func updateItem(id int, item Item) error {
	if !isValidId(id) {
		return fmt.Errorf("invalid id: %d", id)
	}

	Items[id-1].Name = item.Name
	Items[id-1].Desc = item.Desc
	Items[id-1].Status = item.Status

	return nil
}

func deleteItem(id int) error {
	if !isValidId(id) {
		return fmt.Errorf("invalid id: %d", id)
	}

	Items = append(Items[:id-1], Items[id:]...)

	return nil
}

func isValidId(id int) bool {
	return id >= 1 && id <= len(Items)
}

func loadItems() (err error) {
	var file *os.File
	file, err = os.Open(ItemsFilename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			Items = make([]Item, 0)
			if err = saveItems(); err != nil {
				return errors.Wrap(err, "save items")
			}

			return
		}

		return errors.Wrapf(err, "open %s", ItemsFilename)
	}

	defer func() {
		closeError := file.Close()
		if err == nil {
			err = closeError
		}
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Items)

	return errors.Wrapf(err, "decode %s", ItemsFilename)
}

func saveItems() (err error) {
	var file *os.File
	file, err = os.Create(ItemsFilename)
	if err != nil {
		return errors.Wrapf(err, "create %s", ItemsFilename)
	}

	defer func() {
		closeError := file.Close()
		if err == nil {
			err = errors.Wrapf(closeError, "close %s", ItemsFilename)
		}
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(Items)

	return errors.Wrapf(err, "encode %s", ItemsFilename)
}
