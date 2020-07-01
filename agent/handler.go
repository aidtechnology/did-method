package agent

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bryk-io/did-method/info"
	protov1 "github.com/bryk-io/did-method/proto/v1"
	"go.bryk.io/x/ccg/did"
	xlog "go.bryk.io/x/log"
	"go.bryk.io/x/net/rpc"
	"google.golang.org/grpc"
)

// Handler provides the required functionality for the DID method
type Handler struct {
	methods    []string
	store      Storage
	log        xlog.Logger
	difficulty uint
}

// NewHandler starts a new DID method handler instance
func NewHandler(logger xlog.Logger, methods []string, difficulty uint, store Storage) (*Handler, error) {
	return &Handler{
		log:        logger,
		store:      store,
		methods:    methods,
		difficulty: difficulty,
	}, nil
}

// Close the instance and safely terminate any internal processing
func (h *Handler) Close() error {
	h.log.Info("closing agent handler")
	return h.store.Close()
}

// Retrieve an existing DID instance based on its subject string
func (h *Handler) Retrieve(req *protov1.QueryRequest) (*did.Identifier, error) {
	logFields := xlog.Fields{
		"method":  req.Method,
		"subject": req.Subject,
	}
	h.log.WithFields(logFields).Debug("retrieve request")

	// Verify method is supported
	if !h.isSupported(req.Method) {
		h.log.WithFields(logFields).Warning("non supported method")
		return nil, errors.New("non supported method")
	}

	// Retrieve document from storage
	id, err := h.store.Get(req)
	if err != nil {
		h.log.WithFields(logFields).Warning(err.Error())
		return nil, err
	}
	return id, nil
}

// Process an incoming request ticket
func (h *Handler) Process(req *protov1.ProcessRequest) error {
	// Empty request
	if req == nil {
		return errors.New("empty request")
	}

	// Validate ticket
	if err := req.Ticket.Verify(nil, h.difficulty); err != nil {
		h.log.WithFields(xlog.Fields{"error": err.Error()}).Error("invalid ticket")
		return err
	}

	// Load DID document
	id, err := req.Ticket.GetDID()
	if err != nil {
		h.log.WithFields(xlog.Fields{"error": err.Error()}).Error("invalid DID contents")
		return err
	}

	// Verify method is supported
	if !h.isSupported(id.Method()) {
		h.log.WithFields(xlog.Fields{"method": id.Method()}).Warning("non supported method")
		return errors.New("non supported method")
	}

	// Update operations require another validation step using the original record
	isUpdate := h.store.Exists(id)
	if isUpdate {
		if err := req.Ticket.Verify(id.Key(req.Ticket.KeyId), h.difficulty); err != nil {
			h.log.WithFields(xlog.Fields{"error": err.Error()}).Error("invalid ticket")
			return err
		}
	}

	h.log.WithFields(xlog.Fields{
		"subject": id.Subject(),
		"update":  isUpdate,
		"task":    req.Task,
	}).Debug("write operation")
	switch req.Task {
	case protov1.ProcessRequest_TASK_PUBLISH:
		err = h.store.Save(id)
	case protov1.ProcessRequest_TASK_DEACTIVATE:
		err = h.store.Delete(id)
	default:
		return errors.New("invalid request task")
	}
	return err
}

// ServiceDefinition allows the handler instance to be exposed using a RPC server
func (h *Handler) ServiceDefinition() *rpc.Service {
	return &rpc.Service{
		GatewaySetup: protov1.RegisterAgentAPIHandlerFromEndpoint,
		ServerSetup: func(s *grpc.Server) {
			protov1.RegisterAgentAPIServer(s, &rpcHandler{handler: h})
		},
	}
}

// QueryResponseFilter provides custom encoding of HTTP query results.
func (h *Handler) QueryResponseFilter() rpc.HTTPGatewayFilter {
	return func(res http.ResponseWriter, req *http.Request) error {
		// Filter query requests
		if !strings.HasPrefix(req.URL.Path, "/v1/retrieve/") {
			return nil
		}
		seg := strings.Split(strings.TrimPrefix(req.URL.Path, "/v1/retrieve/"), "/")
		if len(seg) != 2 {
			return nil
		}

		// Submit query
		var (
			status   = http.StatusNotFound
			response []byte
		)
		rr := &protov1.QueryRequest{
			Method:  seg[0],
			Subject: seg[1],
		}
		id, err := h.Retrieve(rr)
		if err != nil {
			response, _ = json.MarshalIndent(map[string]string{"error": err.Error()}, "", "  ")
		} else {
			response, _ = json.MarshalIndent(id.Document(true), "", "  ")
			status = http.StatusOK
			res.Header().Set("Etag", fmt.Sprintf("W/%x", sha256.Sum256(response)))
		}

		// Return result
		res.Header().Set("content-type", "application/json")
		res.Header().Set("x-content-type-options", "nosniff")
		res.Header().Set("x-didctl-build-code", info.BuildCode)
		res.Header().Set("x-didctl-build-timestamp", info.BuildTimestamp)
		res.Header().Set("x-didctl-version", info.CoreVersion)
		res.WriteHeader(status)
		_, _ = res.Write(response)
		return errors.New("prevent further processing")
	}
}

// Verify a specific method is supported
func (h *Handler) isSupported(method string) bool {
	for _, m := range h.methods {
		if method == m {
			return true
		}
	}
	return false
}
