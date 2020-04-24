package agent

import (
	"context"
	"encoding/json"

	protov1 "github.com/bryk-io/did-method/proto/v1"
	"github.com/gogo/protobuf/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Wrapper to enable RPC access to an underlying method handler instance
type rpcHandler struct {
	handler *Handler
}

func getHeaders() metadata.MD {
	return metadata.New(map[string]string{
		"x-content-type-options": "nosniff",
	})
}

func (rh *rpcHandler) Ping(ctx context.Context, _ *types.Empty) (*protov1.PingResponse, error) {
	if err := grpc.SendHeader(ctx, getHeaders()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &protov1.PingResponse{Ok: true}, nil
}

func (rh *rpcHandler) Process(ctx context.Context, req *protov1.ProcessRequest) (*protov1.ProcessResponse, error) {
	if err := grpc.SendHeader(ctx, getHeaders()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := rh.handler.Process(req); err != nil {
		return &protov1.ProcessResponse{Ok: false}, status.Error(codes.InvalidArgument, err.Error())
	}
	return &protov1.ProcessResponse{Ok: true}, nil
}

func (rh *rpcHandler) Query(ctx context.Context, req *protov1.QueryRequest) (*protov1.QueryResponse, error) {
	if err := grpc.SendHeader(ctx, getHeaders()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, err := rh.handler.Retrieve(req)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	js, _ := json.Marshal(id.SafeDocument())
	return &protov1.QueryResponse{
		Document: js,
	}, nil
}
