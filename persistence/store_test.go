package persistence

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"
	"time"

	testify "github.com/stretchr/testify/require"
)

func TestGetEntryStore(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			t.Fatalf("Panic recovered: %v\n", r)
		}
	}()

	trackerFileName = "testfile.json"
	assert := testify.New(t)

	store, err := GetEntryStore()
	assert.NoError(err, "Failed to get store: %v\n", err)
	assert.NotNil(store)
}

func TestFormatEntry(t *testing.T) {
	assert := testify.New(t)

	start := time.Now().Add(-5 * time.Second)
	end := time.Now()
	entry := &TimeEntry{
		Start:       start,
		End:         end,
		Description: "Some description",
	}

	output := formatEntry(entry)
	assert.Equal(fmt.Sprintf("%s - %s: Some description", start.Format("15:04:05"), end.Format("15:04:05")), output)

	assert.Panics(func() {
		formatEntry(nil)
	})
}

func TestJsonEntryStore_SaveAndRetrieve(t *testing.T) {
	assert := testify.New(t)
	trackerFileName = "testfile.json"
	defer cleanupTestFile()

	defer func() {
		if r := recover(); r != nil {
			_ = os.Remove(trackerFileName)
			debug.PrintStack()
			t.Fatalf("Panic recovered: %v\n", r)
		}
	}()

	store, err := GetEntryStore()
	assert.NoError(err, "Failed to get store: %v\n", err)
	assert.NotNil(store)

	entry := &TimeEntry{
		Start:       time.Now().Add(-5 * time.Second),
		End:         time.Now(),
		Description: "A helpful description",
	}

	err = store.SaveEntry(entry)
	assert.NoError(err)

	// Check IO operations as well
	newStore, err := GetEntryStore()
	assert.NoError(err, "Failed to get new store: %v\n", err)
	assert.NotNil(newStore)
	entries, err := newStore.ListEntriesToday()
	assert.NoError(err)
	assert.Len(entries, 1)
	assert.Equal(formatEntry(entry), entries[0])
}

func cleanupTestFile() {
	// Delete the file when we're done
	storePath, err := getStorePath()
	if err != nil {
		fmt.Printf("Unable to find tracker path: %v\n", err)
	}
	fmt.Printf("Deleting store file %s\n", storePath)
	_ = os.Remove(storePath)
}
