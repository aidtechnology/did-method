// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        buf-v1.17.0
// source: did/v1/agent_api.proto

package didv1

import (
	reflect "reflect"
	sync "sync"

	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Available mutations.
type ProcessRequest_Task int32

const (
	ProcessRequest_TASK_UNSPECIFIED ProcessRequest_Task = 0
	ProcessRequest_TASK_PUBLISH     ProcessRequest_Task = 1
	ProcessRequest_TASK_DEACTIVATE  ProcessRequest_Task = 2
)

// Enum value maps for ProcessRequest_Task.
var (
	ProcessRequest_Task_name = map[int32]string{
		0: "TASK_UNSPECIFIED",
		1: "TASK_PUBLISH",
		2: "TASK_DEACTIVATE",
	}
	ProcessRequest_Task_value = map[string]int32{
		"TASK_UNSPECIFIED": 0,
		"TASK_PUBLISH":     1,
		"TASK_DEACTIVATE":  2,
	}
)

func (x ProcessRequest_Task) Enum() *ProcessRequest_Task {
	p := new(ProcessRequest_Task)
	*p = x
	return p
}

func (x ProcessRequest_Task) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProcessRequest_Task) Descriptor() protoreflect.EnumDescriptor {
	return file_did_v1_agent_api_proto_enumTypes[0].Descriptor()
}

func (ProcessRequest_Task) Type() protoreflect.EnumType {
	return &file_did_v1_agent_api_proto_enumTypes[0]
}

func (x ProcessRequest_Task) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProcessRequest_Task.Descriptor instead.
func (ProcessRequest_Task) EnumDescriptor() ([]byte, []int) {
	return file_did_v1_agent_api_proto_rawDescGZIP(), []int{2, 0}
}

// Ticket required for write operations.
type Ticket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// UNIX timestamp (in UTC) when the ticket was generated.
	// All ticket automatically expire after 5 minutes to
	// prevent replay attacks.
	Timestamp int64 `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Nonce selected to solve the PoW challenge.
	NonceValue int64 `protobuf:"varint,2,opt,name=nonce_value,json=nonceValue,proto3" json:"nonce_value,omitempty"`
	// Cryptographic key identifier. Must be a valid 'authentication' method
	// on the DID document. The key will be used to generate the DID proof
	// and to sign the ticket itself.
	KeyId string `protobuf:"bytes,3,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	// JSON encoded DID document.
	Document []byte `protobuf:"bytes,4,opt,name=document,proto3" json:"document,omitempty"`
	// JSON encoded Proof document.
	Proof []byte `protobuf:"bytes,5,opt,name=proof,proto3" json:"proof,omitempty"`
	// Digital signature for the ticket, it's calculated using the
	// PoW solution as input.
	Signature []byte `protobuf:"bytes,6,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *Ticket) Reset() {
	*x = Ticket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_did_v1_agent_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ticket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ticket) ProtoMessage() {}

func (x *Ticket) ProtoReflect() protoreflect.Message {
	mi := &file_did_v1_agent_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ticket.ProtoReflect.Descriptor instead.
func (*Ticket) Descriptor() ([]byte, []int) {
	return file_did_v1_agent_api_proto_rawDescGZIP(), []int{0}
}

func (x *Ticket) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Ticket) GetNonceValue() int64 {
	if x != nil {
		return x.NonceValue
	}
	return 0
}

func (x *Ticket) GetKeyId() string {
	if x != nil {
		return x.KeyId
	}
	return ""
}

func (x *Ticket) GetDocument() []byte {
	if x != nil {
		return x.Document
	}
	return nil
}

func (x *Ticket) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

func (x *Ticket) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

// Basic reachability test response.
type PingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Responsiveness result.
	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_did_v1_agent_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_did_v1_agent_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_did_v1_agent_api_proto_rawDescGZIP(), []int{1}
}

func (x *PingResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

// Mutation request, either to publish or deactivate a DID record.
type ProcessRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Mutation type.
	Task ProcessRequest_Task `protobuf:"varint,1,opt,name=task,proto3,enum=did.v1.ProcessRequest_Task" json:"task,omitempty"`
	// Request ticket.
	Ticket *Ticket `protobuf:"bytes,2,opt,name=ticket,proto3" json:"ticket,omitempty"`
}

func (x *ProcessRequest) Reset() {
	*x = ProcessRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_did_v1_agent_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessRequest) ProtoMessage() {}

func (x *ProcessRequest) ProtoReflect() protoreflect.Message {
	mi := &file_did_v1_agent_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessRequest.ProtoReflect.Descriptor instead.
func (*ProcessRequest) Descriptor() ([]byte, []int) {
	return file_did_v1_agent_api_proto_rawDescGZIP(), []int{2}
}

func (x *ProcessRequest) GetTask() ProcessRequest_Task {
	if x != nil {
		return x.Task
	}
	return ProcessRequest_TASK_UNSPECIFIED
}

func (x *ProcessRequest) GetTicket() *Ticket {
	if x != nil {
		return x.Ticket
	}
	return nil
}

// Mutation result.
type ProcessResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Processing result, must be 'true' if the mutation was
	// properly applied.
	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *ProcessResponse) Reset() {
	*x = ProcessResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_did_v1_agent_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessResponse) ProtoMessage() {}

func (x *ProcessResponse) ProtoReflect() protoreflect.Message {
	mi := &file_did_v1_agent_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessResponse.ProtoReflect.Descriptor instead.
func (*ProcessResponse) Descriptor() ([]byte, []int) {
	return file_did_v1_agent_api_proto_rawDescGZIP(), []int{3}
}

func (x *ProcessResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

// Queries allow to resolve a previously registered DID document.
type QueryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// DID method.
	Method string `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	// DID subject.
	Subject string `protobuf:"bytes,2,opt,name=subject,proto3" json:"subject,omitempty"`
}

func (x *QueryRequest) Reset() {
	*x = QueryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_did_v1_agent_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRequest) ProtoMessage() {}

func (x *QueryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_did_v1_agent_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRequest.ProtoReflect.Descriptor instead.
func (*QueryRequest) Descriptor() ([]byte, []int) {
	return file_did_v1_agent_api_proto_rawDescGZIP(), []int{4}
}

func (x *QueryRequest) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *QueryRequest) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

// Query response.
type QueryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// JSON-encoded DID document.
	Document []byte `protobuf:"bytes,1,opt,name=document,proto3" json:"document,omitempty"`
	// JSON-encoded DID proof.
	Proof []byte `protobuf:"bytes,2,opt,name=proof,proto3" json:"proof,omitempty"`
}

func (x *QueryResponse) Reset() {
	*x = QueryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_did_v1_agent_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryResponse) ProtoMessage() {}

func (x *QueryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_did_v1_agent_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryResponse.ProtoReflect.Descriptor instead.
func (*QueryResponse) Descriptor() ([]byte, []int) {
	return file_did_v1_agent_api_proto_rawDescGZIP(), []int{5}
}

func (x *QueryResponse) GetDocument() []byte {
	if x != nil {
		return x.Document
	}
	return nil
}

func (x *QueryResponse) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

var File_did_v1_agent_api_proto protoreflect.FileDescriptor

var file_did_v1_agent_api_proto_rawDesc = []byte{
	0x0a, 0x16, 0x64, 0x69, 0x64, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x64, 0x69, 0x64, 0x2e, 0x76, 0x31,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x28, 0x74,
	0x68, 0x69, 0x72, 0x64, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x79, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x27, 0x74, 0x68, 0x69, 0x72, 0x64, 0x5f, 0x70,
	0x61, 0x72, 0x74, 0x79, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xae, 0x01, 0x0a, 0x06, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x6f, 0x6e,
	0x63, 0x65, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a,
	0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x6b, 0x65,
	0x79, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6b, 0x65, 0x79, 0x49,
	0x64, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x08, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x70, 0x72,
	0x6f, 0x6f, 0x66, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x22, 0x1e, 0x0a, 0x0c, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f,
	0x6b, 0x22, 0xae, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x64, 0x69, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52,
	0x04, 0x74, 0x61, 0x73, 0x6b, 0x12, 0x26, 0x0a, 0x06, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x64, 0x69, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x54,
	0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x06, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x43, 0x0a,
	0x04, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x54,
	0x41, 0x53, 0x4b, 0x5f, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x53, 0x48, 0x10, 0x01, 0x12, 0x13, 0x0a,
	0x0f, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x44, 0x45, 0x41, 0x43, 0x54, 0x49, 0x56, 0x41, 0x54, 0x45,
	0x10, 0x02, 0x22, 0x21, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x02, 0x6f, 0x6b, 0x22, 0x40, 0x0a, 0x0c, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x41, 0x0a, 0x0d, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x32, 0x85, 0x02, 0x0a, 0x08, 0x41,
	0x67, 0x65, 0x6e, 0x74, 0x41, 0x50, 0x49, 0x12, 0x46, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x14, 0x2e, 0x64, 0x69, 0x64, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x10, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x0a, 0x12, 0x08, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x69, 0x6e, 0x67, 0x12,
	0x52, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x12, 0x16, 0x2e, 0x64, 0x69, 0x64,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x17, 0x2e, 0x64, 0x69, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x10, 0x3a, 0x01, 0x2a, 0x22, 0x0b, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x5d, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x14, 0x2e, 0x64,
	0x69, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x15, 0x2e, 0x64, 0x69, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x27, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x21, 0x12, 0x1f, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x2f,
	0x7b, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x7d, 0x2f, 0x7b, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x7d, 0x42, 0xec, 0x02, 0x92, 0x41, 0x83, 0x02, 0x0a, 0x03, 0x32, 0x2e, 0x30, 0x12, 0x39,
	0x0a, 0x0f, 0x44, 0x49, 0x44, 0x20, 0x62, 0x72, 0x79, 0x6b, 0x20, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x22, 0x1f, 0x0a, 0x09, 0x42, 0x65, 0x6e, 0x20, 0x43, 0x65, 0x73, 0x73, 0x61, 0x1a, 0x12,
	0x62, 0x65, 0x6e, 0x40, 0x61, 0x69, 0x64, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x6e, 0x6f, 0x6c, 0x6f,
	0x67, 0x79, 0x32, 0x05, 0x30, 0x2e, 0x39, 0x2e, 0x32, 0x1a, 0x0b, 0x64, 0x69, 0x64, 0x2e, 0x62,
	0x72, 0x79, 0x6b, 0x2e, 0x69, 0x6f, 0x2a, 0x01, 0x02, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x32, 0x14, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a,
	0x73, 0x6f, 0x6e, 0x3a, 0x14, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5a, 0x53, 0x0a, 0x51, 0x0a, 0x06, 0x62,
	0x65, 0x61, 0x72, 0x65, 0x72, 0x12, 0x47, 0x08, 0x02, 0x12, 0x32, 0x41, 0x75, 0x74, 0x68, 0x65,
	0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x20,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x64, 0x20, 0x61, 0x73, 0x3a, 0x20, 0x27, 0x42, 0x65,
	0x61, 0x72, 0x65, 0x72, 0x20, 0x7b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x7d, 0x27, 0x1a, 0x0d, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x02, 0x62, 0x0c,
	0x0a, 0x0a, 0x0a, 0x06, 0x62, 0x65, 0x61, 0x72, 0x65, 0x72, 0x12, 0x00, 0x0a, 0x0a, 0x63, 0x6f,
	0x6d, 0x2e, 0x64, 0x69, 0x64, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x41,
	0x70, 0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x0c, 0x64, 0x69, 0x64, 0x2f, 0x76,
	0x31, 0x3b, 0x64, 0x69, 0x64, 0x76, 0x31, 0xf8, 0x01, 0x00, 0xa2, 0x02, 0x03, 0x44, 0x58, 0x58,
	0xaa, 0x02, 0x06, 0x44, 0x69, 0x64, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x06, 0x44, 0x69, 0x64, 0x5c,
	0x56, 0x31, 0xe2, 0x02, 0x12, 0x44, 0x69, 0x64, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x07, 0x44, 0x69, 0x64, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_did_v1_agent_api_proto_rawDescOnce sync.Once
	file_did_v1_agent_api_proto_rawDescData = file_did_v1_agent_api_proto_rawDesc
)

func file_did_v1_agent_api_proto_rawDescGZIP() []byte {
	file_did_v1_agent_api_proto_rawDescOnce.Do(func() {
		file_did_v1_agent_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_did_v1_agent_api_proto_rawDescData)
	})
	return file_did_v1_agent_api_proto_rawDescData
}

var file_did_v1_agent_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_did_v1_agent_api_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_did_v1_agent_api_proto_goTypes = []interface{}{
	(ProcessRequest_Task)(0), // 0: did.v1.ProcessRequest.Task
	(*Ticket)(nil),           // 1: did.v1.Ticket
	(*PingResponse)(nil),     // 2: did.v1.PingResponse
	(*ProcessRequest)(nil),   // 3: did.v1.ProcessRequest
	(*ProcessResponse)(nil),  // 4: did.v1.ProcessResponse
	(*QueryRequest)(nil),     // 5: did.v1.QueryRequest
	(*QueryResponse)(nil),    // 6: did.v1.QueryResponse
	(*emptypb.Empty)(nil),    // 7: google.protobuf.Empty
}
var file_did_v1_agent_api_proto_depIdxs = []int32{
	0, // 0: did.v1.ProcessRequest.task:type_name -> did.v1.ProcessRequest.Task
	1, // 1: did.v1.ProcessRequest.ticket:type_name -> did.v1.Ticket
	7, // 2: did.v1.AgentAPI.Ping:input_type -> google.protobuf.Empty
	3, // 3: did.v1.AgentAPI.Process:input_type -> did.v1.ProcessRequest
	5, // 4: did.v1.AgentAPI.Query:input_type -> did.v1.QueryRequest
	2, // 5: did.v1.AgentAPI.Ping:output_type -> did.v1.PingResponse
	4, // 6: did.v1.AgentAPI.Process:output_type -> did.v1.ProcessResponse
	6, // 7: did.v1.AgentAPI.Query:output_type -> did.v1.QueryResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_did_v1_agent_api_proto_init() }
func file_did_v1_agent_api_proto_init() {
	if File_did_v1_agent_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_did_v1_agent_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ticket); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_did_v1_agent_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_did_v1_agent_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_did_v1_agent_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_did_v1_agent_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_did_v1_agent_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_did_v1_agent_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_did_v1_agent_api_proto_goTypes,
		DependencyIndexes: file_did_v1_agent_api_proto_depIdxs,
		EnumInfos:         file_did_v1_agent_api_proto_enumTypes,
		MessageInfos:      file_did_v1_agent_api_proto_msgTypes,
	}.Build()
	File_did_v1_agent_api_proto = out.File
	file_did_v1_agent_api_proto_rawDesc = nil
	file_did_v1_agent_api_proto_goTypes = nil
	file_did_v1_agent_api_proto_depIdxs = nil
}
