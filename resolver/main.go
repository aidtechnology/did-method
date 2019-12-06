package resolver

import (
	"errors"
	"fmt"

	"go.bryk.io/x/did"
)

type methodHandler interface {
	Resolve(val string) ([]byte, error)
}

var catalog = make(map[string]methodHandler)

// Get will attempt to retrieve the DID document corresponding to the provided DID
func Get(value string) ([]byte, error) {
	// Verify the provided value is a valid DID string
	id, err := did.Parse(value)
	if err != nil {
		return nil, err
	}

	// Get method handler
	mh, ok := catalog[id.Method()]
	if !ok {
		return nil, errors.New("unsupported did method")
	}

	// Return handler result
	return mh.Resolve(value)
}

// Common verification steps
func verify(value string, method string) (*did.Identifier, error) {
	// Verify provided value
	id, err := did.Parse(value)
	if err != nil {
		return nil, err
	}

	// Validate method value
	if id.Method() != method {
		return nil, fmt.Errorf("invalid method value: %s", id.Method())
	}

	return id, nil
}
