package persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var trackerFileName = ".tracker.json"

var _ EntryStore = (*jsonEntryStore)(nil)

type jsonEntryStore struct {
	entries []*TimeEntry
}

func (j *jsonEntryStore) ListEntriesToday() ([]string, error) {
	panic("implement me")
}

func (j *jsonEntryStore) SaveEntry(entry *TimeEntry) error {
	panic("implement me")
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, ErrFailRead(errors.New("unable to locate user home dir"))
	}
	storePath := path.Join(homeDir, trackerFileName)
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
