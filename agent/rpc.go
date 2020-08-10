package agent

import (
	"context"
	"encoding/json"

	protov1 "github.com/bryk-io/did-method/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Wrapper to enable RPC access to an underlying method handler instance
type rpcHandler struct {
	protov1.UnimplementedAgentAPIServer
	handler *Handler
}

func getHeaders() metadata.MD {
	return metadata.New(map[string]string{
		"x-content-type-options": "nosniff",
	})
}

func (rh *rpcHandler) Ping(ctx context.Context, _ *emptypb.Empty) (*protov1.PingResponse, error) {
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
	id, proof, err := rh.handler.Retrieve(req)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	doc, _ := json.Marshal(id.Document(true))
	pp, _ := json.Marshal(proof)
	return &protov1.QueryResponse{
		Document: doc,
		Proof:    pp,
	}, nil
}
