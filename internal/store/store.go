package store

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/fs"
	"os"
)

type Item struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"description"`
	Status string `json:"status"`
}

func (item Item) String() string {
	return fmt.Sprintf("Name: %s, Description: %s, Status: %s", item.Name, item.Desc, item.Status)
}

type ItemId string

func NewItemID() ItemId {
	return ItemId(uuid.New().String())
}

type Store struct {
	items map[ItemId]Item
}

const ItemsFilename = "items.json"

func NewStore() *Store {
	return &Store{
		items: make(map[ItemId]Item),
	}
}

func (s *Store) Create(item Item) (Item, error) {
	if err := s.LoadItems(); err != nil {
		return Item{}, errors.Wrap(err, "load items")
	}

	id := NewItemID()
	item.Id = string(id)
	s.items[id] = item

	if err := s.SaveItems(); err != nil {
		return Item{}, errors.Wrap(err, "save items")
	}

	return item, nil
}

func (s *Store) ReadAll() ([]Item, error) {
	if err := s.LoadItems(); err != nil {
		return []Item{}, errors.Wrap(err, "load items")
	}

	items := make([]Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}

	return items, nil
}

func (s *Store) Read(id ItemId) (Item, error) {
	if err := s.LoadItems(); err != nil {
		return Item{}, errors.Wrap(err, "load items")
	}

	item, found := s.items[id]
	if !found {
		return Item{}, errors.New("item not found")
	}

	return item, nil
}

func (s *Store) Update(id ItemId, item Item) (Item, error) {
	if err := s.LoadItems(); err != nil {
		return Item{}, errors.Wrap(err, "load items")
	}

	if _, found := s.items[id]; !found {
		return Item{}, errors.New("item not found")
	}

	item.Id = string(id)
	s.items[id] = item

	if err := s.SaveItems(); err != nil {
		return Item{}, errors.Wrap(err, "save items")
	}

	return item, nil
}

func (s *Store) Delete(id ItemId) error {
	if err := s.LoadItems(); err != nil {
		return errors.Wrap(err, "load items")
	}

	if _, found := s.items[id]; !found {
		return errors.New("item not found")
	}

	delete(s.items, id)

	if err := s.SaveItems(); err != nil {
		return errors.Wrap(err, "save items")
	}

	return nil
}

func (s *Store) LoadItems() (err error) {
	var file *os.File
	file, err = os.Open(ItemsFilename)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return errors.Wrapf(err, "opena %s", ItemsFilename)
		}

		s.items = make(map[ItemId]Item)

		if err = s.SaveItems(); err != nil {
			return errors.Wrap(err, "save items")
		}

		file, err = os.Open(ItemsFilename)
		if err != nil {
			return errors.Wrapf(err, "openb %s", ItemsFilename)
		}
	}

	defer func() {
		closeError := file.Close()
		if err == nil {
			err = errors.Wrapf(closeError, "close %s", ItemsFilename)
		}
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s.items)

	if err != nil {
		return errors.Wrapf(err, "decode %s", ItemsFilename)
	}

	return err
}

func (s *Store) SaveItems() (err error) {
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
	err = encoder.Encode(s.items)

	if err != nil {
		return errors.Wrapf(err, "encode %s", ItemsFilename)
	}

	return err
}

func (s *Store) PrintItems() error {
	if err := s.LoadItems(); err != nil {
		return errors.Wrap(err, "load items")
	}

	for _, item := range s.items {
		fmt.Println(item)
	}

	return nil
}
