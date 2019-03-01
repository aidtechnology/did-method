package store

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// LocalStore provides a filesystem-backed store
type LocalStore struct {
	home string
}

// NewLocalStore returns a local store handler
func NewLocalStore(home string) (*LocalStore, error) {
	if !dirExist(home) {
		if err := os.Mkdir(home, 0700); err != nil {
			return nil, fmt.Errorf("failed to create new home directory: %s", err)
		}
	}

	return &LocalStore{
		home: home,
	}, nil
}

// Save add a new entry to the store
func (ls *LocalStore) Save(id string, record *Entry) error {
	if exist(path.Join(ls.home, id)) {
		return errors.New("duplicated entry")
	}
	return ls.save(id, record)
}

// List currently registered entries
func (ls *LocalStore) List() []*Entry {
	var list []*Entry
	_ = filepath.Walk(ls.home, func(f string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			entry, err := loadEntry(filepath.Clean(f))
			if err != nil {
				return err
			}
			list = append(list, entry)
		}
		return nil
	})
	return list
}

// Get an existing entry based on its reference name
func (ls *LocalStore) Get(name string) *Entry {
	for _, e := range ls.List() {
		if e.Name == name {
			return e
		}
	}
	return nil
}

// Update the contents of an existing entry
func (ls *LocalStore) Update(id string, contents []byte) error {
	if !exist(path.Join(ls.home, id)) {
		return errors.New("duplicated entry")
	}
	entry, err := loadEntry(path.Join(ls.home, filepath.Clean(id)))
	if err != nil {
		return err
	}
	entry.Contents = contents
	return ls.save(id, entry)
}

// Close the store instance and free resources
func (ls *LocalStore) Close() error {
	return nil
}

// Store a record
func (ls *LocalStore) save(id string, record *Entry) error {
	contents, err := record.Encode()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(ls.home, id), contents, 0600)
}

// Verify the provided path is either a file or directory that exists
func exist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

// Verify the provided path exists and is a directory
func dirExist(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Get an entry from an existing file
func loadEntry(f string) (*Entry, error) {
	contents, err := ioutil.ReadFile(filepath.Clean(f))
	if err != nil {
		return nil, err
	}
	entry := &Entry{}
	if err := entry.Decode(contents); err != nil {
		return nil, err
	}
	return entry, nil
}
