package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	didpb "github.com/bryk-io/did-method/proto"
	"go.bryk.io/x/did"
	"go.bryk.io/x/net/rpc"
)

func init() {
	catalog["bryk"] = &brykResolver{
		endpoint: "did.bryk.io:443",
	}
}

// "bryk" DID method
// https://github.com/bryk-io/did-method/blob/master/README.md
type brykResolver struct {
	endpoint string
}

// Resolve a specific DID instance as defined in the "bryk" method specification
func (br *brykResolver) Resolve(value string) ([]byte, error) {
	id, err := verify(value, "bryk")
	if err != nil {
		return nil, err
	}

	// Get network connection
	var opts []rpc.ClientOption
	opts = append(opts, rpc.WaitForReady())
	opts = append(opts, rpc.WithTimeout(5*time.Second))
	opts = append(opts, rpc.WithClientTLS(rpc.ClientTLSConfig{IncludeSystemCAs:true}))
	conn, err := rpc.NewClientConnection(br.endpoint, opts...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	// Retrieve element
	client := didpb.NewAgentAPIClient(conn)
	res, err := client.Retrieve(context.TODO(), &didpb.Query{Subject: id.Subject()})
	if err != nil {
		return nil, err
	}

	// Decode document
	doc := &did.Document{}
	if err = json.Unmarshal(res.Source, doc); err != nil {
		return nil, fmt.Errorf("failed to decode received DID Document: %s", err)
	}
	return json.MarshalIndent(doc, "", "  ")
}
