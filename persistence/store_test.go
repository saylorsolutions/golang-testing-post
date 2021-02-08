package persistence

import (
	"runtime/debug"
	"testing"
)

func TestGetEntryStore(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			t.Fatalf("Panic recovered: %v\n", r)
		}
	}()

	trackerFileName = "testfile.json"

	store, err := GetEntryStore()
	if err != nil {
		t.Fatalf("Failed to get store: %v\n", err)
	}
	if store == nil {
		t.Fatalf("Store is nil")
	}
}
