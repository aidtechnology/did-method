package agent

import (
	"context"
	"encoding/json"

	didpb "github.com/bryk-io/did-method/proto"
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

func (rh *rpcHandler) Ping(ctx context.Context, _ *types.Empty) (*didpb.Pong, error) {
	if err := grpc.SendHeader(ctx, getHeaders()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &didpb.Pong{Ok: true}, nil
}

func (rh *rpcHandler) Process(ctx context.Context, req *didpb.Request) (*didpb.ProcessResponse, error) {
	if err := grpc.SendHeader(ctx, getHeaders()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := rh.handler.Process(req); err != nil {
		return &didpb.ProcessResponse{Ok: false}, status.Error(codes.InvalidArgument, err.Error())
	}
	return &didpb.ProcessResponse{Ok: true}, nil
}

func (rh *rpcHandler) Retrieve(ctx context.Context, req *didpb.Query) (*didpb.Response, error) {
	if err := grpc.SendHeader(ctx, getHeaders()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, err := rh.handler.Retrieve(req.Subject)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	js, _ := json.Marshal(id.Document())
	return &didpb.Response{
		Document: populateDocument(id),
		Source:   js,
	}, nil
}
