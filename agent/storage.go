package agent

import (
	protov1 "github.com/aidtechnology/did-method/proto/did/v1"
	"go.bryk.io/pkg/did"
)

// Storage defines an abstract component that provides and manage
// persistent data requirements for DID documents.
type Storage interface {
	// Open will prepare the instance for usage.
	Open(info string) error

	// Close the storage instance, free resources and finish processing.
	Close() error

	// Description returns a brief information summary for the storage instance.
	Description() string

	// Exists will check if a record exists for the specified DID.
	Exists(id *did.Identifier) bool

	// Get will return a previously stored DID instance.
	Get(req *protov1.QueryRequest) (*did.Identifier, *did.ProofLD, error)

	// Save or update the record for the given DID instance.
	Save(id *did.Identifier, proof *did.ProofLD) error

	// Delete any existing records for the given DID instance.
	Delete(id *did.Identifier) error
}
