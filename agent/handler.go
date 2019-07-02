package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/bryk-io/did-method/proto"
	"github.com/bryk-io/x/did"
	"github.com/bryk-io/x/net/rpc"
	"github.com/bryk-io/x/storage/kv"
	log "github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"google.golang.org/grpc"
)

// Handler provides the required functionality for the DID method
type Handler struct {
	db     *kv.Store
	server *rpc.Server
	output *log.Logger
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
	return &Handler{
		db:     db,
		output: getLogger(),
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
	contents, err := h.db.Get([]byte(subject))
	if err != nil {
		return nil, errors.New("no information available for the subject")
	}
	doc := &did.Document{}
	if err = doc.Decode(contents); err != nil {
		return nil, err
	}
	return did.FromDocument(doc)
}

// Process an incoming request ticket
func (h *Handler) Process(req *proto.Request) error {
	// Validate ticket
	if err := req.Ticket.Verify(nil); err != nil {
		h.output.WithField("error", err.Error()).Error("invalid ticket")
		return err
	}

	// Load DID document
	id, err := req.Ticket.LoadDID()
	if err != nil {
		h.output.WithField("error", err.Error()).Error("invalid DID contents")
		return err
	}

	// Update operations require another validation step using the original record
	isUpdate := false
	if r, err := h.db.Get([]byte(id.Subject())); err == nil && len(r) != 0 {
		orig, err := h.Retrieve(id.Subject())
		if err != nil {
			return fmt.Errorf("failed to recover original record for update: %s", err)
		}
		if err := req.Ticket.Verify(orig.Key(req.Ticket.KeyId)); err != nil {
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
	case proto.Request_PUBLISH:
		return h.db.Update(&kv.Item{
			Key:   []byte(id.Subject()),
			Value: req.Ticket.Content,
		})
	case proto.Request_DEACTIVATE:
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
		GatewaySetup: proto.RegisterAgentHandlerFromEndpoint,
		Setup: func(s *grpc.Server) {
			proto.RegisterAgentServer(s, &rpcHandler{handler: h})
		},
	}))

	// Custom HTTP handler method for data retrieval, the response should be the JSON-LD encoded
	// document of the requested DID instance
	opts = append(opts, rpc.WithHandlerFunc("/v1/retrieve", h.queryHTTP))

	// Create server instance
	srv, err := rpc.NewServer(opts...)
	if err != nil {
		return nil, err
	}

	h.server = srv
	return h.server, nil
}

// GetConnection returns an RPC client connection to the method handler instance
func (h *Handler) GetConnection(opts ...rpc.ClientOption) (*grpc.ClientConn, error) {
	if h.server == nil {
		return nil, errors.New("no server initialized")
	}
	return rpc.NewClientConnection(h.server.GetEndpoint(), opts...)
}

// Log will print a message to the handler's output
func (h *Handler) Log(message string) {
	h.output.Info(message)
}

// Handle data queries done via HTTP. Will return the pretty formatted JSON-LD document
// of the DID instance, if available. A simple error message otherwise.
func (h *Handler) queryHTTP(writer http.ResponseWriter, request *http.Request) {
	eh := map[string]interface{}{"ok": false}
	writer.Header().Set("Content-Type", "application/json")

	// Retrieve entry
	id, err := h.Retrieve(request.URL.Query().Get("subject"))
	if err != nil {
		writer.WriteHeader(400)
		eh["error"] = err.Error()
		_ = json.NewEncoder(writer).Encode(eh)
		return
	}

	// Prepare output
	output, err := json.MarshalIndent(id.Document(), "", "  ")
	if err != nil {
		writer.WriteHeader(400)
		eh["error"] = err.Error()
		_ = json.NewEncoder(writer).Encode(eh)
		return
	}

	// Send response
	writer.WriteHeader(200)
	_, _ = fmt.Fprintf(writer, "%s", output)
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
