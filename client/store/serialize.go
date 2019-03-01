package store

import (
	"bytes"

	"github.com/vmihailenco/msgpack"
)

// Decode message pack binary data to a given 'dest' element
func decodeMsg(data []byte, dest interface{}) error {
	dec := msgpack.NewDecoder(bytes.NewReader(data))
	return dec.Decode(dest)
}

// Encode a given source element to message pack binary data
func encodeMsg(src interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf).SortMapKeys(true)
	err := enc.Encode(src)
	return buf.Bytes(), err
}
