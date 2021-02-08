package persistence

// EntryStore provides simple CRUD operations.
type EntryStore interface {
	// ListEntriesToday retrieves the entries from today and returns their descriptions prefixed with start and end time.
	// An ErrFailRead is returned if the operation fails for any reason.
	ListEntriesToday() ([]string, error)

	// SaveEntry persists the given entry to te data store. An ErrFailWrite is returned if the operation fails.
	SaveEntry(entry *TimeEntry) error
}

// ErrFailRead indicates that a read operation failed.
type ErrFailRead error

// ErrFailWrite indicates that a write operation failed.
type ErrFailWrite error
