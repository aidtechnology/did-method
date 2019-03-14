package proto

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/bryk-io/x/crypto/pow"
	"github.com/bryk-io/x/did"
	"golang.org/x/crypto/sha3"
)

const ticketDifficultyLevel = 24

// NewTicket returns a properly initialized new ticket instance
func NewTicket(contents []byte, keyID string) *Ticket {
	return &Ticket{
		Timestamp:  time.Now().Unix(),
		Content:    contents,
		NonceValue: 0,
		KeyId:      keyID,
	}
}

// ResetNonce returns the internal nonce value back to 0
func (t *Ticket) ResetNonce() {
	t.NonceValue = 0
}

// IncrementNonce will adjust the internal nonce value by 1
func (t *Ticket) IncrementNonce() {
	t.NonceValue++
}

// Nonce returns the current value set on the nonce attribute
func (t *Ticket) Nonce() int64 {
	return t.NonceValue
}

// Encode returns a deterministic binary encoding for the ticket instance using a
// concatenation of values of the form 'timestamp | nonce | content'; where both
// timestamp and nonce are individually encoded using little endian byte order
func (t *Ticket) Encode() ([]byte, error) {
	var tc []byte
	nb := bytes.NewBuffer(nil)
	tb := bytes.NewBuffer(nil)
	kb := make([]byte, hex.EncodedLen(len([]byte(t.KeyId))))
	if err := binary.Write(nb, binary.LittleEndian, t.Nonce()); err != nil {
		return nil, fmt.Errorf("failed to encode nonce value: %s", err)
	}
	if err := binary.Write(tb, binary.LittleEndian, t.GetTimestamp()); err != nil {
		return nil, fmt.Errorf("failed to encode nonce value: %s", err)
	}
	hex.Encode(kb, []byte(t.KeyId))
	tc = append(tc, tb.Bytes()...)
	tc = append(tc, nb.Bytes()...)
	tc = append(tc, kb...)
	return append(tc, t.Content...), nil
}

// Solve the ticket challenge using the proof-of-work mechanism
func (t *Ticket) Solve(ctx context.Context) (string, error) {
	res, err := pow.Solve(ctx, t, sha3.New256(), ticketDifficultyLevel)
	if err != nil {
		return "", err
	}
	return <-res, nil
}

// Verify perform all the required validations to ensure the request ticket is
// ready for further processing
// - Challenge is valid
// - Contents are a properly encoded DID instance
// - DID instance’s “method” value is set to “bryk”
// - Contents don’t include any private key, for security reasons no private keys should
//   ever be published on the network
// - Signature is valid
func (t *Ticket) Verify(k *did.PublicKey) (err error) {
	// Challenge is valid
	if !pow.Verify(t, sha3.New256(), ticketDifficultyLevel) {
		return errors.New("invalid ticket challenge")
	}

	// Contents are a properly encoded DID instance
	id, err := t.LoadDID()
	if err != nil {
		return err
	}

	// DID instance’s “method” value is set to “bryk”
	if id.Method() != "bryk" {
		return errors.New("invalid DID method")
	}

	// Verify private keys are not included
	for _, k := range id.Keys() {
		if len(k.Private) != 0 {
			return errors.New("private keys included on the DID")
		}
	}

	var key *did.PublicKey
	if k != nil {
		// Use provided key
		key = k
	} else {
		// Retrieve DID's key
		key = id.Key(t.KeyId)
	}
	if key == nil {
		return errors.New("the selected key is not available on the DID")
	}

	// Get digest
	data, err := t.Encode()
	if err != nil {
		return errors.New("failed to re-encode ticket instance")
	}
	digest := sha3.New256()
	digest.Write(data)

	// Verify signature
	if !key.Verify(digest.Sum(nil), t.Signature) {
		return errors.New("invalid ticket signature")
	}
	return
}

// LoadDID obtain the DID instance encoded in the ticket contents
func (t *Ticket) LoadDID() (*did.Identifier, error) {
	id := &did.Identifier{}
	if err := id.Decode(t.Content); err != nil {
		return nil, errors.New("invalid ticket contents")
	}
	return id, nil
}
