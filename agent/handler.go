package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	didpb "github.com/bryk-io/did-method/proto"
	log "github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"go.bryk.io/x/did"
	"go.bryk.io/x/net/rpc"
	"go.bryk.io/x/storage/kv"
	"google.golang.org/grpc"
)

// Handler provides the required functionality for the DID method
type Handler struct {
	db         *kv.Store
	server     *rpc.Server
	output     *log.Logger
	difficulty uint
}

// NewHandler starts a new DID method handler instance
func NewHandler(home string, difficulty uint) (*Handler, error) {
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
	return &Handler{
		db:         db,
		output:     getLogger(),
		difficulty: difficulty,
	}, nil
}

// Close the instance and safely terminate any internal processing
func (h *Handler) Close() (err error) {
	h.output.Info("closing agent handler")
	err = h.db.Close()
	if h.server != nil {
		err = h.server.Stop()
	}
	return
}

// Retrieve an existing DID instance based on its subject string
func (h *Handler) Retrieve(subject string) (*did.Identifier, error) {
	h.output.WithField("subject", subject).Debug("retrieve request")
	contents, err := h.db.Read([]byte(subject))
	if err != nil {
		return nil, errors.New("no information available for subject")
	}
	doc := &did.Document{}
	if err = json.Unmarshal(contents, doc); err != nil {
		return nil, err
	}
	return did.FromDocument(doc)
}

// Process an incoming request ticket
func (h *Handler) Process(req *didpb.Request) error {
	// Empty request
	if req == nil {
		return errors.New("empty request")
	}

	// Validate ticket
	if err := req.Ticket.Verify(nil, h.difficulty); err != nil {
		h.output.WithField("error", err.Error()).Error("invalid ticket")
		return err
	}

	// Load DID document
	id, err := req.Ticket.GetDID()
	if err != nil {
		h.output.WithField("error", err.Error()).Error("invalid DID contents")
		return err
	}

	// Update operations require another validation step using the original record
	isUpdate := false
	if r, err := h.db.Read([]byte(id.Subject())); err == nil && len(r) != 0 {
		orig, err := h.Retrieve(id.Subject())
		if err != nil {
			return fmt.Errorf("failed to recover original record for update: %s", err)
		}
		if err := req.Ticket.Verify(orig.Key(req.Ticket.KeyId), h.difficulty); err != nil {
			h.output.WithField("error", err.Error()).Error("invalid ticket")
			return err
		}
		isUpdate = true
	}

	h.output.WithFields(log.Fields{
		"subject": id.Subject(),
		"update":  isUpdate,
		"task":    req.Task,
	}).Debug("write operation")
	switch req.Task {
	case didpb.Request_PUBLISH:
		return h.db.Update([]byte(id.Subject()), req.Ticket.Content)
	case didpb.Request_DEACTIVATE:
		return h.db.Delete([]byte(id.Subject()))
	default:
		return errors.New("invalid request task")
	}
}

// GetServer returns a ready-to-use DID method handler RPC server
func (h *Handler) GetServer(opts ...rpc.ServerOption) (*rpc.Server, error) {
	// Return existing server instance
	if h.server != nil {
		return h.server, nil
	}

	// Add RPC service handler
	opts = append(opts, rpc.WithService(&rpc.Service{
		GatewaySetup: didpb.RegisterAgentAPIHandlerFromEndpoint,
		Setup: func(s *grpc.Server) {
			didpb.RegisterAgentAPIServer(s, &rpcHandler{handler: h})
		},
	}))

	// Create server instance
	srv, err := rpc.NewServer(opts...)
	if err != nil {
		return nil, err
	}

	h.server = srv
	return h.server, nil
}

// Log will print a message to the handler's output
func (h *Handler) Log(message string) {
	h.output.Info(message)
}

// Verify the provided path exists and is a directory
func dirExist(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func getLogger() *log.Logger {
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
