package agent

import (
	"context"
	"encoding/json"

	didpb "github.com/bryk-io/did-method/proto"
	"github.com/gogo/protobuf/types"
)

// Wrapper to enable RPC access to an underlying method handler instance
type rpcHandler struct {
	handler *Handler
}

func (rh *rpcHandler) Ping(ctx context.Context, _ *types.Empty) (*didpb.Pong, error) {
	return &didpb.Pong{Ok: true}, nil
}

func (rh *rpcHandler) Process(ctx context.Context, req *didpb.Request) (*didpb.ProcessResponse, error) {
	err := rh.handler.Process(req)
	return &didpb.ProcessResponse{Ok: err == nil}, err
}

func (rh *rpcHandler) Retrieve(ctx context.Context, req *didpb.Query) (*didpb.Response, error) {
	id, err := rh.handler.Retrieve(req.Subject)
	if err != nil {
		return nil, err
	}
	js, _ := json.Marshal(id.Document())
	return &didpb.Response{
		Document: populateDocument(id),
		Source:   js,
	}, nil
}
