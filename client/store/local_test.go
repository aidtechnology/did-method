package store

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/bryk-io/x/did"
)

func TestLocalStore(t *testing.T) {
	home, _ := ioutil.TempDir("", "testing_")
	st, err := NewLocalStore(home)
	if err != nil {
		t.Fatal("failed to create local store:", err)
	}
	defer func() {
		_ = st.Close()
	}()

	// Sample entry
	id, _ := did.NewIdentifierWithMode("bryk", "", did.ModeUUID)
	if err := id.AddNewKey("master", did.KeyTypeEd, did.EncodingHex); err != nil {
		t.Error(err)
	}
	if err := id.AddAuthenticationKey("master"); err != nil {
		t.Error(err)
	}
	if err := id.AddProof("master", "sample.acme.com"); err != nil {
		t.Error(err)
	}
	contents, _ := id.Encode()
	entry := &Entry{
		Name:     id.Subject(),
		Recovery: "passphrase",
		Contents: contents,
	}

	// Save
	if err := st.Save(id.Subject(), entry); err != nil {
		t.Fatal("failed to save entry:", err)
	}

	// GET
	if res := st.Get("invalid-name"); res != nil {
		t.Fatal("failed to return nil for invalid entry names")
	}
	r2 := st.Get(id.Subject())
	id2 := &did.Identifier{}
	if err = id2.Decode(r2.Contents); err != nil {
		t.Fatal("failed to decode identifier from entry:", err)
	}
	if !bytes.Equal(id.Key("master").Private, id2.Key("master").Private) {
		t.Fatal("failed to decode master key")
	}

	// List
	if len(st.List()) != 1 {
		t.Fatal("invalid entry count")
	}

	// Update
	if err := id.AddNewKey("iadb-provider", did.KeyTypeEd, did.EncodingHex); err != nil {
		t.Error(err)
	}
	if err := id.AddProof("master", "sample.acme.com"); err != nil {
		t.Error(err)
	}
	newContents, _ := id.Encode()
	if err = st.Update("invalid-entry", newContents); err == nil {
		t.Fatal("failed to catch invalid entry to update")
	}
	if err = st.Update(id.Subject(), newContents); err != nil {
		t.Fatal("failed to update entry:", err)
	}

	// Validate updated data
	r3 := st.Get(id.Subject())
	id3 := &did.Identifier{}
	if err = id3.Decode(r3.Contents); err != nil {
		t.Fatal("failed to decode identifier from entry:", err)
	}
	if len(id3.Keys()) != 2 {
		t.Fatal("invalid keys count for updated entry")
	}
}
