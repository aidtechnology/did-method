package agent

import (
	"context"
	"encoding/json"

	protov1 "github.com/aidtechnology/did-method/proto/did/v1"
	"go.bryk.io/pkg/otel"
	otelApi "go.bryk.io/pkg/otel/api"
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
	return &protov1.PingResponse{Ok: true}, nil
}

func (rh *rpcHandler) Process(ctx context.Context, req *protov1.ProcessRequest) (res *protov1.ProcessResponse, err error) { // nolint: lll
	// Track operation
	task := otelApi.Start(ctx, "rpc.Process", otelApi.WithSpanKind(otelApi.SpanKindServer))
	defer task.End(nil)

	// Process request
	if err = rh.handler.Process(task.Context(), req); err != nil {
		res.Ok = false
		task.End(err)
		err = status.Error(codes.InvalidArgument, err.Error())
		return
	}
	res.Ok = true
	return
}

func (rh *rpcHandler) Query(ctx context.Context, req *protov1.QueryRequest) (res *protov1.QueryResponse, err error) {
	// Track operation
	task := otelApi.Start(
		ctx,
		"rpc.Query",
		otelApi.WithSpanKind(otelApi.SpanKindServer),
		otelApi.WithAttributes(otel.Attributes{
			"method": req.Method,
		}))
	defer task.End(nil)

	// Process request
	id, proof, err := rh.handler.Retrieve(task.Context(), req)
	if err != nil {
		task.End(err)
		err = status.Error(codes.NotFound, err.Error())
		return
	}
	res.Document, _ = json.Marshal(id.Document(true))
	res.Proof, _ = json.Marshal(proof)
	return
}
