package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"strconv"
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

	fmt.Println("Options")
	fmt.Println("1. Create")
	fmt.Println("2. Print all")
	fmt.Println("3. Print one")
	fmt.Println("4. Update")
	fmt.Println("5. Delete")
	fmt.Println("6. Exit")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter choice (1, 2, 3, 4, 5, 6): ")
		scanner.Scan()
		choice, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Invalid choice:", choice)
			continue
		}

		switch choice {
		case 1:
			fmt.Print("Enter item name: ")
			scanner.Scan()
			name := scanner.Text()
			fmt.Print("Enter item description: ")
			scanner.Scan()
			desc := scanner.Text()
			fmt.Print("Enter item status: ")
			scanner.Scan()
			status := scanner.Text()

			item := store.Item{Name: name, Desc: desc, Status: status}
			if _, err = itemStore.Create(item); err != nil {
				slog.ErrorContext(ctx, err.Error())
				continue
			}

			slog.InfoContext(ctx, "item added")
		case 2:
			var items []store.Item
			items, err = itemStore.ReadAll()
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				continue
			}
			for _, item := range items {
				fmt.Println(item)
			}
		case 3:
			fmt.Print("Enter item id: ")
			scanner.Scan()
			id := scanner.Text()

			itemID := store.ItemID(id)
			var item store.Item
			item, err = itemStore.Read(itemID)
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				continue
			}
			fmt.Println(item)
		case 4:
			fmt.Print("Enter item id: ")
			scanner.Scan()
			id := scanner.Text()
			fmt.Print("Enter new item name: ")
			scanner.Scan()
			name := scanner.Text()
			fmt.Print("Enter new item description: ")
			scanner.Scan()
			desc := scanner.Text()
			fmt.Print("Enter new item status: ")
			scanner.Scan()
			status := scanner.Text()

			itemID := store.ItemID(id)
			item := store.Item{Name: name, Desc: desc, Status: status}
			if _, err = itemStore.Update(itemID, item); err != nil {
				slog.ErrorContext(ctx, err.Error())
				continue
			}

			slog.InfoContext(ctx, "item updated")
		case 5:
			fmt.Print("Enter item id: ")
			scanner.Scan()
			id := scanner.Text()

			itemID := store.ItemID(id)
			if err = itemStore.Delete(itemID); err != nil {
				slog.ErrorContext(ctx, err.Error())
				continue
			}

			slog.InfoContext(ctx, "item deleted")
		case 6:
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice:", choice)
		}
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
