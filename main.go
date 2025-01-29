package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

const ItemsFilename = "items.json"

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var Items []Item

func (item Item) String() string {
	return fmt.Sprintf("Name: %s, Description: %s, Status: %s", item.Name, item.Description, item.Status)
}

func main() {
	loadItemsFromDisk()

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addName := addCmd.String("name", "", "item name")
	addDescription := addCmd.String("description", "", "item description")
	addStatus := addCmd.String("status", "", "item status")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateId := updateCmd.Int("id", 0, "item id")
	updateName := updateCmd.String("name", "", "item name")
	updateDescription := updateCmd.String("description", "", "item description")
	updateStatus := updateCmd.String("status", "", "item status")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteId := deleteCmd.Int("id", 0, "item id")

	if len(os.Args) < 2 {
		log.Fatal("Expected 'add', 'list', 'update' or 'delete' subcommands")
	}

	switch os.Args[1] {
	case "add":
		parseSubcommands(addCmd)
		addItem(addName, addDescription, addStatus)
		fmt.Println("Item added")
		listItems()
	case "list":
		listItems()
	case "update":
		parseSubcommands(updateCmd)
		updateItem(updateId, updateName, updateDescription, updateStatus)
		fmt.Println("Item updated")
		listItems()
	case "delete":
		parseSubcommands(deleteCmd)
		deleteItem(deleteId)
		fmt.Println("Item deleted")
		listItems()

	default:
		log.Fatal("Expected 'add', 'list', 'update' or 'delete' subcommands")
	}

	saveItemsToDisk()
}

func loadItemsFromDisk() {
	file, err := os.Open(ItemsFilename)
	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&Items); err != nil {
		log.Fatal(err)
	}
}

func parseSubcommands(cmd *flag.FlagSet) {
	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}
}

func addItem(addName *string, addDescription *string, addStatus *string) {
	Items = append(Items, Item{*addName, *addDescription, *addStatus})
}

func listItems() {
	for index, item := range Items {
		fmt.Printf("%d: %s\n", index+1, item)
	}
}

func updateItem(updateId *int, updateName *string, updateDescription *string, updateStatus *string) {
	validateId(updateId)

	Items[*updateId-1].Name = *updateName
	Items[*updateId-1].Description = *updateDescription
	Items[*updateId-1].Status = *updateStatus
}

func deleteItem(deleteId *int) {
	validateId(deleteId)

	Items = append(Items[:*deleteId-1], Items[*deleteId:]...)
}

func validateId(id *int) {
	if *id < 1 || *id > len(Items) {
		log.Fatal("Invalid id: ", *id)
	}
}

func saveItemsToDisk() {
	file, err := os.Create(ItemsFilename)
	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(Items); err != nil {
		log.Fatal(err)
	}
}
