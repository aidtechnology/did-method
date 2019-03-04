package store

// Entry represents an existing record on the store
type Entry struct {
	// Reference name for the DID instance
	Name string

	// Recovery method specified
	Recovery string

	// DID instance encoded contents
	Contents []byte
}

// Encode the entry to its binary messagepack representation
func (e *Entry) Encode() ([]byte, error) {
	return encodeMsg(e)
}

// Decode will restore the entry instance from its existing binary messagepack representation
func (e *Entry) Decode(val []byte) error {
	return decodeMsg(val, e)
}
