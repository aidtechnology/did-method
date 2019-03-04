package method

import (
	"context"
	"encoding/json"

	"github.com/bryk-io/id/proto"
	"github.com/gogo/protobuf/types"
)

// Wrapper to enable RPC access to an underlying method handler instance
type rpcHandler struct {
	handler *Handler
}

func (rh *rpcHandler) Ping(ctx context.Context, _ *types.Empty) (*proto.Pong, error) {
	return &proto.Pong{Ok: true}, nil
}

func (rh *rpcHandler) Process(ctx context.Context, ticket *proto.Ticket) (*proto.Response, error) {
	err := rh.handler.Process(ticket)
	return &proto.Response{Ok: err == nil}, err
}

func (rh *rpcHandler) Retrieve(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	id, err := rh.handler.Retrieve(req.Subject)
	if err != nil {
		return nil, err
	}
	js, err := json.MarshalIndent(id.Document(), "", "  ")
	if err != nil {
		return nil, err
	}
	return &proto.Response{
		Ok:       true,
		Contents: js,
	}, nil
}
