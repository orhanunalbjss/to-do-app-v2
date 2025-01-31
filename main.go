package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log/slog"
	"os"
)

const ItemsFilename = "items.json"

type ContextHandler struct {
	slog.Handler
}

func (handler *ContextHandler) Handle(context context.Context, record slog.Record) error {
	if traceId, ok := context.Value("TraceID").(string); ok {
		record.AddAttrs(slog.String("TraceID", traceId))
	}
	return handler.Handler.Handle(context, record)
}

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (item Item) String() string {
	return fmt.Sprintf("Name: %s, Description: %s, Status: %s", item.Name, item.Description, item.Status)
}

var Items []Item

var (
	addFlagSet     = flag.NewFlagSet("add", flag.ExitOnError)
	addName        = addFlagSet.String("name", "", "item name")
	addDescription = addFlagSet.String("description", "", "item description")
	addStatus      = addFlagSet.String("status", "", "item status")

	updateFlagSet     = flag.NewFlagSet("update", flag.ExitOnError)
	updateId          = updateFlagSet.Int("id", 0, "item id")
	updateName        = updateFlagSet.String("name", "", "item name")
	updateDescription = updateFlagSet.String("description", "", "item description")
	updateStatus      = updateFlagSet.String("status", "", "item status")

	deleteFlagSet = flag.NewFlagSet("delete", flag.ExitOnError)
	deleteId      = deleteFlagSet.Int("id", 0, "item id")
)

func main() {
	var err error

	var handler slog.Handler
	handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})
	handler = &ContextHandler{handler}

	slog.SetDefault(slog.New(handler))

	traceId := uuid.NewString()
	ctx := context.WithValue(context.Background(), "TraceID", traceId)

	if len(os.Args) < 2 {
		slog.ErrorContext(ctx, "expected one of the following commands: add, list, update, delete")
		os.Exit(1)
	}

	if err = loadItems(); err != nil {
		logErrorAndExit(ctx, err, "failed to load items")
	}

	command := os.Args[1]
	commandArgs := os.Args[2:]

	switch command {
	case "add":
		if err = processAddCommand(commandArgs); err != nil {
			logErrorAndExit(ctx, err, "failed to process add command")
		}
		slog.InfoContext(ctx, "item added")
	case "list":
		listItems()
	case "update":
		if err = processUpdateCommand(commandArgs); err != nil {
			logErrorAndExit(ctx, err, "failed to process update command")
		}
		slog.InfoContext(ctx, "item updated")
	case "delete":
		if err = processDeleteCommand(commandArgs); err != nil {
			logErrorAndExit(ctx, err, "failed to process delete command")
		}
		slog.InfoContext(ctx, "item deleted")
	default:
		slog.ErrorContext(ctx, "expected one of the following commands: add, list, update, delete")
		os.Exit(1)
	}
}

func logErrorAndExit(ctx context.Context, err error, message string) {
	err = errors.Wrap(err, message)
	slog.ErrorContext(ctx, err.Error())
	os.Exit(1)
}

func processAddCommand(args []string) error {
	var err error
	if err = addFlagSet.Parse(args); err != nil {
		return errors.Wrapf(err, "failed to parse arguments: %v", args)
	}
	addItem(addName, addDescription, addStatus)
	listItems()
	err = saveItems()
	return errors.Wrap(err, "failed to save items")
}

func processUpdateCommand(args []string) error {
	var err error
	if err = updateFlagSet.Parse(args); err != nil {
		return errors.Wrapf(err, "failed to parse arguments: %v", args)
	}
	if err = updateItem(updateId, updateName, updateDescription, updateStatus); err != nil {
		return errors.Wrap(err, "failed to update item")
	}
	listItems()
	err = saveItems()
	return errors.Wrap(err, "failed to save items")
}

func processDeleteCommand(args []string) error {
	var err error
	if err = deleteFlagSet.Parse(args); err != nil {
		return errors.Wrapf(err, "failed to parse arguments: %v", args)
	}
	if err = deleteItem(deleteId); err != nil {
		return errors.Wrap(err, "failed to delete item")
	}
	listItems()
	err = saveItems()
	return errors.Wrap(err, "failed to save items")
}

func addItem(name *string, description *string, status *string) {
	Items = append(Items, Item{*name, *description, *status})
}

func listItems() {
	for index, item := range Items {
		fmt.Printf("%d: %s\n", index+1, item)
	}
}

func updateItem(id *int, name *string, description *string, status *string) error {
	if !isValidId(id) {
		return fmt.Errorf("invalid id: %d", *id)
	}

	Items[*id-1].Name = *name
	Items[*id-1].Description = *description
	Items[*id-1].Status = *status

	return nil
}

func deleteItem(id *int) error {
	if !isValidId(id) {
		return fmt.Errorf("invalid id: %d", *id)
	}

	Items = append(Items[:*id-1], Items[*id:]...)

	return nil
}

func isValidId(id *int) bool {
	return id != nil && *id >= 1 && *id <= len(Items)
}

func loadItems() (err error) {
	var file *os.File
	file, err = os.Open(ItemsFilename)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", ItemsFilename)
	}

	defer func() {
		closeError := file.Close()
		if closeError == nil {
			err = errors.Wrapf(err, "failed to close %s", ItemsFilename)
		}
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Items)

	return errors.Wrapf(err, "failed to decode %s", ItemsFilename)
}

func saveItems() (err error) {
	var file *os.File
	file, err = os.Create(ItemsFilename)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s", ItemsFilename)
	}

	defer func() {
		closeError := file.Close()
		if err == nil {
			err = errors.Wrapf(closeError, "failed to close %s", ItemsFilename)
		}
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(Items)

	return errors.Wrapf(err, "failed to encode %s", ItemsFilename)
}
