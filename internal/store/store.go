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
	ID     string `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"description"`
	Status string `json:"status"`
}

func (item Item) String() string {
	return fmt.Sprintf("Name: %s, Description: %s, Status: %s", item.Name, item.Desc, item.Status)
}

type ItemID string

func NewItemID() ItemID {
	return ItemID(uuid.New().String())
}

type response struct {
	items []Item
	item  Item
	err   error
}

type request struct {
	action       string
	responseChan chan response
	id           ItemID
	item         Item
}

type Store struct {
	items       map[ItemID]Item
	requestChan chan request
}

const ItemsFilename = "items.json"

func NewStore() *Store {
	s := &Store{
		items:       make(map[ItemID]Item),
		requestChan: make(chan request, 100),
	}

	go s.processRequests()

	return s
}

func (s *Store) processRequests() {
	for req := range s.requestChan {
		switch req.action {
		case "create":
			req.responseChan <- s.create(req.item)
		case "readAll":
			req.responseChan <- s.readAll()
		case "read":
			req.responseChan <- s.read(req.id)
		case "update":
			req.responseChan <- s.update(req.id, req.item)
		case "delete":
			req.responseChan <- s.delete(req.id)
		}
	}
}

func (s *Store) create(item Item) response {
	if err := s.loadItems(); err != nil {
		return response{
			err: errors.Wrap(err, "load items"),
		}
	}

	id := NewItemID()
	item.ID = string(id)
	s.items[id] = item

	if err := s.saveItems(); err != nil {
		return response{
			err: errors.Wrap(err, "save items"),
		}
	}

	return response{
		item: item,
	}
}

func (s *Store) Create(item Item) (Item, error) {
	responseChan := make(chan response, 1)
	req := request{
		action:       "create",
		responseChan: responseChan,
		item:         item,
	}
	s.requestChan <- req
	res := <-responseChan

	return res.item, res.err
}

func (s *Store) readAll() response {
	if err := s.loadItems(); err != nil {
		return response{
			err: errors.Wrap(err, "load items"),
		}
	}

	items := make([]Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}

	return response{
		items: items,
	}
}

func (s *Store) ReadAll() ([]Item, error) {
	responseChan := make(chan response, 1)
	req := request{
		action:       "readAll",
		responseChan: responseChan,
	}
	s.requestChan <- req
	res := <-responseChan

	return res.items, res.err
}

func (s *Store) read(id ItemID) response {
	if err := s.loadItems(); err != nil {
		return response{
			err: errors.Wrap(err, "load items"),
		}
	}

	item, found := s.items[id]
	if !found {
		return response{
			err: errors.Errorf("item '%s' not found", id),
		}
	}

	return response{
		item: item,
	}
}

func (s *Store) Read(id ItemID) (Item, error) {
	responseChan := make(chan response, 1)
	req := request{
		action:       "read",
		responseChan: responseChan,
		id:           id,
	}
	s.requestChan <- req
	res := <-responseChan

	return res.item, res.err
}

func (s *Store) update(id ItemID, item Item) response {
	if err := s.loadItems(); err != nil {
		return response{
			err: errors.Wrap(err, "load items"),
		}
	}

	if _, found := s.items[id]; !found {
		return response{
			err: errors.Errorf("item '%s' not found", id),
		}
	}

	item.ID = string(id)
	s.items[id] = item

	if err := s.saveItems(); err != nil {
		return response{
			err: errors.Wrap(err, "save items"),
		}
	}

	return response{
		item: item,
	}
}

func (s *Store) Update(id ItemID, item Item) (Item, error) {
	responseChan := make(chan response, 1)
	req := request{
		action:       "update",
		responseChan: responseChan,
		id:           id,
		item:         item,
	}
	s.requestChan <- req
	res := <-responseChan

	return res.item, res.err
}

func (s *Store) delete(id ItemID) response {
	if err := s.loadItems(); err != nil {
		return response{
			err: errors.Wrap(err, "load items"),
		}
	}

	if _, found := s.items[id]; !found {
		return response{
			err: errors.Errorf("item '%s' not found", id),
		}
	}

	delete(s.items, id)

	if err := s.saveItems(); err != nil {
		return response{
			err: errors.Wrap(err, "save items"),
		}
	}

	return response{}
}

func (s *Store) Delete(id ItemID) error {
	responseChan := make(chan response, 1)
	req := request{
		action:       "delete",
		responseChan: responseChan,
		id:           id,
	}
	s.requestChan <- req
	res := <-responseChan

	return res.err
}

func (s *Store) loadItems() (err error) {
	var file *os.File
	file, err = os.Open(ItemsFilename)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return errors.Wrapf(err, "open %s", ItemsFilename)
		}

		if err = s.saveItems(); err != nil {
			return errors.Wrap(err, "save items")
		}

		file, err = os.Open(ItemsFilename)
		if err != nil {
			return errors.Wrapf(err, "openb %s", ItemsFilename)
		}
	}

	defer func() {
		closeErr := file.Close()
		if err == nil {
			err = errors.Wrapf(closeErr, "close %s", ItemsFilename)
		}
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s.items)

	if err != nil {
		return errors.Wrapf(err, "decode %s", ItemsFilename)
	}

	return err
}

func (s *Store) saveItems() (err error) {
	var file *os.File
	file, err = os.Create(ItemsFilename)
	if err != nil {
		return errors.Wrapf(err, "create %s", ItemsFilename)
	}

	defer func() {
		closeErr := file.Close()
		if err == nil {
			err = errors.Wrapf(closeErr, "close %s", ItemsFilename)
		}
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s.items)

	if err != nil {
		return errors.Wrapf(err, "encode %s", ItemsFilename)
	}

	return err
}
