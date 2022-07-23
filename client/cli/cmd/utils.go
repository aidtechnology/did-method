package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/aidtechnology/did-method/client/store"
	"github.com/aidtechnology/did-method/info"
	"github.com/aidtechnology/did-method/resolver"
	"github.com/spf13/viper"
	"go.bryk.io/pkg/crypto/ed25519"
	xlog "go.bryk.io/pkg/log"
	"go.bryk.io/pkg/net/rpc"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"
)

// When reading contents from standard input a maximum of 4MB is expected.
const maxPipeInputSize = 4096

// Securely expand the provided secret material.
func expand(secret []byte, size int, info []byte) ([]byte, error) {
	salt := make([]byte, sha256.Size)
	buf := make([]byte, size)
	h := hkdf.New(sha3.New256, secret, salt[:], info)
	if _, err := io.ReadFull(h, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// Restore key pair from the provided material.
func keyFromMaterial(material []byte) (*ed25519.KeyPair, error) {
	m, err := expand(material, 32, nil)
	if err != nil {
		return nil, err
	}
	seed := [32]byte{}
	copy(seed[:], m)
	return ed25519.FromSeed(seed[:])
}

// Accessor to the local storage handler.
func getClientStore() (*store.LocalStore, error) {
	home := viper.GetString("client.home")
	log.WithField("home", home).Debug("local store handler")
	return store.NewLocalStore(home)
}

// Get an RPC network connection.
func getClientConnection(op ...rpc.ClientOption) (*grpc.ClientConn, error) {
	node := viper.GetString("client.node")
	timeout := viper.GetInt("client.timeout")
	agentV := fmt.Sprintf("didctl-client/%s", info.CoreVersion)
	log.WithFields(xlog.Fields{
		"node":       node,
		"timeout":    timeout,
		"user-agent": agentV,
	}).Info("establishing connection to the network")

	// establish new connection
	opts, err := conf.ClientRPC()
	if err != nil {
		return nil, err
	}
	opts = append(opts, op...)
	return rpc.NewClientConnection(node, opts...)
}

// Use the global resolver to obtain the DID document for the requested
// identifier.
func resolve(id string) ([]byte, error) {
	return resolver.Get(id, conf.Resolver)
}
