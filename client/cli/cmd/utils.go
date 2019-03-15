package cmd

import (
	"bufio"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bryk-io/did-method/client/store"
	"github.com/bryk-io/did-method/proto"
	"github.com/bryk-io/x/crypto/ed25519"
	"github.com/bryk-io/x/did"
	"github.com/bryk-io/x/net/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/x-cray/logrus-prefixed-formatter"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
)

// When reading contents from standard input a maximum of 4MB is expected
const maxPipeInputSize = 4096

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

// Accessor to the local storage handler
func getClientStore() (*store.LocalStore, error) {
	return store.NewLocalStore(viper.GetString("home"))
}

// Get an RPC network connection
func getClientConnection(ll *log.Logger) (*grpc.ClientConn, error) {
	node := viper.GetString("node")
	if ll != nil {
		ll.Infof("establishing connection to the network with node: %s", node)
	}
	var opts []rpc.ClientOption
	opts = append(opts, rpc.WaitForReady())
	opts = append(opts, rpc.WithUserAgent("bryk-id-client"))
	opts = append(opts, rpc.WithTimeout(5*time.Second))
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

// Read contents passed in from standard input, if any
func getPipedInput() ([]byte, error) {
	var input []byte

	// Fail to read stdin
	info, err := os.Stdin.Stat()
	if err != nil {
		return input, err
	}

	// No input passed in
	if info.Mode()&os.ModeCharDevice != 0 {
		return input, errors.New("no piped input")
	}

	// Read input
	reader := bufio.NewReader(os.Stdin)
	for {
		b, err := reader.ReadByte()
		if err != nil && err == io.EOF {
			break
		}
		input = append(input, b)
		if len(input) == maxPipeInputSize {
			break
		}
	}

	// Return provided input
	return input, nil
}

// Retrieve a DID instance from the network
func retrieveSubject(subject string, ll *log.Logger) (*did.Identifier, error) {
	// Get network connection
	conn, err := getClientConnection(ll)
	if err != nil {
		return nil, err
	}

	client := proto.NewAgentClient(conn)
	res, err := client.Retrieve(context.TODO(), &proto.Query{Subject: subject})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve DID records: %s", err)
	}
	if !res.Ok {
		return nil, errors.New("no information available for the provided DID")
	}

	// Decode contents
	ll.Debug("decoding contents")
	doc := &did.Document{}
	if err = doc.Decode(res.Contents); err != nil {
		return nil, fmt.Errorf("failed to decode received DID Document: %s", err)
	}
	return did.FromDocument(doc)
}
