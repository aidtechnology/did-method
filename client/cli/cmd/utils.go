package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/bryk-io/did-method/client/store"
	didpb "github.com/bryk-io/did-method/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"go.bryk.io/x/crypto/ed25519"
	"go.bryk.io/x/did"
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
	return store.NewLocalStore(viper.GetString("home"))
}

// Get an RPC network connection
func getClientConnection(ll *log.Logger) (*grpc.ClientConn, error) {
	node := viper.GetString("client.node")
	if ll != nil {
		ll.Infof("establishing connection to the network with node: %s", node)
	}
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

// Output handler
func getLogger() *log.Logger {
	// Set formatter
	output := log.New()
	formatter := &prefixed.TextFormatter{}
	formatter.FullTimestamp = true
	formatter.TimestampFormat = time.StampMilli
	formatter.SetColorScheme(&prefixed.ColorScheme{
		DebugLevelStyle: "black",
		TimestampStyle:  "white+h",
	})
	output.Formatter = formatter
	output.SetLevel(log.DebugLevel)
	return output
}

// Retrieve a DID instance from the network
func retrieveSubject(subject string, ll *log.Logger) (*did.Identifier, error) {
	// Get network connection
	conn, err := getClientConnection(ll)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	client := didpb.NewAgentAPIClient(conn)
	res, err := client.Retrieve(context.TODO(), &didpb.Query{Subject: subject})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve DID records: %s", err)
	}

	// Decode contents
	ll.Debug("decoding contents")
	doc := &did.Document{}
	if err = json.Unmarshal(res.Source, doc); err != nil {
		return nil, fmt.Errorf("failed to decode received DID Document: %s", err)
	}
	return did.FromDocument(doc)
}
