package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var Items []Item

func main() {
	nameFlag := flag.String("name", "Buy a new backpack", "item name")
	descriptionFlag := flag.String("description", "Prefer red or blue", "item description")
	statusFlag := flag.String("status", "Not Started", "item status")

	flag.Parse()

	item := Item{Name: *nameFlag, Description: *descriptionFlag, Status: *statusFlag}
	Items = append(Items, item)

	fmt.Println(Items)

	file, err := os.Create("items.json")
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
