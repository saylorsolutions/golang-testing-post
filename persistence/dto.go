package persistence

import "time"

// TimeEntry represents an entry to be tracked and stored in the EntryStore
type TimeEntry struct {
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Description string    `json:"description"`
}
