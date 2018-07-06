// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc.proto

package grpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type OrderType int32

const (
	OrderType_Midpoint OrderType = 0
	OrderType_Limit    OrderType = 1
)

var OrderType_name = map[int32]string{
	0: "Midpoint",
	1: "Limit",
}
var OrderType_value = map[string]int32{
	"Midpoint": 0,
	"Limit":    1,
}

func (x OrderType) String() string {
	return proto.EnumName(OrderType_name, int32(x))
}
func (OrderType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{0}
}

type OrderParity int32

const (
	OrderParity_Buy  OrderParity = 0
	OrderParity_Sell OrderParity = 1
)

var OrderParity_name = map[int32]string{
	0: "Buy",
	1: "Sell",
}
var OrderParity_value = map[string]int32{
	"Buy":  0,
	"Sell": 1,
}

func (x OrderParity) String() string {
	return proto.EnumName(OrderParity_name, int32(x))
}
func (OrderParity) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{1}
}

type OrderSettlement int32

const (
	OrderSettlement_Nil         OrderSettlement = 0
	OrderSettlement_RenEx       OrderSettlement = 1
	OrderSettlement_RenExAtomic OrderSettlement = 2
)

var OrderSettlement_name = map[int32]string{
	0: "Nil",
	1: "RenEx",
	2: "RenExAtomic",
}
var OrderSettlement_value = map[string]int32{
	"Nil":         0,
	"RenEx":       1,
	"RenExAtomic": 2,
}

func (x OrderSettlement) String() string {
	return proto.EnumName(OrderSettlement_name, int32(x))
}
func (OrderSettlement) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{2}
}

type PingRequest struct {
	Signature            []byte   `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	MultiAddress         string   `protobuf:"bytes,2,opt,name=multiAddress,proto3" json:"multiAddress,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingRequest) Reset()         { *m = PingRequest{} }
func (m *PingRequest) String() string { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()    {}
func (*PingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{0}
}
func (m *PingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingRequest.Unmarshal(m, b)
}
func (m *PingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingRequest.Marshal(b, m, deterministic)
}
func (dst *PingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingRequest.Merge(dst, src)
}
func (m *PingRequest) XXX_Size() int {
	return xxx_messageInfo_PingRequest.Size(m)
}
func (m *PingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PingRequest proto.InternalMessageInfo

func (m *PingRequest) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *PingRequest) GetMultiAddress() string {
	if m != nil {
		return m.MultiAddress
	}
	return ""
}

type PingResponse struct {
	Signature            []byte   `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	MultiAddress         string   `protobuf:"bytes,2,opt,name=multiAddress,proto3" json:"multiAddress,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingResponse) Reset()         { *m = PingResponse{} }
func (m *PingResponse) String() string { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()    {}
func (*PingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{1}
}
func (m *PingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingResponse.Unmarshal(m, b)
}
func (m *PingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingResponse.Marshal(b, m, deterministic)
}
func (dst *PingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingResponse.Merge(dst, src)
}
func (m *PingResponse) XXX_Size() int {
	return xxx_messageInfo_PingResponse.Size(m)
}
func (m *PingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PingResponse proto.InternalMessageInfo

func (m *PingResponse) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *PingResponse) GetMultiAddress() string {
	if m != nil {
		return m.MultiAddress
	}
	return ""
}

type QueryRequest struct {
	Signature            []byte   `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	Address              string   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryRequest) Reset()         { *m = QueryRequest{} }
func (m *QueryRequest) String() string { return proto.CompactTextString(m) }
func (*QueryRequest) ProtoMessage()    {}
func (*QueryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{2}
}
func (m *QueryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryRequest.Unmarshal(m, b)
}
func (m *QueryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryRequest.Marshal(b, m, deterministic)
}
func (dst *QueryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryRequest.Merge(dst, src)
}
func (m *QueryRequest) XXX_Size() int {
	return xxx_messageInfo_QueryRequest.Size(m)
}
func (m *QueryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryRequest proto.InternalMessageInfo

func (m *QueryRequest) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *QueryRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type QueryResponse struct {
	Signature            []byte   `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	MultiAddress         string   `protobuf:"bytes,2,opt,name=multiAddress,proto3" json:"multiAddress,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryResponse) Reset()         { *m = QueryResponse{} }
func (m *QueryResponse) String() string { return proto.CompactTextString(m) }
func (*QueryResponse) ProtoMessage()    {}
func (*QueryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{3}
}
func (m *QueryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryResponse.Unmarshal(m, b)
}
func (m *QueryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryResponse.Marshal(b, m, deterministic)
}
func (dst *QueryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResponse.Merge(dst, src)
}
func (m *QueryResponse) XXX_Size() int {
	return xxx_messageInfo_QueryResponse.Size(m)
}
func (m *QueryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResponse proto.InternalMessageInfo

func (m *QueryResponse) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *QueryResponse) GetMultiAddress() string {
	if m != nil {
		return m.MultiAddress
	}
	return ""
}

type StreamMessage struct {
	Authentication       *StreamAuthentication `protobuf:"bytes,1,opt,name=authentication,proto3" json:"authentication,omitempty"`
	Data                 []byte                `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *StreamMessage) Reset()         { *m = StreamMessage{} }
func (m *StreamMessage) String() string { return proto.CompactTextString(m) }
func (*StreamMessage) ProtoMessage()    {}
func (*StreamMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{4}
}
func (m *StreamMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamMessage.Unmarshal(m, b)
}
func (m *StreamMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamMessage.Marshal(b, m, deterministic)
}
func (dst *StreamMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamMessage.Merge(dst, src)
}
func (m *StreamMessage) XXX_Size() int {
	return xxx_messageInfo_StreamMessage.Size(m)
}
func (m *StreamMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamMessage.DiscardUnknown(m)
}

var xxx_messageInfo_StreamMessage proto.InternalMessageInfo

func (m *StreamMessage) GetAuthentication() *StreamAuthentication {
	if m != nil {
		return m.Authentication
	}
	return nil
}

func (m *StreamMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type StreamAuthentication struct {
	Signature            []byte   `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	Address              string   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamAuthentication) Reset()         { *m = StreamAuthentication{} }
func (m *StreamAuthentication) String() string { return proto.CompactTextString(m) }
func (*StreamAuthentication) ProtoMessage()    {}
func (*StreamAuthentication) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{5}
}
func (m *StreamAuthentication) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamAuthentication.Unmarshal(m, b)
}
func (m *StreamAuthentication) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamAuthentication.Marshal(b, m, deterministic)
}
func (dst *StreamAuthentication) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamAuthentication.Merge(dst, src)
}
func (m *StreamAuthentication) XXX_Size() int {
	return xxx_messageInfo_StreamAuthentication.Size(m)
}
func (m *StreamAuthentication) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamAuthentication.DiscardUnknown(m)
}

var xxx_messageInfo_StreamAuthentication proto.InternalMessageInfo

func (m *StreamAuthentication) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *StreamAuthentication) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type OpenOrderRequest struct {
	OrderFragment        *EncryptedOrderFragment `protobuf:"bytes,1,opt,name=orderFragment,proto3" json:"orderFragment,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *OpenOrderRequest) Reset()         { *m = OpenOrderRequest{} }
func (m *OpenOrderRequest) String() string { return proto.CompactTextString(m) }
func (*OpenOrderRequest) ProtoMessage()    {}
func (*OpenOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{6}
}
func (m *OpenOrderRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OpenOrderRequest.Unmarshal(m, b)
}
func (m *OpenOrderRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OpenOrderRequest.Marshal(b, m, deterministic)
}
func (dst *OpenOrderRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OpenOrderRequest.Merge(dst, src)
}
func (m *OpenOrderRequest) XXX_Size() int {
	return xxx_messageInfo_OpenOrderRequest.Size(m)
}
func (m *OpenOrderRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OpenOrderRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OpenOrderRequest proto.InternalMessageInfo

func (m *OpenOrderRequest) GetOrderFragment() *EncryptedOrderFragment {
	if m != nil {
		return m.OrderFragment
	}
	return nil
}

type OpenOrderResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OpenOrderResponse) Reset()         { *m = OpenOrderResponse{} }
func (m *OpenOrderResponse) String() string { return proto.CompactTextString(m) }
func (*OpenOrderResponse) ProtoMessage()    {}
func (*OpenOrderResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{7}
}
func (m *OpenOrderResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OpenOrderResponse.Unmarshal(m, b)
}
func (m *OpenOrderResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OpenOrderResponse.Marshal(b, m, deterministic)
}
func (dst *OpenOrderResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OpenOrderResponse.Merge(dst, src)
}
func (m *OpenOrderResponse) XXX_Size() int {
	return xxx_messageInfo_OpenOrderResponse.Size(m)
}
func (m *OpenOrderResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OpenOrderResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OpenOrderResponse proto.InternalMessageInfo

type EncryptedOrderFragment struct {
	OrderId              []byte               `protobuf:"bytes,1,opt,name=orderId,proto3" json:"orderId,omitempty"`
	OrderType            OrderType            `protobuf:"varint,2,opt,name=orderType,proto3,enum=grpc.OrderType" json:"orderType,omitempty"`
	OrderParity          OrderParity          `protobuf:"varint,3,opt,name=orderParity,proto3,enum=grpc.OrderParity" json:"orderParity,omitempty"`
	OrderSettlement      OrderSettlement      `protobuf:"varint,4,opt,name=orderSettlement,proto3,enum=grpc.OrderSettlement" json:"orderSettlement,omitempty"`
	OrderExpiry          int64                `protobuf:"varint,5,opt,name=orderExpiry,proto3" json:"orderExpiry,omitempty"`
	Id                   []byte               `protobuf:"bytes,6,opt,name=id,proto3" json:"id,omitempty"`
	EpochDepth           int32                `protobuf:"varint,7,opt,name=epochDepth,proto3" json:"epochDepth,omitempty"`
	Tokens               []byte               `protobuf:"bytes,8,opt,name=tokens,proto3" json:"tokens,omitempty"`
	Price                *EncryptedCoExpShare `protobuf:"bytes,9,opt,name=price,proto3" json:"price,omitempty"`
	Volume               *EncryptedCoExpShare `protobuf:"bytes,10,opt,name=volume,proto3" json:"volume,omitempty"`
	MinimumVolume        *EncryptedCoExpShare `protobuf:"bytes,11,opt,name=minimumVolume,proto3" json:"minimumVolume,omitempty"`
	Nonce                []byte               `protobuf:"bytes,12,opt,name=nonce,proto3" json:"nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *EncryptedOrderFragment) Reset()         { *m = EncryptedOrderFragment{} }
func (m *EncryptedOrderFragment) String() string { return proto.CompactTextString(m) }
func (*EncryptedOrderFragment) ProtoMessage()    {}
func (*EncryptedOrderFragment) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{8}
}
func (m *EncryptedOrderFragment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EncryptedOrderFragment.Unmarshal(m, b)
}
func (m *EncryptedOrderFragment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EncryptedOrderFragment.Marshal(b, m, deterministic)
}
func (dst *EncryptedOrderFragment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncryptedOrderFragment.Merge(dst, src)
}
func (m *EncryptedOrderFragment) XXX_Size() int {
	return xxx_messageInfo_EncryptedOrderFragment.Size(m)
}
func (m *EncryptedOrderFragment) XXX_DiscardUnknown() {
	xxx_messageInfo_EncryptedOrderFragment.DiscardUnknown(m)
}

var xxx_messageInfo_EncryptedOrderFragment proto.InternalMessageInfo

func (m *EncryptedOrderFragment) GetOrderId() []byte {
	if m != nil {
		return m.OrderId
	}
	return nil
}

func (m *EncryptedOrderFragment) GetOrderType() OrderType {
	if m != nil {
		return m.OrderType
	}
	return OrderType_Midpoint
}

func (m *EncryptedOrderFragment) GetOrderParity() OrderParity {
	if m != nil {
		return m.OrderParity
	}
	return OrderParity_Buy
}

func (m *EncryptedOrderFragment) GetOrderSettlement() OrderSettlement {
	if m != nil {
		return m.OrderSettlement
	}
	return OrderSettlement_Nil
}

func (m *EncryptedOrderFragment) GetOrderExpiry() int64 {
	if m != nil {
		return m.OrderExpiry
	}
	return 0
}

func (m *EncryptedOrderFragment) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *EncryptedOrderFragment) GetEpochDepth() int32 {
	if m != nil {
		return m.EpochDepth
	}
	return 0
}

func (m *EncryptedOrderFragment) GetTokens() []byte {
	if m != nil {
		return m.Tokens
	}
	return nil
}

func (m *EncryptedOrderFragment) GetPrice() *EncryptedCoExpShare {
	if m != nil {
		return m.Price
	}
	return nil
}

func (m *EncryptedOrderFragment) GetVolume() *EncryptedCoExpShare {
	if m != nil {
		return m.Volume
	}
	return nil
}

func (m *EncryptedOrderFragment) GetMinimumVolume() *EncryptedCoExpShare {
	if m != nil {
		return m.MinimumVolume
	}
	return nil
}

func (m *EncryptedOrderFragment) GetNonce() []byte {
	if m != nil {
		return m.Nonce
	}
	return nil
}

type EncryptedCoExpShare struct {
	Co                   []byte   `protobuf:"bytes,1,opt,name=co,proto3" json:"co,omitempty"`
	Exp                  []byte   `protobuf:"bytes,2,opt,name=exp,proto3" json:"exp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EncryptedCoExpShare) Reset()         { *m = EncryptedCoExpShare{} }
func (m *EncryptedCoExpShare) String() string { return proto.CompactTextString(m) }
func (*EncryptedCoExpShare) ProtoMessage()    {}
func (*EncryptedCoExpShare) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{9}
}
func (m *EncryptedCoExpShare) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EncryptedCoExpShare.Unmarshal(m, b)
}
func (m *EncryptedCoExpShare) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EncryptedCoExpShare.Marshal(b, m, deterministic)
}
func (dst *EncryptedCoExpShare) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncryptedCoExpShare.Merge(dst, src)
}
func (m *EncryptedCoExpShare) XXX_Size() int {
	return xxx_messageInfo_EncryptedCoExpShare.Size(m)
}
func (m *EncryptedCoExpShare) XXX_DiscardUnknown() {
	xxx_messageInfo_EncryptedCoExpShare.DiscardUnknown(m)
}

var xxx_messageInfo_EncryptedCoExpShare proto.InternalMessageInfo

func (m *EncryptedCoExpShare) GetCo() []byte {
	if m != nil {
		return m.Co
	}
	return nil
}

func (m *EncryptedCoExpShare) GetExp() []byte {
	if m != nil {
		return m.Exp
	}
	return nil
}

type StatusRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusRequest) Reset()         { *m = StatusRequest{} }
func (m *StatusRequest) String() string { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()    {}
func (*StatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{10}
}
func (m *StatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusRequest.Unmarshal(m, b)
}
func (m *StatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusRequest.Marshal(b, m, deterministic)
}
func (dst *StatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusRequest.Merge(dst, src)
}
func (m *StatusRequest) XXX_Size() int {
	return xxx_messageInfo_StatusRequest.Size(m)
}
func (m *StatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StatusRequest proto.InternalMessageInfo

type StatusResponse struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Bootstrapped         bool     `protobuf:"varint,2,opt,name=bootstrapped,proto3" json:"bootstrapped,omitempty"`
	Peers                int64    `protobuf:"varint,3,opt,name=peers,proto3" json:"peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusResponse) Reset()         { *m = StatusResponse{} }
func (m *StatusResponse) String() string { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()    {}
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_619b6be9bea7a31c, []int{11}
}
func (m *StatusResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusResponse.Unmarshal(m, b)
}
func (m *StatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusResponse.Marshal(b, m, deterministic)
}
func (dst *StatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusResponse.Merge(dst, src)
}
func (m *StatusResponse) XXX_Size() int {
	return xxx_messageInfo_StatusResponse.Size(m)
}
func (m *StatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StatusResponse proto.InternalMessageInfo

func (m *StatusResponse) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *StatusResponse) GetBootstrapped() bool {
	if m != nil {
		return m.Bootstrapped
	}
	return false
}

func (m *StatusResponse) GetPeers() int64 {
	if m != nil {
		return m.Peers
	}
	return 0
}

func init() {
	proto.RegisterType((*PingRequest)(nil), "grpc.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "grpc.PingResponse")
	proto.RegisterType((*QueryRequest)(nil), "grpc.QueryRequest")
	proto.RegisterType((*QueryResponse)(nil), "grpc.QueryResponse")
	proto.RegisterType((*StreamMessage)(nil), "grpc.StreamMessage")
	proto.RegisterType((*StreamAuthentication)(nil), "grpc.StreamAuthentication")
	proto.RegisterType((*OpenOrderRequest)(nil), "grpc.OpenOrderRequest")
	proto.RegisterType((*OpenOrderResponse)(nil), "grpc.OpenOrderResponse")
	proto.RegisterType((*EncryptedOrderFragment)(nil), "grpc.EncryptedOrderFragment")
	proto.RegisterType((*EncryptedCoExpShare)(nil), "grpc.EncryptedCoExpShare")
	proto.RegisterType((*StatusRequest)(nil), "grpc.StatusRequest")
	proto.RegisterType((*StatusResponse)(nil), "grpc.StatusResponse")
	proto.RegisterEnum("grpc.OrderType", OrderType_name, OrderType_value)
	proto.RegisterEnum("grpc.OrderParity", OrderParity_name, OrderParity_value)
	proto.RegisterEnum("grpc.OrderSettlement", OrderSettlement_name, OrderSettlement_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SwarmServiceClient is the client API for SwarmService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SwarmServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (SwarmService_QueryClient, error)
}

type swarmServiceClient struct {
	cc *grpc.ClientConn
}

func NewSwarmServiceClient(cc *grpc.ClientConn) SwarmServiceClient {
	return &swarmServiceClient{cc}
}

func (c *swarmServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/grpc.SwarmService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swarmServiceClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (SwarmService_QueryClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SwarmService_serviceDesc.Streams[0], "/grpc.SwarmService/Query", opts...)
	if err != nil {
		return nil, err
	}
	x := &swarmServiceQueryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SwarmService_QueryClient interface {
	Recv() (*QueryResponse, error)
	grpc.ClientStream
}

type swarmServiceQueryClient struct {
	grpc.ClientStream
}

func (x *swarmServiceQueryClient) Recv() (*QueryResponse, error) {
	m := new(QueryResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SwarmServiceServer is the server API for SwarmService service.
type SwarmServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Query(*QueryRequest, SwarmService_QueryServer) error
}

func RegisterSwarmServiceServer(s *grpc.Server, srv SwarmServiceServer) {
	s.RegisterService(&_SwarmService_serviceDesc, srv)
}

func _SwarmService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwarmServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.SwarmService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwarmServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwarmService_Query_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SwarmServiceServer).Query(m, &swarmServiceQueryServer{stream})
}

type SwarmService_QueryServer interface {
	Send(*QueryResponse) error
	grpc.ServerStream
}

type swarmServiceQueryServer struct {
	grpc.ServerStream
}

func (x *swarmServiceQueryServer) Send(m *QueryResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _SwarmService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.SwarmService",
	HandlerType: (*SwarmServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _SwarmService_Ping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Query",
			Handler:       _SwarmService_Query_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "grpc.proto",
}

// StreamServiceClient is the client API for StreamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StreamServiceClient interface {
	Connect(ctx context.Context, opts ...grpc.CallOption) (StreamService_ConnectClient, error)
}

type streamServiceClient struct {
	cc *grpc.ClientConn
}

func NewStreamServiceClient(cc *grpc.ClientConn) StreamServiceClient {
	return &streamServiceClient{cc}
}

func (c *streamServiceClient) Connect(ctx context.Context, opts ...grpc.CallOption) (StreamService_ConnectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_StreamService_serviceDesc.Streams[0], "/grpc.StreamService/Connect", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamServiceConnectClient{stream}
	return x, nil
}

type StreamService_ConnectClient interface {
	Send(*StreamMessage) error
	Recv() (*StreamMessage, error)
	grpc.ClientStream
}

type streamServiceConnectClient struct {
	grpc.ClientStream
}

func (x *streamServiceConnectClient) Send(m *StreamMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *streamServiceConnectClient) Recv() (*StreamMessage, error) {
	m := new(StreamMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamServiceServer is the server API for StreamService service.
type StreamServiceServer interface {
	Connect(StreamService_ConnectServer) error
}

func RegisterStreamServiceServer(s *grpc.Server, srv StreamServiceServer) {
	s.RegisterService(&_StreamService_serviceDesc, srv)
}

func _StreamService_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamServiceServer).Connect(&streamServiceConnectServer{stream})
}

type StreamService_ConnectServer interface {
	Send(*StreamMessage) error
	Recv() (*StreamMessage, error)
	grpc.ServerStream
}

type streamServiceConnectServer struct {
	grpc.ServerStream
}

func (x *streamServiceConnectServer) Send(m *StreamMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *streamServiceConnectServer) Recv() (*StreamMessage, error) {
	m := new(StreamMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _StreamService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.StreamService",
	HandlerType: (*StreamServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _StreamService_Connect_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "grpc.proto",
}

// OrderbookServiceClient is the client API for OrderbookService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrderbookServiceClient interface {
	OpenOrder(ctx context.Context, in *OpenOrderRequest, opts ...grpc.CallOption) (*OpenOrderResponse, error)
}

type orderbookServiceClient struct {
	cc *grpc.ClientConn
}

func NewOrderbookServiceClient(cc *grpc.ClientConn) OrderbookServiceClient {
	return &orderbookServiceClient{cc}
}

func (c *orderbookServiceClient) OpenOrder(ctx context.Context, in *OpenOrderRequest, opts ...grpc.CallOption) (*OpenOrderResponse, error) {
	out := new(OpenOrderResponse)
	err := c.cc.Invoke(ctx, "/grpc.OrderbookService/OpenOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderbookServiceServer is the server API for OrderbookService service.
type OrderbookServiceServer interface {
	OpenOrder(context.Context, *OpenOrderRequest) (*OpenOrderResponse, error)
}

func RegisterOrderbookServiceServer(s *grpc.Server, srv OrderbookServiceServer) {
	s.RegisterService(&_OrderbookService_serviceDesc, srv)
}

func _OrderbookService_OpenOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderbookServiceServer).OpenOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.OrderbookService/OpenOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderbookServiceServer).OpenOrder(ctx, req.(*OpenOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrderbookService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.OrderbookService",
	HandlerType: (*OrderbookServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OpenOrder",
			Handler:    _OrderbookService_OpenOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc.proto",
}

// StatusServiceClient is the client API for StatusService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StatusServiceClient interface {
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
}

type statusServiceClient struct {
	cc *grpc.ClientConn
}

func NewStatusServiceClient(cc *grpc.ClientConn) StatusServiceClient {
	return &statusServiceClient{cc}
}

func (c *statusServiceClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/grpc.StatusService/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatusServiceServer is the server API for StatusService service.
type StatusServiceServer interface {
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
}

func RegisterStatusServiceServer(s *grpc.Server, srv StatusServiceServer) {
	s.RegisterService(&_StatusService_serviceDesc, srv)
}

func _StatusService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.StatusService/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusServiceServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _StatusService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.StatusService",
	HandlerType: (*StatusServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _StatusService_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc.proto",
}

func init() { proto.RegisterFile("grpc.proto", fileDescriptor_grpc_619b6be9bea7a31c) }

var fileDescriptor_grpc_619b6be9bea7a31c = []byte{
	// 741 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0x4d, 0x6f, 0xe3, 0x36,
	0x10, 0x5d, 0xf9, 0xdb, 0x63, 0xd9, 0xd6, 0xd2, 0x69, 0xaa, 0x1a, 0x8b, 0xc2, 0x10, 0x7a, 0x30,
	0x02, 0xec, 0x76, 0xeb, 0x00, 0xcd, 0xa5, 0x40, 0x90, 0x0f, 0x07, 0x2d, 0xd0, 0xc4, 0x8e, 0x5c,
	0xe4, 0xae, 0x48, 0x03, 0x9b, 0x88, 0x45, 0x32, 0x14, 0x95, 0xda, 0xff, 0xbd, 0x87, 0x42, 0x94,
	0x14, 0x4b, 0xae, 0x81, 0x14, 0x8b, 0xdc, 0x38, 0x4f, 0xf3, 0x1e, 0x1f, 0x87, 0xc3, 0x11, 0xc0,
	0x52, 0x0a, 0xff, 0x8b, 0x90, 0x5c, 0x71, 0x52, 0x4b, 0xd6, 0xce, 0x0c, 0x3a, 0x73, 0xca, 0x96,
	0x2e, 0x3e, 0xc7, 0x18, 0x29, 0xf2, 0x09, 0xda, 0x11, 0x5d, 0x32, 0x4f, 0xc5, 0x12, 0x6d, 0x63,
	0x64, 0x8c, 0x4d, 0x77, 0x07, 0x10, 0x07, 0xcc, 0x30, 0x5e, 0x2b, 0x7a, 0x11, 0x04, 0x12, 0xa3,
	0xc8, 0xae, 0x8c, 0x8c, 0x71, 0xdb, 0x2d, 0x61, 0xce, 0x1c, 0xcc, 0x54, 0x30, 0x12, 0x9c, 0x45,
	0xf8, 0x0e, 0x8a, 0x37, 0x60, 0xde, 0xc7, 0x28, 0xb7, 0xff, 0xcf, 0xa3, 0x0d, 0x4d, 0xaf, 0x24,
	0x96, 0x87, 0xce, 0x3d, 0x74, 0x33, 0x9d, 0x77, 0xb3, 0xb6, 0x84, 0xee, 0x42, 0x49, 0xf4, 0xc2,
	0x5b, 0x8c, 0x22, 0x6f, 0x89, 0xe4, 0x12, 0x7a, 0x5e, 0xac, 0x56, 0xc8, 0x14, 0xf5, 0x3d, 0x45,
	0x39, 0xd3, 0xba, 0x9d, 0xc9, 0xf0, 0x8b, 0xae, 0x7c, 0x9a, 0x7c, 0x51, 0xca, 0x70, 0xf7, 0x18,
	0x84, 0x40, 0x2d, 0xf0, 0x94, 0xa7, 0x37, 0x34, 0x5d, 0xbd, 0x76, 0xee, 0xe0, 0xe8, 0x10, 0xf7,
	0x9b, 0x6b, 0xf1, 0x00, 0xd6, 0x4c, 0x20, 0x9b, 0xc9, 0x00, 0x65, 0x5e, 0xd7, 0x4b, 0xe8, 0xf2,
	0x24, 0xbe, 0x91, 0xde, 0x32, 0x44, 0xa6, 0x32, 0xeb, 0x9f, 0x52, 0xeb, 0x53, 0xe6, 0xcb, 0xad,
	0x50, 0x18, 0xcc, 0x8a, 0x39, 0x6e, 0x99, 0xe2, 0x0c, 0xe0, 0x63, 0x41, 0x37, 0xad, 0xb3, 0xf3,
	0x4f, 0x15, 0x8e, 0x0f, 0xd3, 0x13, 0x87, 0x5a, 0xe0, 0x8f, 0x20, 0x73, 0x9f, 0x87, 0xe4, 0x33,
	0xb4, 0xf5, 0xf2, 0xaf, 0xad, 0x40, 0xed, 0xbe, 0x37, 0xe9, 0xa7, 0x4e, 0x66, 0x39, 0xec, 0xee,
	0x32, 0xc8, 0x29, 0x74, 0x74, 0x30, 0xf7, 0x24, 0x55, 0x5b, 0xbb, 0xaa, 0x09, 0x1f, 0x0b, 0x84,
	0xf4, 0x83, 0x5b, 0xcc, 0x22, 0xe7, 0xd0, 0xd7, 0xe1, 0x02, 0x95, 0x5a, 0xa3, 0x3e, 0x73, 0x4d,
	0x13, 0xbf, 0x2b, 0x10, 0x77, 0x1f, 0xdd, 0xfd, 0x6c, 0x32, 0xca, 0x76, 0x9d, 0x6e, 0x04, 0x95,
	0x5b, 0xbb, 0x3e, 0x32, 0xc6, 0x55, 0xb7, 0x08, 0x91, 0x1e, 0x54, 0x68, 0x60, 0x37, 0xf4, 0xd9,
	0x2a, 0x34, 0x20, 0x3f, 0x02, 0xa0, 0xe0, 0xfe, 0xea, 0x1a, 0x85, 0x5a, 0xd9, 0xcd, 0x91, 0x31,
	0xae, 0xbb, 0x05, 0x84, 0x1c, 0x43, 0x43, 0xf1, 0x27, 0x64, 0x91, 0xdd, 0xd2, 0x9c, 0x2c, 0x22,
	0x3f, 0x43, 0x5d, 0x48, 0xea, 0xa3, 0xdd, 0xd6, 0x97, 0xf2, 0xc3, 0xde, 0xa5, 0x5c, 0xf1, 0xe9,
	0x46, 0x2c, 0x56, 0x9e, 0x44, 0x37, 0xcd, 0x23, 0xbf, 0x40, 0xe3, 0x85, 0xaf, 0xe3, 0x10, 0x6d,
	0x78, 0x8b, 0x91, 0x25, 0x92, 0x73, 0xe8, 0x86, 0x94, 0xd1, 0x30, 0x0e, 0x1f, 0x52, 0x66, 0xe7,
	0x2d, 0x66, 0x39, 0x9f, 0x1c, 0x41, 0x9d, 0x71, 0xe6, 0xa3, 0x6d, 0x6a, 0xef, 0x69, 0xe0, 0x9c,
	0xc1, 0xe0, 0x00, 0x37, 0xa9, 0x8c, 0xcf, 0xb3, 0x5b, 0xaf, 0xf8, 0x9c, 0x58, 0x50, 0xc5, 0x8d,
	0xc8, 0xba, 0x3e, 0x59, 0x3a, 0xfd, 0xe4, 0x75, 0x79, 0x2a, 0x8e, 0xb2, 0x0e, 0x75, 0x02, 0xe8,
	0xe5, 0x40, 0xf6, 0x84, 0x0b, 0x1d, 0x6e, 0x94, 0x3a, 0x3c, 0x79, 0xbe, 0x8f, 0x9c, 0xab, 0x48,
	0x49, 0x4f, 0x08, 0x0c, 0xb4, 0x6e, 0xcb, 0x2d, 0x61, 0x89, 0x5f, 0x81, 0x28, 0x23, 0xdd, 0x2e,
	0x55, 0x37, 0x0d, 0x4e, 0x7e, 0x82, 0xf6, 0x6b, 0x8b, 0x11, 0x13, 0x5a, 0xb7, 0x34, 0x10, 0x9c,
	0x32, 0x65, 0x7d, 0x20, 0x6d, 0xa8, 0xff, 0x49, 0x43, 0xaa, 0x2c, 0xe3, 0x64, 0x04, 0x9d, 0x42,
	0x5f, 0x91, 0x26, 0x54, 0x2f, 0xe3, 0xad, 0xf5, 0x81, 0xb4, 0xa0, 0xb6, 0xc0, 0xf5, 0xda, 0x32,
	0x4e, 0x7e, 0x85, 0xfe, 0x5e, 0x03, 0x25, 0x59, 0x77, 0x74, 0x9d, 0x0a, 0xb9, 0xc8, 0xa6, 0x1b,
	0xcb, 0x20, 0x7d, 0xe8, 0xe8, 0xe5, 0x85, 0xe2, 0x21, 0xf5, 0xad, 0xca, 0xe4, 0x19, 0xcc, 0xc5,
	0xdf, 0x9e, 0x0c, 0x17, 0x28, 0x5f, 0x92, 0x9b, 0xfc, 0x0c, 0xb5, 0x64, 0xa2, 0x92, 0xac, 0x9b,
	0x0b, 0xe3, 0x7a, 0x48, 0x8a, 0x50, 0x56, 0x92, 0x09, 0xd4, 0xf5, 0x98, 0x23, 0xd9, 0xc7, 0xe2,
	0xec, 0x1c, 0x0e, 0x4a, 0x58, 0xca, 0xf8, 0x6a, 0x4c, 0x7e, 0xcf, 0xe7, 0x58, 0xbe, 0xe7, 0x19,
	0x34, 0xaf, 0x38, 0x63, 0xe8, 0x2b, 0x32, 0x28, 0x8e, 0xae, 0x6c, 0xce, 0x0d, 0x0f, 0x81, 0x63,
	0xe3, 0xab, 0x31, 0x99, 0x83, 0xa5, 0x0f, 0xfd, 0xc8, 0xf9, 0x53, 0x2e, 0xf6, 0x1b, 0xb4, 0x5f,
	0x87, 0x02, 0x39, 0xce, 0x9e, 0xd6, 0xde, 0xf4, 0x19, 0x7e, 0xff, 0x1f, 0x3c, 0x75, 0x37, 0xb9,
	0xce, 0xbb, 0x20, 0x97, 0x3b, 0x85, 0x46, 0x0a, 0xec, 0xac, 0x15, 0x9a, 0x64, 0x78, 0x54, 0x06,
	0x53, 0x95, 0xc7, 0x86, 0xfe, 0xe9, 0x9d, 0xfe, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x9b, 0x81, 0xa2,
	0xb2, 0x02, 0x07, 0x00, 0x00,
}
