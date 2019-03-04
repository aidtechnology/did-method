package proto

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/bryk-io/x/crypto/pow"
	"github.com/bryk-io/x/did"
	e "golang.org/x/crypto/ed25519"
)

const ticketDifficultyLevel = 24

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
	if err := binary.Write(nb, binary.LittleEndian, t.Nonce()); err != nil {
		return nil, fmt.Errorf("failed to encode nonce value: %s", err)
	}
	if err := binary.Write(tb, binary.LittleEndian, t.GetTimestamp()); err != nil {
		return nil, fmt.Errorf("failed to encode nonce value: %s", err)
	}
	tc = append(tc, tb.Bytes()...)
	tc = append(tc, nb.Bytes()...)
	return append(tc, t.Content...), nil
}

// LoadDID obtain the DID instance encoded in the ticket contents
func (t *Ticket) LoadDID() (*did.Identifier, error) {
	id := &did.Identifier{}
	if err := id.Decode(t.Content); err != nil {
		return nil, errors.New("invalid ticket contents")
	}
	return id, nil
}

// Solve the ticket challenge using the proof-of-work mechanism
func (t *Ticket) Solve(ctx context.Context) (string, error) {
	res, err := pow.Solve(ctx, t, sha256.New(), ticketDifficultyLevel)
	if err != nil {
		return "", err
	}
	return <-res, nil
}

// Verify perform all the required validations to ensure the request ticket is
// ready for further processing.
// - Challenge is valid
// - Contents are a properly encoded DID instance
// - DID instance’s “method” value is set to “bryk”
// - Contents don’t include any private key, for security reasons no private keys should
//   ever be published on the network
// - Signature is valid
func (t *Ticket) Verify() (err error) {
	// Challenge is valid
	if !pow.Verify(t, sha256.New(), ticketDifficultyLevel) {
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

	// Retrieve DID's master key
	key := id.Key("master")
	if key == nil {
		return errors.New("no master key available on the DID")
	}

	// Decode public key
	var pubBytes []byte
	if key.ValueHex != "" {
		pubBytes, err = hex.DecodeString(key.ValueHex)
		if err != nil {
			return errors.New("invalid key hex encoding")
		}
	}
	if key.ValueBase64 != "" {
		pubBytes, err = base64.StdEncoding.DecodeString(key.ValueBase64)
		if err != nil {
			return errors.New("invalid key base64 encoding")
		}
	}

	// Get digest
	data, err := t.Encode()
	if err != nil {
		return errors.New("failed to re-encode ticket instance")
	}
	digest := sha256.New()
	digest.Write(data)

	// Verify signature
	pub := e.PublicKey(pubBytes)
	if !e.Verify(pub, digest.Sum(nil), t.Signature) {
		return errors.New("invalid ticket signature")
	}
	return
}
