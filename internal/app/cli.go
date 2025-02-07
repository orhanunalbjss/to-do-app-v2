package app

import (
	"flag"
	"github.com/pkg/errors"
	"to-do-app-v2/internal/store"
)

type Cli struct {
	store *store.Store
}

func NewCli(store *store.Store) *Cli {
	return &Cli{
		store: store,
	}
}

func (c *Cli) AddCommand(args []string) error {
	cmd := flag.NewFlagSet("add", flag.ExitOnError)
	var (
		name, desc, status string
	)

	cmd.StringVar(&name, "name", "", "store name")
	cmd.StringVar(&desc, "description", "", "store description")
	cmd.StringVar(&status, "status", "", "store status")

	if err := cmd.Parse(args); err != nil {
		return errors.Wrapf(err, "parse arguments: %v", args)
	}

	newItem := store.Item{Name: name, Desc: desc, Status: status}

	if _, err := c.store.Create(newItem); err != nil {
		return errors.Wrap(err, "create item")
	}

	if err := c.store.PrintItems(); err != nil {
		return errors.Wrap(err, "print items")
	}

	return nil
}

func (c *Cli) ListCommand() error {
	if err := c.store.PrintItems(); err != nil {
		return errors.Wrap(err, "print items")
	}

	return nil
}

func (c *Cli) UpdateCommand(args []string) error {
	cmd := flag.NewFlagSet("update", flag.ExitOnError)
	var (
		id                 string
		name, desc, status string
	)

	cmd.StringVar(&id, "id", "", "store id")
	cmd.StringVar(&name, "name", "", "store name")
	cmd.StringVar(&desc, "description", "", "store description")
	cmd.StringVar(&status, "status", "", "store status")

	if err := cmd.Parse(args); err != nil {
		return errors.Wrapf(err, "parse arguments: %v", args)
	}

	itemID := store.ItemID(id)
	item := store.Item{Name: name, Desc: desc, Status: status}

	if _, err := c.store.Update(itemID, item); err != nil {
		return errors.Wrap(err, "update item")
	}

	if err := c.store.PrintItems(); err != nil {
		return errors.Wrap(err, "print items")
	}

	return nil
}

func (c *Cli) DeleteCommand(args []string) error {
	cmd := flag.NewFlagSet("delete", flag.ExitOnError)

	var id string
	cmd.StringVar(&id, "id", "", "store id")

	if err := cmd.Parse(args); err != nil {
		return errors.Wrapf(err, "parse arguments: %v", args)
	}

	itemID := store.ItemID(id)

	if err := c.store.Delete(itemID); err != nil {
		return errors.Wrap(err, "delete item")
	}

	if err := c.store.PrintItems(); err != nil {
		return errors.Wrap(err, "print items")
	}

	return nil
}
