// Code generated by protoc-gen-go-drpc. DO NOT EDIT.
// protoc-gen-go-drpc version: v0.0.30
// source: did/v1/agent_api.proto

package didv1

import (
	context "context"
	errors "errors"

	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	drpc "storj.io/drpc"
	drpcerr "storj.io/drpc/drpcerr"
)

type drpcEncoding_File_did_v1_agent_api_proto struct{}

func (drpcEncoding_File_did_v1_agent_api_proto) Marshal(msg drpc.Message) ([]byte, error) {
	return proto.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_did_v1_agent_api_proto) MarshalAppend(buf []byte, msg drpc.Message) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(buf, msg.(proto.Message))
}

func (drpcEncoding_File_did_v1_agent_api_proto) Unmarshal(buf []byte, msg drpc.Message) error {
	return proto.Unmarshal(buf, msg.(proto.Message))
}

func (drpcEncoding_File_did_v1_agent_api_proto) JSONMarshal(msg drpc.Message) ([]byte, error) {
	return protojson.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_did_v1_agent_api_proto) JSONUnmarshal(buf []byte, msg drpc.Message) error {
	return protojson.Unmarshal(buf, msg.(proto.Message))
}

type DRPCAgentAPIClient interface {
	DRPCConn() drpc.Conn

	Ping(ctx context.Context, in *emptypb.Empty) (*PingResponse, error)
	Process(ctx context.Context, in *ProcessRequest) (*ProcessResponse, error)
	Query(ctx context.Context, in *QueryRequest) (*QueryResponse, error)
}

type drpcAgentAPIClient struct {
	cc drpc.Conn
}

func NewDRPCAgentAPIClient(cc drpc.Conn) DRPCAgentAPIClient {
	return &drpcAgentAPIClient{cc}
}

func (c *drpcAgentAPIClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcAgentAPIClient) Ping(ctx context.Context, in *emptypb.Empty) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/did.v1.AgentAPI/Ping", drpcEncoding_File_did_v1_agent_api_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcAgentAPIClient) Process(ctx context.Context, in *ProcessRequest) (*ProcessResponse, error) {
	out := new(ProcessResponse)
	err := c.cc.Invoke(ctx, "/did.v1.AgentAPI/Process", drpcEncoding_File_did_v1_agent_api_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcAgentAPIClient) Query(ctx context.Context, in *QueryRequest) (*QueryResponse, error) {
	out := new(QueryResponse)
	err := c.cc.Invoke(ctx, "/did.v1.AgentAPI/Query", drpcEncoding_File_did_v1_agent_api_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type DRPCAgentAPIServer interface {
	Ping(context.Context, *emptypb.Empty) (*PingResponse, error)
	Process(context.Context, *ProcessRequest) (*ProcessResponse, error)
	Query(context.Context, *QueryRequest) (*QueryResponse, error)
}

type DRPCAgentAPIUnimplementedServer struct{}

func (s *DRPCAgentAPIUnimplementedServer) Ping(context.Context, *emptypb.Empty) (*PingResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCAgentAPIUnimplementedServer) Process(context.Context, *ProcessRequest) (*ProcessResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCAgentAPIUnimplementedServer) Query(context.Context, *QueryRequest) (*QueryResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

type DRPCAgentAPIDescription struct{}

func (DRPCAgentAPIDescription) NumMethods() int { return 3 }

func (DRPCAgentAPIDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/did.v1.AgentAPI/Ping", drpcEncoding_File_did_v1_agent_api_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentAPIServer).
					Ping(
						ctx,
						in1.(*emptypb.Empty),
					)
			}, DRPCAgentAPIServer.Ping, true
	case 1:
		return "/did.v1.AgentAPI/Process", drpcEncoding_File_did_v1_agent_api_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentAPIServer).
					Process(
						ctx,
						in1.(*ProcessRequest),
					)
			}, DRPCAgentAPIServer.Process, true
	case 2:
		return "/did.v1.AgentAPI/Query", drpcEncoding_File_did_v1_agent_api_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCAgentAPIServer).
					Query(
						ctx,
						in1.(*QueryRequest),
					)
			}, DRPCAgentAPIServer.Query, true
	default:
		return "", nil, nil, nil, false
	}
}

func DRPCRegisterAgentAPI(mux drpc.Mux, impl DRPCAgentAPIServer) error {
	return mux.Register(impl, DRPCAgentAPIDescription{})
}

type DRPCAgentAPI_PingStream interface {
	drpc.Stream
	SendAndClose(*PingResponse) error
}

type drpcAgentAPI_PingStream struct {
	drpc.Stream
}

func (x *drpcAgentAPI_PingStream) SendAndClose(m *PingResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_did_v1_agent_api_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCAgentAPI_ProcessStream interface {
	drpc.Stream
	SendAndClose(*ProcessResponse) error
}

type drpcAgentAPI_ProcessStream struct {
	drpc.Stream
}

func (x *drpcAgentAPI_ProcessStream) SendAndClose(m *ProcessResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_did_v1_agent_api_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCAgentAPI_QueryStream interface {
	drpc.Stream
	SendAndClose(*QueryResponse) error
}

type drpcAgentAPI_QueryStream struct {
	drpc.Stream
}

func (x *drpcAgentAPI_QueryStream) SendAndClose(m *QueryResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_did_v1_agent_api_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}
