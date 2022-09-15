package storage

import (
	"fmt"

	protov1 "github.com/aidtechnology/did-method/proto/did/v1"
	"go.bryk.io/pkg/errors"
)

// NotFound is returned when a resolve operation targets a non-existing DID.
type NotFound struct {
	// DID subject.
	Subject string

	// DID method.
	Method string
}

// Error returns a string representation of the `not found` instance.
func (nf *NotFound) Error() string {
	return fmt.Sprintf("no information available for subject: '%s' on method: '%s'", nf.Subject, nf.Method)
}

// Is provides the custom comparison logic for the error type.
func (nf *NotFound) Is(target error) bool {
	var e *NotFound
	if errors.As(target, &e) {
		return e.Method == nf.Method && e.Subject == nf.Subject
	}
	return false
}

// NotFoundError is a utility function returning an error instance that indicates
// a resolve request target's doesn't exist.
func NotFoundError(req *protov1.QueryRequest) error {
	return &NotFound{
		Method:  req.Method,
		Subject: req.Subject,
	}
}
