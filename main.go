package main

import (
	"flag"
	"fmt"
)

type Item struct {
	name        string
	description string
	status      string
}

var items []Item

func main() {
	nameFlag := flag.String("name", "Buy a new backpack", "item name")
	descriptionFlag := flag.String("description", "item description", "item description")
	statusFlag := flag.String("status", "not started", "item status")

	flag.Parse()

	item := Item{name: *nameFlag, description: *descriptionFlag, status: *statusFlag}
	items = append(items, item)

	fmt.Println(items)
}
