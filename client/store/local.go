package store

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/bryk-io/x/storage/kv"
)

// LocalStore provides a filesystem-backed store
type LocalStore struct {
	db *kv.Store
}

// NewLocalStore returns a local store handler
func NewLocalStore(home string) (*LocalStore, error) {
	h := filepath.Clean(home)
	if !dirExist(h) {
		if err := os.Mkdir(h, 0700); err != nil {
			return nil, fmt.Errorf("failed to create new home directory: %s", err)
		}
	}
	db, err := kv.Open(path.Join(h, "data"), false)
	if err != nil {
		return nil, err
	}
	return &LocalStore{db: db}, nil
}

// Save add a new entry to the store
func (ls *LocalStore) Save(name string, record *Entry) error {
	contents, err := record.Encode()
	if err != nil {
		return err
	}
	return ls.db.Create([]byte(name), contents)
}

// Get an existing entry based on its reference name
func (ls *LocalStore) Get(name string) *Entry {
	contents, err := ls.db.Read([]byte(name))
	if err != nil {
		return nil
	}
	rec := &Entry{}
	if err = rec.Decode(contents); err != nil {
		return nil
	}
	return rec
}

// List currently registered entries
func (ls *LocalStore) List() []*Entry {
	var err error
	var list []*Entry
	cursor := ls.db.Iterate(context.TODO(), &kv.CursorOptions{KeysOnly: false})
	for i := range cursor {
		rec := &Entry{}
		if err = rec.Decode(i.Value); err != nil {
			continue
		}
		list = append(list, rec)
	}
	return list
}

// Update the contents of an existing entry
func (ls *LocalStore) Update(name string, contents []byte) error {
	rec := ls.Get(name)
	if rec == nil {
		return fmt.Errorf("no entry for the id: %s", name)
	}
	rec.Contents = contents
	data, err := rec.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode entry for storage: %s", err)
	}
	return ls.db.Update([]byte(name), data)
}

// Delete a previously stored entry
func (ls *LocalStore) Delete(name string) error {
	return ls.db.Delete([]byte(name))
}

// Close the store instance and free resources
func (ls *LocalStore) Close() error {
	return ls.db.Close()
}

// Verify the provided path exists and is a directory
func dirExist(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}
