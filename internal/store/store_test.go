package store

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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

func Test_ParallelTests(t *testing.T) {
	t.Cleanup(setupTest())

	data := map[ItemID]Item{
		"readID":   {"readID", "readName", "readDesc", "readStatus"},
		"updateID": {"updateID", "updateName", "updateDesc", "updateStatus"},
		"deleteID": {"deleteID", "deleteName", "deleteDesc", "deleteStatus"},
	}
	store := Store{
		items: data,
	}

	t.Run("ParallelTests", func(t *testing.T) {
		t.Run("Create", func(t *testing.T) {
			t.Parallel()
			item := Item{
				Name:   "createID",
				Desc:   "createDesc",
				Status: "createStatus",
			}
			createdItem, err := store.Create(item)
			assert.NoError(t, err)
			assert.Contains(t, store.items, ItemID(createdItem.ID))
		})
		t.Run("ReadAll", func(t *testing.T) {
			t.Parallel()
			_, err := store.ReadAll()
			assert.NoError(t, err)
		})
		t.Run("Read", func(t *testing.T) {
			t.Parallel()
			_, err := store.Read("readID")
			assert.NoError(t, err)
		})
		t.Run("Update", func(t *testing.T) {
			t.Parallel()
			item := Item{
				ID:     "updateID",
				Name:   "newUpdateID",
				Desc:   "newUpdateDesc",
				Status: "newUpdateStatus",
			}
			updatedItem, err := store.Update("updateID", item)
			assert.NoError(t, err)
			assert.Equal(t, item, updatedItem)

		})
		t.Run("Delete", func(t *testing.T) {
			t.Parallel()
			err := store.Delete("deleteID")
			assert.NoError(t, err)
			assert.NotContains(t, store.items, "deleteID")
		})
	})
}

func Test_Create_ReturnsItem(t *testing.T) {
	t.Cleanup(setupTest())

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

	assert.NoError(t, err)
	assert.Equal(t, expectedItem, actualItem)
}

func Test_Create_AddsToItems(t *testing.T) {
	t.Cleanup(setupTest())

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

	assert.NoError(t, err)
	assert.Len(t, store.items, 3)
}

func Test_ReadAll_ReturnsItems(t *testing.T) {
	t.Cleanup(setupTest())

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

	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedItems, actualItems)
}

func Test_ReadAll_DoesNotChangeItems(t *testing.T) {
	t.Cleanup(setupTest())

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

	assert.NoError(t, err)
	assert.Equal(t, expectedItems, store.items)
}

func Test_Read_ReturnsItem(t *testing.T) {
	t.Cleanup(setupTest())

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

	assert.NoError(t, err)
	assert.Equal(t, expectedItem, actualItem)
}

func Test_Read_DoesNotChangeItems(t *testing.T) {
	t.Cleanup(setupTest())

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

	assert.NoError(t, err)
	assert.Equal(t, expectedItems, store.items)
}

func Test_Update_ReturnsItem(t *testing.T) {
	t.Cleanup(setupTest())

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

	assert.NoError(t, err)
	assert.Equal(t, expectedItem, actualItem)
}

func Test_Update_UpdatesItems(t *testing.T) {
	t.Cleanup(setupTest())

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

	assert.NoError(t, err)
	assert.Equal(t, expectedItems, store.items)
}

func Test_Delete_DoesNotReturnError(t *testing.T) {
	t.Cleanup(setupTest())

	data := map[ItemID]Item{
		"id1": {"id1", "name1", "desc1", "status1"},
		"id2": {"id2", "name2", "desc2", "status2"},
	}
	store := Store{
		items: data,
	}

	err := store.Delete("id2")

	assert.NoError(t, err)
}

func Test_Delete_RemovesFromItems(t *testing.T) {
	t.Cleanup(setupTest())

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

	assert.NoError(t, err)
	assert.Equal(t, expectedItems, store.items)
}
