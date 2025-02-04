package jsonstore

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/fs"
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

const ItemsFilename = "items.json"

var Items []Item

func AddItem(item Item) {
	Items = append(Items, item)
}

func ListItems() {
	for index, item := range Items {
		fmt.Printf("%d: %s\n", index+1, item)
	}
}

func UpdateItem(id int, item Item) error {
	if !isValidId(id) {
		return fmt.Errorf("invalid id: %d", id)
	}

	Items[id-1].Name = item.Name
	Items[id-1].Desc = item.Desc
	Items[id-1].Status = item.Status

	return nil
}

func DeleteItem(id int) error {
	if !isValidId(id) {
		return fmt.Errorf("invalid id: %d", id)
	}

	Items = append(Items[:id-1], Items[id:]...)

	return nil
}

func isValidId(id int) bool {
	return id >= 1 && id <= len(Items)
}

func LoadItems() (err error) {
	var file *os.File
	file, err = os.Open(ItemsFilename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			Items = make([]Item, 0)
			if err = SaveItems(); err != nil {
				return errors.Wrap(err, "save items")
			}

			return err
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

	if err != nil {
		return errors.Wrapf(err, "decode %s", ItemsFilename)
	}

	return err
}

func SaveItems() (err error) {
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

	if err != nil {
		return errors.Wrapf(err, "encode %s", ItemsFilename)
	}

	return err
}
