package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.bryk.io/x/ccg/did"
)

// LocalStore provides a filesystem-backed store
type LocalStore struct {
	home string
}

// NewLocalStore returns a local store handler
func NewLocalStore(home string) (*LocalStore, error) {
	h := filepath.Clean(home)
	if !dirExist(h) {
		if err := os.Mkdir(h, 0700); err != nil {
			return nil, fmt.Errorf("failed to create new home directory: %s", err)
		}
	}
	return &LocalStore{home: h}, nil
}

// Save add a new entry to the store
func (ls *LocalStore) Save(name string, id *did.Identifier) error {
	data, err := json.Marshal(id.Document(false))
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(ls.home, name), data, 0600)
}

// Get an existing entry based on its reference name
func (ls *LocalStore) Get(name string) (*did.Identifier, error) {
	data, err := ioutil.ReadFile(filepath.Clean(filepath.Join(ls.home, name)))
	if err != nil {
		return nil, err
	}
	doc := &did.Document{}
	if err := json.Unmarshal(data, doc); err != nil {
		return nil, err
	}
	return did.FromDocument(doc)
}

// List currently registered entries
func (ls *LocalStore) List() map[string]*did.Identifier {
	// nolint: prealloc
	var list = make(map[string]*did.Identifier)
	_ = filepath.Walk(ls.home, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		id, err := ls.Get(info.Name())
		if err == nil {
			list[info.Name()] = id
		}
		return nil
	})
	return list
}

// Update the contents of an existing entry
func (ls *LocalStore) Update(name string, id *did.Identifier) error {
	return ls.Save(name, id)
}

// Delete a previously stored entry
func (ls *LocalStore) Delete(name string) error {
	return os.Remove(filepath.Join(ls.home, name))
}

// Verify the provided path exists and is a directory
func dirExist(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}
