// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: pkg/core/v1/internals/service_invocation.proto

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package internals

import (
	common "github.com/bhojpur/api/pkg/core/v1/common"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Actor represents compute processing actor using actor_type and actor_id
type Actor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required. The type of actor.
	ActorType string `protobuf:"bytes,1,opt,name=actor_type,json=actorType,proto3" json:"actor_type,omitempty"`
	// Required. The ID of actor type (actor_type)
	ActorId string `protobuf:"bytes,2,opt,name=actor_id,json=actorId,proto3" json:"actor_id,omitempty"`
}

func (x *Actor) Reset() {
	*x = Actor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_core_v1_internals_service_invocation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Actor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Actor) ProtoMessage() {}

func (x *Actor) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_core_v1_internals_service_invocation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Actor.ProtoReflect.Descriptor instead.
func (*Actor) Descriptor() ([]byte, []int) {
	return file_pkg_core_v1_internals_service_invocation_proto_rawDescGZIP(), []int{0}
}

func (x *Actor) GetActorType() string {
	if x != nil {
		return x.ActorType
	}
	return ""
}

func (x *Actor) GetActorId() string {
	if x != nil {
		return x.ActorId
	}
	return ""
}

// InternalInvokeRequest is the message to transfer caller's data to callee for
// service invocation. This includes callee's app id and caller's request data.
type InternalInvokeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required. The version of Bhojpur Application runtime API.
	Ver APIVersion `protobuf:"varint,1,opt,name=ver,proto3,enum=v1.internals.APIVersion" json:"ver,omitempty"`
	// Required. metadata holds caller's HTTP headers or gRPC metadata.
	Metadata map[string]*ListStringValue `protobuf:"bytes,2,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Required. message including caller's invocation request.
	Message *common.InvokeRequest `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	// Actor type and id. This field is used only for actor service invocation.
	Actor *Actor `protobuf:"bytes,4,opt,name=actor,proto3" json:"actor,omitempty"`
}

func (x *InternalInvokeRequest) Reset() {
	*x = InternalInvokeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_core_v1_internals_service_invocation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InternalInvokeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InternalInvokeRequest) ProtoMessage() {}

func (x *InternalInvokeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_core_v1_internals_service_invocation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InternalInvokeRequest.ProtoReflect.Descriptor instead.
func (*InternalInvokeRequest) Descriptor() ([]byte, []int) {
	return file_pkg_core_v1_internals_service_invocation_proto_rawDescGZIP(), []int{1}
}

func (x *InternalInvokeRequest) GetVer() APIVersion {
	if x != nil {
		return x.Ver
	}
	return APIVersion_APIVERSION_UNSPECIFIED
}

func (x *InternalInvokeRequest) GetMetadata() map[string]*ListStringValue {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *InternalInvokeRequest) GetMessage() *common.InvokeRequest {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *InternalInvokeRequest) GetActor() *Actor {
	if x != nil {
		return x.Actor
	}
	return nil
}

// InternalInvokeResponse is the message to transfer callee's response to caller
// for service invocation.
type InternalInvokeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required. HTTP/gRPC service status.
	Status *Status `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	// Required. The application callback response headers.
	Headers map[string]*ListStringValue `protobuf:"bytes,2,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Application callback response trailers. It will be used only for gRPC type
	// application callback
	Trailers map[string]*ListStringValue `protobuf:"bytes,3,rep,name=trailers,proto3" json:"trailers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Callee's invocation response message.
	Message *common.InvokeResponse `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *InternalInvokeResponse) Reset() {
	*x = InternalInvokeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_core_v1_internals_service_invocation_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InternalInvokeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InternalInvokeResponse) ProtoMessage() {}

func (x *InternalInvokeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_core_v1_internals_service_invocation_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InternalInvokeResponse.ProtoReflect.Descriptor instead.
func (*InternalInvokeResponse) Descriptor() ([]byte, []int) {
	return file_pkg_core_v1_internals_service_invocation_proto_rawDescGZIP(), []int{2}
}

func (x *InternalInvokeResponse) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *InternalInvokeResponse) GetHeaders() map[string]*ListStringValue {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *InternalInvokeResponse) GetTrailers() map[string]*ListStringValue {
	if x != nil {
		return x.Trailers
	}
	return nil
}

func (x *InternalInvokeResponse) GetMessage() *common.InvokeResponse {
	if x != nil {
		return x.Message
	}
	return nil
}

// ListStringValue represents string value array
type ListStringValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The array of string.
	Values []string `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *ListStringValue) Reset() {
	*x = ListStringValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_core_v1_internals_service_invocation_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListStringValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListStringValue) ProtoMessage() {}

func (x *ListStringValue) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_core_v1_internals_service_invocation_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListStringValue.ProtoReflect.Descriptor instead.
func (*ListStringValue) Descriptor() ([]byte, []int) {
	return file_pkg_core_v1_internals_service_invocation_proto_rawDescGZIP(), []int{3}
}

func (x *ListStringValue) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_pkg_core_v1_internals_service_invocation_proto protoreflect.FileDescriptor

var file_pkg_core_v1_internals_service_invocation_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f,
	0x69, 0x6e, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0c, 0x76, 0x31, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x73, 0x1a, 0x1f,
	0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x26, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x22, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2f, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x41, 0x0a, 0x05, 0x41,
	0x63, 0x74, 0x6f, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x22, 0xcd,
	0x02, 0x0a, 0x15, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x76, 0x6f, 0x6b,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x03, 0x76, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x76, 0x31, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x73, 0x2e, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x03, 0x76, 0x65, 0x72, 0x12, 0x4d, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x76, 0x31, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49, 0x6e,
	0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x32, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x49, 0x6e, 0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x61, 0x63, 0x74, 0x6f, 0x72,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x76, 0x31, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2e, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x05, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x1a, 0x5a, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x33, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x76, 0x31, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xcf,
	0x03, 0x0a, 0x16, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x76, 0x6f, 0x6b,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x76, 0x31, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x4b, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x76, 0x31, 0x2e, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x49, 0x6e, 0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x12, 0x4e, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x76, 0x31, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49, 0x6e,
	0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x54, 0x72, 0x61,
	0x69, 0x6c, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x74, 0x72, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x73, 0x12, 0x33, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x49, 0x6e, 0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x59, 0x0a, 0x0c, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x33, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x76, 0x31, 0x2e,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x1a, 0x5a, 0x0a, 0x0d, 0x54, 0x72, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x33, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x76, 0x31, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x29, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x32, 0xc7, 0x01, 0x0a, 0x11,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x58, 0x0a, 0x09, 0x43, 0x61, 0x6c, 0x6c, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x23,
	0x2e, 0x76, 0x31, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2e, 0x49, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x76, 0x31, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x73, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x76, 0x6f, 0x6b,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x58, 0x0a, 0x09, 0x43,
	0x61, 0x6c, 0x6c, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x12, 0x23, 0x2e, 0x76, 0x31, 0x2e, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x49, 0x6e, 0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e,
	0x76, 0x31, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x73, 0x2e, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x68, 0x6f, 0x6a, 0x70, 0x75, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x73, 0x3b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_core_v1_internals_service_invocation_proto_rawDescOnce sync.Once
	file_pkg_core_v1_internals_service_invocation_proto_rawDescData = file_pkg_core_v1_internals_service_invocation_proto_rawDesc
)

func file_pkg_core_v1_internals_service_invocation_proto_rawDescGZIP() []byte {
	file_pkg_core_v1_internals_service_invocation_proto_rawDescOnce.Do(func() {
		file_pkg_core_v1_internals_service_invocation_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_core_v1_internals_service_invocation_proto_rawDescData)
	})
	return file_pkg_core_v1_internals_service_invocation_proto_rawDescData
}

var file_pkg_core_v1_internals_service_invocation_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pkg_core_v1_internals_service_invocation_proto_goTypes = []interface{}{
	(*Actor)(nil),                  // 0: v1.internals.Actor
	(*InternalInvokeRequest)(nil),  // 1: v1.internals.InternalInvokeRequest
	(*InternalInvokeResponse)(nil), // 2: v1.internals.InternalInvokeResponse
	(*ListStringValue)(nil),        // 3: v1.internals.ListStringValue
	nil,                            // 4: v1.internals.InternalInvokeRequest.MetadataEntry
	nil,                            // 5: v1.internals.InternalInvokeResponse.HeadersEntry
	nil,                            // 6: v1.internals.InternalInvokeResponse.TrailersEntry
	(APIVersion)(0),                // 7: v1.internals.APIVersion
	(*common.InvokeRequest)(nil),   // 8: v1.common.InvokeRequest
	(*Status)(nil),                 // 9: v1.internals.Status
	(*common.InvokeResponse)(nil),  // 10: v1.common.InvokeResponse
}
var file_pkg_core_v1_internals_service_invocation_proto_depIdxs = []int32{
	7,  // 0: v1.internals.InternalInvokeRequest.ver:type_name -> v1.internals.APIVersion
	4,  // 1: v1.internals.InternalInvokeRequest.metadata:type_name -> v1.internals.InternalInvokeRequest.MetadataEntry
	8,  // 2: v1.internals.InternalInvokeRequest.message:type_name -> v1.common.InvokeRequest
	0,  // 3: v1.internals.InternalInvokeRequest.actor:type_name -> v1.internals.Actor
	9,  // 4: v1.internals.InternalInvokeResponse.status:type_name -> v1.internals.Status
	5,  // 5: v1.internals.InternalInvokeResponse.headers:type_name -> v1.internals.InternalInvokeResponse.HeadersEntry
	6,  // 6: v1.internals.InternalInvokeResponse.trailers:type_name -> v1.internals.InternalInvokeResponse.TrailersEntry
	10, // 7: v1.internals.InternalInvokeResponse.message:type_name -> v1.common.InvokeResponse
	3,  // 8: v1.internals.InternalInvokeRequest.MetadataEntry.value:type_name -> v1.internals.ListStringValue
	3,  // 9: v1.internals.InternalInvokeResponse.HeadersEntry.value:type_name -> v1.internals.ListStringValue
	3,  // 10: v1.internals.InternalInvokeResponse.TrailersEntry.value:type_name -> v1.internals.ListStringValue
	1,  // 11: v1.internals.ServiceInvocation.CallActor:input_type -> v1.internals.InternalInvokeRequest
	1,  // 12: v1.internals.ServiceInvocation.CallLocal:input_type -> v1.internals.InternalInvokeRequest
	2,  // 13: v1.internals.ServiceInvocation.CallActor:output_type -> v1.internals.InternalInvokeResponse
	2,  // 14: v1.internals.ServiceInvocation.CallLocal:output_type -> v1.internals.InternalInvokeResponse
	13, // [13:15] is the sub-list for method output_type
	11, // [11:13] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_pkg_core_v1_internals_service_invocation_proto_init() }
func file_pkg_core_v1_internals_service_invocation_proto_init() {
	if File_pkg_core_v1_internals_service_invocation_proto != nil {
		return
	}
	file_pkg_core_v1_internals_apiversion_proto_init()
	file_pkg_core_v1_internals_status_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_pkg_core_v1_internals_service_invocation_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Actor); i {
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
		file_pkg_core_v1_internals_service_invocation_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InternalInvokeRequest); i {
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
		file_pkg_core_v1_internals_service_invocation_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InternalInvokeResponse); i {
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
		file_pkg_core_v1_internals_service_invocation_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListStringValue); i {
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
			RawDescriptor: file_pkg_core_v1_internals_service_invocation_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_core_v1_internals_service_invocation_proto_goTypes,
		DependencyIndexes: file_pkg_core_v1_internals_service_invocation_proto_depIdxs,
		MessageInfos:      file_pkg_core_v1_internals_service_invocation_proto_msgTypes,
	}.Build()
	File_pkg_core_v1_internals_service_invocation_proto = out.File
	file_pkg_core_v1_internals_service_invocation_proto_rawDesc = nil
	file_pkg_core_v1_internals_service_invocation_proto_goTypes = nil
	file_pkg_core_v1_internals_service_invocation_proto_depIdxs = nil
}