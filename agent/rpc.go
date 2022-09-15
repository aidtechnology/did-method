package agent

import (
	"context"
	"encoding/json"

	"github.com/aidtechnology/did-method/agent/storage"
	protov1 "github.com/aidtechnology/did-method/proto/did/v1"
	"go.bryk.io/pkg/errors"
	"go.bryk.io/pkg/otel"
	otelcodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Wrapper to enable RPC access to an underlying method handler instance.
type rpcHandler struct {
	protov1.UnimplementedAgentAPIServer
	handler *Handler
}

func (rh *rpcHandler) Ping(ctx context.Context, _ *emptypb.Empty) (*protov1.PingResponse, error) {
	// Get parent span reference
	parent := rh.handler.oop.SpanFromContext(ctx)

	// Track operation
	task := rh.handler.oop.Start(parent.Context(), "Ping", otel.WithSpanKind(otel.SpanKindServer))
	defer task.End()
	return &protov1.PingResponse{Ok: true}, nil
}

func (rh *rpcHandler) Process(ctx context.Context, req *protov1.ProcessRequest) (*protov1.ProcessResponse, error) {
	// Get parent span reference
	parent := rh.handler.oop.SpanFromContext(ctx)

	// Track operation
	task := rh.handler.oop.Start(parent.Context(), "rpc.Process", otel.WithSpanKind(otel.SpanKindServer))
	defer task.End()

	// Process request
	if err := rh.handler.Process(task.Context(), req); err != nil {
		task.SetStatus(otelcodes.Error, err.Error())
		return &protov1.ProcessResponse{Ok: false}, status.Error(codes.InvalidArgument, err.Error())
	}
	return &protov1.ProcessResponse{Ok: true}, nil
}

func (rh *rpcHandler) Query(ctx context.Context, req *protov1.QueryRequest) (*protov1.QueryResponse, error) {
	// Get parent span reference
	parent := rh.handler.oop.SpanFromContext(ctx)

	// Track operation
	task := rh.handler.oop.Start(
		parent.Context(),
		"rpc.Query",
		otel.WithSpanKind(otel.SpanKindServer),
		otel.WithSpanAttributes(otel.Attributes{
			"method": req.Method,
		}))
	defer task.End()

	// Process request
	id, proof, err := rh.handler.Retrieve(task.Context(), req)
	if err != nil {
		if !errors.Is(err, storage.NotFoundError(req)) {
			task.SetStatus(otelcodes.Error, err.Error())
		}
		return nil, status.Error(codes.NotFound, err.Error())
	}
	doc, _ := json.Marshal(id.Document(true))
	pp, _ := json.Marshal(proof)
	return &protov1.QueryResponse{
		Document: doc,
		Proof:    pp,
	}, nil
}
