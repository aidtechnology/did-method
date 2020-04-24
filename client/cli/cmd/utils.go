package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	"github.com/bryk-io/did-method/client/store"
	"github.com/bryk-io/did-method/resolver"
	"github.com/spf13/viper"
	"go.bryk.io/x/crypto/ed25519"
	"go.bryk.io/x/net/rpc"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"
)

// When reading contents from standard input a maximum of 4MB is expected
const maxPipeInputSize = 4096

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
	m, err := expand(material, 32, nil)
	if err != nil {
		return nil, err
	}
	seed := [32]byte{}
	copy(seed[:], m)
	return ed25519.FromSeed(seed[:])
}

// Accessor to the local storage handler
func getClientStore() (*store.LocalStore, error) {
	return store.NewLocalStore(viper.GetString("client.home"))
}

// Get an RPC network connection
func getClientConnection() (*grpc.ClientConn, error) {
	node := viper.GetString("client.node")
	log.Infof("establishing connection to the network with node: %s", node)
	timeout := viper.GetInt("client.timeout")
	opts := []rpc.ClientOption{
		rpc.WaitForReady(),
		rpc.WithUserAgent(fmt.Sprintf("didctl-client/%s", coreVersion)),
		rpc.WithTimeout(time.Duration(timeout) * time.Second),
	}
	if viper.GetBool("client.tls") {
		opts = append(opts, rpc.WithClientTLS(rpc.ClientTLSConfig{IncludeSystemCAs: true}))
	}
	if override := viper.GetString("client.override"); override != "" {
		opts = append(opts, rpc.WithServerNameOverride(override))
	}
	return rpc.NewClientConnection(node, opts...)
}

// Use the global resolver to obtain the DID document for the requested
// identifier.
func resolve(id string) ([]byte, error) {
	var conf []*resolver.Provider
	if err := viper.UnmarshalKey("resolver", &conf); err != nil {
		return nil, fmt.Errorf("invalid resolver configuration: %s", err)
	}
	return resolver.Get(id, conf)
}
