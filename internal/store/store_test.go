package store

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func setupTest() func() {
	// set-up code here
	err := os.Remove(ItemsFilename)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Error removing items file: %s, error: %s\n", ItemsFilename, err.Error())
		os.Exit(1)
	}
	// tear down later
	return func() {
		// tear-down code here
		err = os.Remove(ItemsFilename)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Error removing items file: %s, error: %s\n", ItemsFilename, err.Error())
			os.Exit(1)
		}
	}
}

func Test_Create_ReturnsItem(t *testing.T) {
	defer setupTest()()

	data := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}
	store := Store{
		items: data,
	}

	expectedItem := Item{
		ID: "id3", Name: "name3", Desc: "desc3", Status: "status3",
	}

	item := Item{
		Name: "name3", Desc: "desc3", Status: "status3",
	}

	actualItem, err := store.Create(item)
	actualItem.ID = expectedItem.ID

	require.NoError(t, err)
	require.Equal(t, expectedItem, actualItem)
}

func Test_Create_AddsToItems(t *testing.T) {
	defer setupTest()()

	data := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}
	store := Store{
		items: data,
	}

	item := Item{
		Name: "name3", Desc: "desc3", Status: "status3",
	}

	_, err := store.Create(item)

	require.NoError(t, err)
	require.Len(t, store.items, 3)
}

func Test_ReadAll_ReturnsItems(t *testing.T) {
	defer setupTest()()

	data := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}
	store := Store{
		items: data,
	}

	expectedItems := []Item{
		{ID: "id1", Name: "name1", Desc: "desc1", Status: "status1"},
		{ID: "id2", Name: "name2", Desc: "desc2", Status: "status2"},
	}

	actualItems, err := store.ReadAll()

	require.NoError(t, err)
	require.ElementsMatch(t, expectedItems, actualItems)
}

func Test_ReadAll_DoesNotChangeItems(t *testing.T) {
	defer setupTest()()

	data := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}
	store := Store{
		items: data,
	}

	expectedItems := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}

	_, err := store.ReadAll()

	require.NoError(t, err)
	require.Equal(t, expectedItems, store.items)
}

func Test_Read_ReturnsItem(t *testing.T) {
	defer setupTest()()

	data := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}
	store := Store{
		items: data,
	}

	expectedItem := Item{
		ID: "id2", Name: "name2", Desc: "desc2", Status: "status2",
	}

	actualItem, err := store.Read("id2")

	require.NoError(t, err)
	require.Equal(t, expectedItem, actualItem)
}

func Test_Read_DoesNotChangeItems(t *testing.T) {
	defer setupTest()()

	data := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}
	store := Store{
		items: data,
	}

	expectedItems := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}

	_, err := store.Read("id2")

	require.NoError(t, err)
	require.Equal(t, expectedItems, store.items)
}

func Test_Update_ReturnsItem(t *testing.T) {
	defer setupTest()()

	data := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}
	store := Store{
		items: data,
	}

	expectedItem := Item{
		ID: "id2", Name: "name3", Desc: "desc3", Status: "status3",
	}

	item := Item{
		Name: "name3", Desc: "desc3", Status: "status3",
	}

	actualItem, err := store.Update("id2", item)

	require.NoError(t, err)
	require.Equal(t, expectedItem, actualItem)
}

func Test_Update_UpdatesItems(t *testing.T) {
	defer setupTest()()

	data := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}
	store := Store{
		items: data,
	}

	expectedItems := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name3", "desc3", "status3"},
	}

	item := Item{
		Name: "name3", Desc: "desc3", Status: "status3",
	}

	_, err := store.Update("id2", item)

	require.NoError(t, err)
	require.Equal(t, expectedItems, store.items)
}

func Test_Delete_DoesNotReturnError(t *testing.T) {
	defer setupTest()()

	data := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}
	store := Store{
		items: data,
	}

	err := store.Delete("id2")

	require.NoError(t, err)
}

func Test_Delete_RemovesFromItems(t *testing.T) {
	defer setupTest()()

	data := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}
	store := Store{
		items: data,
	}

	expectedItems := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
	}

	err := store.Delete("id2")

	require.NoError(t, err)
	require.Equal(t, expectedItems, store.items)
}
