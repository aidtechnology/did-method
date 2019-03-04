package method

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/bryk-io/id/proto"
	"github.com/bryk-io/x/did"
	"github.com/bryk-io/x/net/rpc"
	"github.com/bryk-io/x/storage/kv"
	"google.golang.org/grpc"
)

// Handler provides the required functionality for the DID method
type Handler struct {
	db     *kv.Store
	server *rpc.Server
}

// NewHandler starts a new DID method handler instance
func NewHandler(home string) (*Handler, error) {
	h := filepath.Clean(home)
	if !dirExist(h) {
		if err := os.Mkdir(h, 0700); err != nil {
			return nil, fmt.Errorf("failed to create new home directory: %s", err)
		}
	}
	db, err := kv.Open(path.Join(h, "data"), false)
	if err != nil {
		return nil, err
	}
	return &Handler{db: db}, nil
}

// Close the instance and safely terminate any internal processing
func (h *Handler) Close() (err error) {
	err = h.db.Close()
	if h.server != nil {
		err = h.server.Stop()
	}
	return
}

// Retrieve an existing DID instance based on its subject strting
func (h *Handler) Retrieve(subject string) (*did.Identifier, error) {
	contents, err := h.db.Get([]byte(subject))
	if err != nil {
		return nil, err
	}
	d := &did.Identifier{}
	if err = d.Decode(contents); err != nil {
		return nil, err
	}
	return d, nil
}

// Process an incoming request ticket
func (h *Handler) Process(ticket *proto.Ticket) error {
	if err := ticket.Verify(); err != nil {
		return err
	}
	id, err := ticket.LoadDID()
	if err != nil {
		return err
	}
	data, err := id.Encode()
	if err != nil {
		return err
	}
	return h.db.Update(&kv.Item{
		Key:   []byte(id.Subject()),
		Value: data,
	})
}

// GetServer returns a ready-to-use DID method handler RPC server
func (h *Handler) GetServer(opts ...rpc.ServerOption) (*rpc.Server, error) {
	// Return existing server instance
	if h.server != nil {
		return h.server, nil
	}

	// Add RPC service handler
	opts = append(opts, rpc.WithService(func(s *grpc.Server) {
		proto.RegisterMethodServer(s, &rpcHandler{handler: h})
	}, proto.RegisterMethodHandlerFromEndpoint))

	// Create server instance
	var err error
	h.server, err = rpc.NewServer(opts...)
	return h.server, err
}

// GetConnection returns an RPC client connection to the method handler instance
func (h *Handler) GetConnection(opts ...rpc.ClientOption) (*grpc.ClientConn, error) {
	if h.server == nil {
		return nil, errors.New("no server initialized")
	}
	return rpc.NewClientConnection(h.server.GetEndpoint(), opts...)
}

// Verify the provided path exists and is a directory
func dirExist(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}
