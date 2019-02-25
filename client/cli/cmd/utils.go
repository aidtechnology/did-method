package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/bryk-io/x/crypto/ed25519"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"
	"golang.org/x/crypto/ssh/terminal"
)

// Helper method to securely read data from stdin
func secureAsk(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	return terminal.ReadPassword(0)
}

// Securely expand the provided secret material
func expand(secret []byte, size int, info []byte) ([]byte, error) {
	salt := make([]byte, sha256.Size)
	buf := make([]byte, size)
	h := hkdf.New(sha3.New256, secret, salt[:], info)
	if _, err := io.ReadFull(h, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// Restore key pair from the provided material
func keyFromMaterial(material []byte) (*ed25519.KeyPair, error) {
	m, err := expand(material, ed25519.SeedSize, nil)
	if err != nil {
		return nil, err
	}
	seed := [ed25519.SeedSize]byte{}
	copy(seed[:], m)
	return ed25519.Restore(seed)
}
