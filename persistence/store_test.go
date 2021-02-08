package persistence

import (
	"runtime/debug"
	"testing"

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
