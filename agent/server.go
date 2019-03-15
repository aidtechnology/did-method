package agent

import (
	"context"

	"github.com/bryk-io/did-method/proto"
	"github.com/gogo/protobuf/types"
)

// Wrapper to enable RPC access to an underlying method handler instance
type rpcHandler struct {
	handler *Handler
}

func (rh *rpcHandler) Ping(ctx context.Context, _ *types.Empty) (*proto.Pong, error) {
	return &proto.Pong{Ok: true}, nil
}

func (rh *rpcHandler) Process(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	err := rh.handler.Process(req)
	return &proto.Response{Ok: err == nil}, err
}

func (rh *rpcHandler) Retrieve(ctx context.Context, req *proto.Query) (*proto.Response, error) {
	id, err := rh.handler.Retrieve(req.Subject)
	if err != nil {
		return nil, err
	}
	data, err := id.Document().Encode()
	if err != nil {
		return nil, err
	}
	return &proto.Response{
		Ok:       true,
		Contents: data,
	}, nil
}
