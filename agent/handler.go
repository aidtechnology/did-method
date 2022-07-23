package agent

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aidtechnology/did-method/info"
	protov1 "github.com/aidtechnology/did-method/proto/did/v1"
	"go.bryk.io/pkg/did"
	xlog "go.bryk.io/pkg/log"
	"go.bryk.io/pkg/net/rpc"
	"go.bryk.io/pkg/otel"
	"google.golang.org/grpc"
)

// Handler provides the required functionality for the DID method.
type Handler struct {
	oop        *otel.Operator
	methods    []string
	store      Storage
	difficulty uint
}

// NewHandler starts a new DID method handler instance.
func NewHandler(methods []string, difficulty uint, store Storage, oop *otel.Operator) (*Handler, error) {
	return &Handler{
		oop:        oop,
		store:      store,
		methods:    methods,
		difficulty: difficulty,
	}, nil
}

// Close the instance and safely terminate any internal processing.
func (h *Handler) Close() error {
	h.oop.Info("closing agent handler")
	return h.store.Close()
}

// Retrieve an existing DID instance based on its subject string.
func (h *Handler) Retrieve(ctx context.Context, req *protov1.QueryRequest) (*did.Identifier, *did.ProofLD, error) {
	// Track operation
	task := h.oop.Start(
		ctx,
		"handler.Retrieve",
		otel.WithSpanKind(otel.SpanKindServer),
		otel.WithSpanAttributes(otel.Attributes{"method": req.Method}))
	defer task.End()

	// Verify method is supported
	if !h.isSupported(req.Method) {
		err := errors.New("unsupported method")
		task.Error(xlog.Error, err, nil)
		return nil, nil, err
	}

	// Retrieve document from storage
	task.Event("database read")
	id, proof, err := h.store.Get(req)
	if err != nil {
		task.Error(xlog.Error, err, nil)
		return nil, nil, err
	}
	return id, proof, nil
}

// Process an incoming request ticket.
func (h *Handler) Process(ctx context.Context, req *protov1.ProcessRequest) error {
	// Track operation
	task := h.oop.Start(
		ctx,
		"handler.Process",
		otel.WithSpanKind(otel.SpanKindServer),
		otel.WithSpanAttributes(otel.Attributes{"task": req.Task.String()}))
	defer task.End()

	// Validate ticket
	if err := req.Ticket.Verify(h.difficulty); err != nil {
		task.Error(xlog.Error, err, nil)
		return err
	}

	// Load DID document and proof
	id, err := req.Ticket.GetDID()
	if err != nil {
		task.Error(xlog.Error, err, nil)
		return err
	}
	proof, err := req.Ticket.GetProofLD()
	if err != nil {
		task.Error(xlog.Error, err, nil)
		return err
	}

	// Verify method is supported
	if !h.isSupported(id.Method()) {
		err := errors.New("unsupported method")
		task.Error(xlog.Error, err, otel.Attributes{
			"method": id.Method(),
		})
		return err
	}

	// Update operations require another validation step using the original record
	isUpdate := h.store.Exists(id)
	if isUpdate {
		if err := req.Ticket.Verify(h.difficulty); err != nil {
			task.Error(xlog.Error, err, nil)
			return err
		}
	}

	fields := otel.Attributes{
		"subject": id.Subject(),
		"update":  isUpdate,
		"task":    req.Task.String(),
	}
	switch req.Task {
	case protov1.ProcessRequest_TASK_PUBLISH:
		task.Event("database save", fields)
		err = h.store.Save(id, proof)
	case protov1.ProcessRequest_TASK_DEACTIVATE:
		task.Event("database delete", fields)
		err = h.store.Delete(id)
	default:
		return errors.New("invalid request task")
	}
	if err != nil {
		task.Error(xlog.Error, err, fields)
	}
	return err
}

// ServerSetup perform all initialization requirements for the
// handler instance to be exposed through the provided gRPC server.
func (h *Handler) ServerSetup(srv *grpc.Server) {
	protov1.RegisterAgentAPIServer(srv, &rpcHandler{handler: h})
}

// GatewaySetup return the HTTP setup required to expose the handler
// instance via HTTP.
func (h *Handler) GatewaySetup() rpc.GatewayRegister {
	return protov1.RegisterAgentAPIHandler
}

// CustomGatewayOptions returns additional settings required when exposing
// the handler instance via HTTP.
func (h *Handler) CustomGatewayOptions() []rpc.GatewayOption {
	return []rpc.GatewayOption{
		rpc.WithSpanFormatter(spanNameFormatter()),
		rpc.WithInterceptor(h.queryResponseFilter()),
	}
}

// QueryResponseFilter provides custom encoding of HTTP query results.
func (h *Handler) queryResponseFilter() rpc.GatewayInterceptor {
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
		id, proof, err := h.Retrieve(req.Context(), rr)
		if err != nil {
			response, _ = json.MarshalIndent(map[string]string{"error": err.Error()}, "", "  ")
		} else {
			response, _ = json.MarshalIndent(map[string]interface{}{
				"document": id.Document(true),
				"proof":    proof,
			}, "", "  ")
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

// Verify a specific method is supported.
func (h *Handler) isSupported(method string) bool {
	for _, m := range h.methods {
		if method == m {
			return true
		}
	}
	return false
}

// SpanNameFormatter determines how transactions are reported to observability
// services.
func spanNameFormatter() func(r *http.Request) string {
	return func(r *http.Request) string {
		if strings.HasPrefix(r.URL.Path, "/v1/retrieve") {
			return fmt.Sprintf("%s %s", r.Method, "/v1/retrieve/{method}/{subject}")
		}
		return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
	}
}
