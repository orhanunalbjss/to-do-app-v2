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
		log.Fatal("Expected 'add' or 'update' subcommands")
	}

	switch os.Args[1] {
	case "add":
		if err := addCmd.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		Items = append(Items, Item{*addName, *addDescription, *addStatus})

		fmt.Println("Item added")
	case "update":
		if err := updateCmd.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		if *updateId < 1 || *updateId > len(Items) {
			log.Fatal("Invalid id: ", *updateId)
		}

		Items[*updateId-1].Name = *updateName
		Items[*updateId-1].Description = *updateDescription
		Items[*updateId-1].Status = *updateStatus

		fmt.Println("Item updated")
	case "delete":
		if err := deleteCmd.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		if *deleteId < 1 || *deleteId > len(Items) {
			log.Fatal("Invalid id: ", *deleteId)
		}

		Items = append(Items[:*deleteId-1], Items[*deleteId:]...)

		fmt.Println("Item deleted")

	default:
		log.Fatal("Expected 'add' subcommand'")
	}

	printItems()

	saveItemsToDisk()
}

func printItems() {
	for index, item := range Items {
		fmt.Printf("%d: %s\n", index+1, item)
	}
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
