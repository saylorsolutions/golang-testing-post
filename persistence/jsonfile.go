package persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

var trackerFileName = ".tracker.json"

var _ EntryStore = (*jsonEntryStore)(nil)

type jsonEntryStore struct {
	entries []*TimeEntry
}

func (j *jsonEntryStore) ListEntriesToday() ([]string, error) {
	var entries []string
	now := time.Now()
	todayThreshold := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for _, e := range j.entries {
		if e.Start.After(todayThreshold) || e.End.After(todayThreshold) {
			entries = append(entries, formatEntry(e))
		}
	}
	return entries, nil
}

func formatEntry(entry *TimeEntry) string {
	return fmt.Sprintf("%s - %s: %s", entry.Start.Format("15:04:05"), entry.End.Format("15:04:05"), entry.Description)
}

func (j *jsonEntryStore) SaveEntry(entry *TimeEntry) error {
	err := validateSaveEntry(entry)
	if err != nil {
		return err
	}
	j.entries = append(j.entries, entry)
	data, err := json.Marshal(j.entries)
	if err != nil {
		return ErrFailWrite(fmt.Errorf("failed to marshal entries to JSON: %v\n", err))
	}
	storePath, err := getStorePath()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(storePath, data, 0644)
	if err != nil {
		return ErrFailWrite(fmt.Errorf("failed to write entries to file: %v\n", err))
	}
	return nil
}

func validateSaveEntry(entry *TimeEntry) error {
	if entry == nil {
		return ErrValidation(errors.New("nil entry"))
	}
	if entry.Start == zero {
		return ErrValidation(errors.New("missing start time"))
	}
	if entry.End == zero {
		return ErrValidation(errors.New("missing end time"))
	}
	if entry.End.Before(entry.Start) {
		return ErrValidation(fmt.Errorf("entry start time is before end time: %s", formatEntry(entry)))
	}
	if entry.Description == "" {
		return ErrValidation(errors.New("empty entry description"))
	}
	return nil
}

func GetEntryStore() (store EntryStore, err error) {
	data, err := loadStoreFile()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return &jsonEntryStore{}, nil
	}
	var entries []*TimeEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, ErrFailRead(fmt.Errorf("incompatible or corrupted store file: %v\n", err))
	}
	return &jsonEntryStore{entries: entries}, nil
}

func loadStoreFile() ([]byte, error) {
	storePath, err := getStorePath()
	if err != nil {
		return nil, err
	}
	storeFile, err := os.OpenFile(storePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, ErrFailRead(fmt.Errorf("unable to open JSON store: %v\n", err))
	}
	defer storeFile.Close()

	data, err := ioutil.ReadAll(storeFile)
	if err != nil {
		return nil, ErrFailRead(fmt.Errorf("unable to read store file: %v\n", err))
	}
	return data, err
}

func getStorePath() (storePath string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", ErrFailRead(errors.New("unable to locate user home dir"))
	}
	storePath = path.Join(homeDir, trackerFileName)
	return storePath, nil
}
