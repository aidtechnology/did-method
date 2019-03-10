package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	"github.com/bryk-io/did-method/client/store"
	"github.com/bryk-io/x/crypto/ed25519"
	"github.com/bryk-io/x/net/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/x-cray/logrus-prefixed-formatter"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
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

func getClientStore() (*store.LocalStore, error) {
	return store.NewLocalStore(viper.GetString("home"))
}

func getClientConnection(ll *log.Logger) (*grpc.ClientConn, error) {
	node := viper.GetString("node")
	if ll != nil {
		ll.Infof("establishing connection to the network with node: %s", node)
	}
	var opts []rpc.ClientOption
	opts = append(opts, rpc.WaitForReady())
	opts = append(opts, rpc.WithUserAgent("bryk-id-client"))
	opts = append(opts, rpc.WithTimeout(5 * time.Second))
	return rpc.NewClientConnection(node, opts...)
}

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
