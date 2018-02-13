// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc.proto

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	rpc.proto

It has these top-level messages:
	Address
	MultiAddress
	MultiAddresses
	Nothing
	Query
	OrderFragment
	ResultFragment
	Result
	Results
	Atom
*/
package rpc

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

// An Address message is the network representation of an identity.Address.
type Address struct {
	Address string `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
}

func (m *Address) Reset()                    { *m = Address{} }
func (m *Address) String() string            { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()               {}
func (*Address) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Address) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

// A MultiAddress is the network representation of an identity.MultiAddress.
type MultiAddress struct {
	Multi string `protobuf:"bytes,1,opt,name=multi" json:"multi,omitempty"`
}

func (m *MultiAddress) Reset()                    { *m = MultiAddress{} }
func (m *MultiAddress) String() string            { return proto.CompactTextString(m) }
func (*MultiAddress) ProtoMessage()               {}
func (*MultiAddress) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MultiAddress) GetMulti() string {
	if m != nil {
		return m.Multi
	}
	return ""
}

// MultiAddresses are the network representation of identity.MultiAddresses.
type MultiAddresses struct {
	Multis []*MultiAddress `protobuf:"bytes,1,rep,name=multis" json:"multis,omitempty"`
}

func (m *MultiAddresses) Reset()                    { *m = MultiAddresses{} }
func (m *MultiAddresses) String() string            { return proto.CompactTextString(m) }
func (*MultiAddresses) ProtoMessage()               {}
func (*MultiAddresses) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MultiAddresses) GetMultis() []*MultiAddress {
	if m != nil {
		return m.Multis
	}
	return nil
}

// Nothing is in this message. It is used to send nothing, or signal a
// successful response.
type Nothing struct {
}

func (m *Nothing) Reset()                    { *m = Nothing{} }
func (m *Nothing) String() string            { return proto.CompactTextString(m) }
func (*Nothing) ProtoMessage()               {}
func (*Nothing) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// A Query message contains the Address of a Node that needs to be found and
// the MultiAddress of the Node from which the Query originated.
type Query struct {
	// Network data.
	From *MultiAddress `protobuf:"bytes,1,opt,name=from" json:"from,omitempty"`
	// Public data.
	Query *Address `protobuf:"bytes,2,opt,name=query" json:"query,omitempty"`
}

func (m *Query) Reset()                    { *m = Query{} }
func (m *Query) String() string            { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()               {}
func (*Query) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Query) GetFrom() *MultiAddress {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Query) GetQuery() *Address {
	if m != nil {
		return m.Query
	}
	return nil
}

// An OrderFragment is a message contains the details of an order fragment.
type OrderFragment struct {
	// Network data.
	To   *Address      `protobuf:"bytes,1,opt,name=to" json:"to,omitempty"`
	From *MultiAddress `protobuf:"bytes,2,opt,name=from" json:"from,omitempty"`
	// Public data.
	Id          []byte `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	OrderId     []byte `protobuf:"bytes,4,opt,name=orderId,proto3" json:"orderId,omitempty"`
	OrderType   int64  `protobuf:"varint,5,opt,name=orderType" json:"orderType,omitempty"`
	OrderParity int64  `protobuf:"varint,6,opt,name=orderParity" json:"orderParity,omitempty"`
	// Secure data.
	FstCodeShare   []byte `protobuf:"bytes,7,opt,name=fstCodeShare,proto3" json:"fstCodeShare,omitempty"`
	SndCodeShare   []byte `protobuf:"bytes,8,opt,name=sndCodeShare,proto3" json:"sndCodeShare,omitempty"`
	PriceShare     []byte `protobuf:"bytes,9,opt,name=priceShare,proto3" json:"priceShare,omitempty"`
	MaxVolumeShare []byte `protobuf:"bytes,10,opt,name=maxVolumeShare,proto3" json:"maxVolumeShare,omitempty"`
	MinVolumeShare []byte `protobuf:"bytes,11,opt,name=minVolumeShare,proto3" json:"minVolumeShare,omitempty"`
}

func (m *OrderFragment) Reset()                    { *m = OrderFragment{} }
func (m *OrderFragment) String() string            { return proto.CompactTextString(m) }
func (*OrderFragment) ProtoMessage()               {}
func (*OrderFragment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *OrderFragment) GetTo() *Address {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *OrderFragment) GetFrom() *MultiAddress {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *OrderFragment) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *OrderFragment) GetOrderId() []byte {
	if m != nil {
		return m.OrderId
	}
	return nil
}

func (m *OrderFragment) GetOrderType() int64 {
	if m != nil {
		return m.OrderType
	}
	return 0
}

func (m *OrderFragment) GetOrderParity() int64 {
	if m != nil {
		return m.OrderParity
	}
	return 0
}

func (m *OrderFragment) GetFstCodeShare() []byte {
	if m != nil {
		return m.FstCodeShare
	}
	return nil
}

func (m *OrderFragment) GetSndCodeShare() []byte {
	if m != nil {
		return m.SndCodeShare
	}
	return nil
}

func (m *OrderFragment) GetPriceShare() []byte {
	if m != nil {
		return m.PriceShare
	}
	return nil
}

func (m *OrderFragment) GetMaxVolumeShare() []byte {
	if m != nil {
		return m.MaxVolumeShare
	}
	return nil
}

func (m *OrderFragment) GetMinVolumeShare() []byte {
	if m != nil {
		return m.MinVolumeShare
	}
	return nil
}

// A ResultFragment message is the network representation of a
// compute.ResultFragment and the metadata needed to distribute it through the
// network.
type ResultFragment struct {
	// Network data.
	To   *Address      `protobuf:"bytes,1,opt,name=to" json:"to,omitempty"`
	From *MultiAddress `protobuf:"bytes,2,opt,name=from" json:"from,omitempty"`
	// Public data.
	Id                  []byte `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	BuyOrderId          []byte `protobuf:"bytes,4,opt,name=buyOrderId,proto3" json:"buyOrderId,omitempty"`
	SellOrderId         []byte `protobuf:"bytes,5,opt,name=sellOrderId,proto3" json:"sellOrderId,omitempty"`
	BuyOrderFragmentId  []byte `protobuf:"bytes,6,opt,name=buyOrderFragmentId,proto3" json:"buyOrderFragmentId,omitempty"`
	SellOrderFragmentId []byte `protobuf:"bytes,7,opt,name=sellOrderFragmentId,proto3" json:"sellOrderFragmentId,omitempty"`
	// Secure data.
	FstCodeShare   []byte `protobuf:"bytes,8,opt,name=fstCodeShare,proto3" json:"fstCodeShare,omitempty"`
	SndCodeShare   []byte `protobuf:"bytes,9,opt,name=sndCodeShare,proto3" json:"sndCodeShare,omitempty"`
	PriceShare     []byte `protobuf:"bytes,10,opt,name=priceShare,proto3" json:"priceShare,omitempty"`
	MaxVolumeShare []byte `protobuf:"bytes,11,opt,name=maxVolumeShare,proto3" json:"maxVolumeShare,omitempty"`
	MinVolumeShare []byte `protobuf:"bytes,12,opt,name=minVolumeShare,proto3" json:"minVolumeShare,omitempty"`
}

func (m *ResultFragment) Reset()                    { *m = ResultFragment{} }
func (m *ResultFragment) String() string            { return proto.CompactTextString(m) }
func (*ResultFragment) ProtoMessage()               {}
func (*ResultFragment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ResultFragment) GetTo() *Address {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *ResultFragment) GetFrom() *MultiAddress {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *ResultFragment) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *ResultFragment) GetBuyOrderId() []byte {
	if m != nil {
		return m.BuyOrderId
	}
	return nil
}

func (m *ResultFragment) GetSellOrderId() []byte {
	if m != nil {
		return m.SellOrderId
	}
	return nil
}

func (m *ResultFragment) GetBuyOrderFragmentId() []byte {
	if m != nil {
		return m.BuyOrderFragmentId
	}
	return nil
}

func (m *ResultFragment) GetSellOrderFragmentId() []byte {
	if m != nil {
		return m.SellOrderFragmentId
	}
	return nil
}

func (m *ResultFragment) GetFstCodeShare() []byte {
	if m != nil {
		return m.FstCodeShare
	}
	return nil
}

func (m *ResultFragment) GetSndCodeShare() []byte {
	if m != nil {
		return m.SndCodeShare
	}
	return nil
}

func (m *ResultFragment) GetPriceShare() []byte {
	if m != nil {
		return m.PriceShare
	}
	return nil
}

func (m *ResultFragment) GetMaxVolumeShare() []byte {
	if m != nil {
		return m.MaxVolumeShare
	}
	return nil
}

func (m *ResultFragment) GetMinVolumeShare() []byte {
	if m != nil {
		return m.MinVolumeShare
	}
	return nil
}

// Result messages are sent to signal that a successful order computation has
// happened.
type Result struct {
	// Public data.
	Id          []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	BuyOrderId  []byte `protobuf:"bytes,2,opt,name=buyOrderId,proto3" json:"buyOrderId,omitempty"`
	SellOrderId []byte `protobuf:"bytes,3,opt,name=sellOrderId,proto3" json:"sellOrderId,omitempty"`
	// Secure data.
	FstCode   []byte `protobuf:"bytes,4,opt,name=fstCode,proto3" json:"fstCode,omitempty"`
	SndCode   []byte `protobuf:"bytes,5,opt,name=sndCode,proto3" json:"sndCode,omitempty"`
	Price     []byte `protobuf:"bytes,6,opt,name=price,proto3" json:"price,omitempty"`
	MaxVolume []byte `protobuf:"bytes,7,opt,name=maxVolume,proto3" json:"maxVolume,omitempty"`
	MinVolume []byte `protobuf:"bytes,8,opt,name=minVolume,proto3" json:"minVolume,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Result) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Result) GetBuyOrderId() []byte {
	if m != nil {
		return m.BuyOrderId
	}
	return nil
}

func (m *Result) GetSellOrderId() []byte {
	if m != nil {
		return m.SellOrderId
	}
	return nil
}

func (m *Result) GetFstCode() []byte {
	if m != nil {
		return m.FstCode
	}
	return nil
}

func (m *Result) GetSndCode() []byte {
	if m != nil {
		return m.SndCode
	}
	return nil
}

func (m *Result) GetPrice() []byte {
	if m != nil {
		return m.Price
	}
	return nil
}

func (m *Result) GetMaxVolume() []byte {
	if m != nil {
		return m.MaxVolume
	}
	return nil
}

func (m *Result) GetMinVolume() []byte {
	if m != nil {
		return m.MinVolume
	}
	return nil
}

// Results message is a list of Result message
type Results struct {
	Result []*Result `protobuf:"bytes,1,rep,name=result" json:"result,omitempty"`
}

func (m *Results) Reset()                    { *m = Results{} }
func (m *Results) String() string            { return proto.CompactTextString(m) }
func (*Results) ProtoMessage()               {}
func (*Results) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Results) GetResult() []*Result {
	if m != nil {
		return m.Result
	}
	return nil
}

// An Atom message is the network representation of a atom.Atom and the
// metadata needed to distribute it through the network.
type Atom struct {
	// Network data.
	To   *Address      `protobuf:"bytes,1,opt,name=to" json:"to,omitempty"`
	From *MultiAddress `protobuf:"bytes,2,opt,name=from" json:"from,omitempty"`
	// Secure data.
	Ledger    int64  `protobuf:"varint,3,opt,name=ledger" json:"ledger,omitempty"`
	Data      []byte `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	Signature []byte `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *Atom) Reset()                    { *m = Atom{} }
func (m *Atom) String() string            { return proto.CompactTextString(m) }
func (*Atom) ProtoMessage()               {}
func (*Atom) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *Atom) GetTo() *Address {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Atom) GetFrom() *MultiAddress {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Atom) GetLedger() int64 {
	if m != nil {
		return m.Ledger
	}
	return 0
}

func (m *Atom) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Atom) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*Address)(nil), "rpc.Address")
	proto.RegisterType((*MultiAddress)(nil), "rpc.MultiAddress")
	proto.RegisterType((*MultiAddresses)(nil), "rpc.MultiAddresses")
	proto.RegisterType((*Nothing)(nil), "rpc.Nothing")
	proto.RegisterType((*Query)(nil), "rpc.Query")
	proto.RegisterType((*OrderFragment)(nil), "rpc.OrderFragment")
	proto.RegisterType((*ResultFragment)(nil), "rpc.ResultFragment")
	proto.RegisterType((*Result)(nil), "rpc.Result")
	proto.RegisterType((*Results)(nil), "rpc.Results")
	proto.RegisterType((*Atom)(nil), "rpc.Atom")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SwarmNode service

type SwarmNodeClient interface {
	// Ping the connection and swap MultiAddresses.
	Ping(ctx context.Context, in *MultiAddress, opts ...grpc.CallOption) (*Nothing, error)
	// Find the MultiAddresses of peers closer to some target Node.
	QueryCloserPeers(ctx context.Context, in *Query, opts ...grpc.CallOption) (*MultiAddresses, error)
	// Find the MultiAddresses of peers closer to some target Node using a
	// frontier search.
	QueryCloserPeersOnFrontier(ctx context.Context, in *Query, opts ...grpc.CallOption) (SwarmNode_QueryCloserPeersOnFrontierClient, error)
}

type swarmNodeClient struct {
	cc *grpc.ClientConn
}

func NewSwarmNodeClient(cc *grpc.ClientConn) SwarmNodeClient {
	return &swarmNodeClient{cc}
}

func (c *swarmNodeClient) Ping(ctx context.Context, in *MultiAddress, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := grpc.Invoke(ctx, "/rpc.SwarmNode/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swarmNodeClient) QueryCloserPeers(ctx context.Context, in *Query, opts ...grpc.CallOption) (*MultiAddresses, error) {
	out := new(MultiAddresses)
	err := grpc.Invoke(ctx, "/rpc.SwarmNode/QueryCloserPeers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swarmNodeClient) QueryCloserPeersOnFrontier(ctx context.Context, in *Query, opts ...grpc.CallOption) (SwarmNode_QueryCloserPeersOnFrontierClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_SwarmNode_serviceDesc.Streams[0], c.cc, "/rpc.SwarmNode/QueryCloserPeersOnFrontier", opts...)
	if err != nil {
		return nil, err
	}
	x := &swarmNodeQueryCloserPeersOnFrontierClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SwarmNode_QueryCloserPeersOnFrontierClient interface {
	Recv() (*MultiAddress, error)
	grpc.ClientStream
}

type swarmNodeQueryCloserPeersOnFrontierClient struct {
	grpc.ClientStream
}

func (x *swarmNodeQueryCloserPeersOnFrontierClient) Recv() (*MultiAddress, error) {
	m := new(MultiAddress)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for SwarmNode service

type SwarmNodeServer interface {
	// Ping the connection and swap MultiAddresses.
	Ping(context.Context, *MultiAddress) (*Nothing, error)
	// Find the MultiAddresses of peers closer to some target Node.
	QueryCloserPeers(context.Context, *Query) (*MultiAddresses, error)
	// Find the MultiAddresses of peers closer to some target Node using a
	// frontier search.
	QueryCloserPeersOnFrontier(*Query, SwarmNode_QueryCloserPeersOnFrontierServer) error
}

func RegisterSwarmNodeServer(s *grpc.Server, srv SwarmNodeServer) {
	s.RegisterService(&_SwarmNode_serviceDesc, srv)
}

func _SwarmNode_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MultiAddress)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwarmNodeServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.SwarmNode/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwarmNodeServer).Ping(ctx, req.(*MultiAddress))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwarmNode_QueryCloserPeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwarmNodeServer).QueryCloserPeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.SwarmNode/QueryCloserPeers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwarmNodeServer).QueryCloserPeers(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwarmNode_QueryCloserPeersOnFrontier_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Query)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SwarmNodeServer).QueryCloserPeersOnFrontier(m, &swarmNodeQueryCloserPeersOnFrontierServer{stream})
}

type SwarmNode_QueryCloserPeersOnFrontierServer interface {
	Send(*MultiAddress) error
	grpc.ServerStream
}

type swarmNodeQueryCloserPeersOnFrontierServer struct {
	grpc.ServerStream
}

func (x *swarmNodeQueryCloserPeersOnFrontierServer) Send(m *MultiAddress) error {
	return x.ServerStream.SendMsg(m)
}

var _SwarmNode_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.SwarmNode",
	HandlerType: (*SwarmNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _SwarmNode_Ping_Handler,
		},
		{
			MethodName: "QueryCloserPeers",
			Handler:    _SwarmNode_QueryCloserPeers_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "QueryCloserPeersOnFrontier",
			Handler:       _SwarmNode_QueryCloserPeersOnFrontier_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rpc.proto",
}

// Client API for XingNode service

type XingNodeClient interface {
	// Send an OrderFragment to some target Node.
	SendOrderFragment(ctx context.Context, in *OrderFragment, opts ...grpc.CallOption) (*Nothing, error)
	// Send a ResultFragment to some target Node, where the ResultFragment is the
	// result of a computation on two OrderFragments.
	SendResultFragment(ctx context.Context, in *ResultFragment, opts ...grpc.CallOption) (*Nothing, error)
	// Get Result messages for successful order matches that have happened since
	// this procedure was last used, as well as new ones that occur.
	Notifications(ctx context.Context, in *Address, opts ...grpc.CallOption) (XingNode_NotificationsClient, error)
	// Get all relevant results that have been computed in the current epoch.
	GetResults(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Results, error)
}

type xingNodeClient struct {
	cc *grpc.ClientConn
}

func NewXingNodeClient(cc *grpc.ClientConn) XingNodeClient {
	return &xingNodeClient{cc}
}

func (c *xingNodeClient) SendOrderFragment(ctx context.Context, in *OrderFragment, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := grpc.Invoke(ctx, "/rpc.XingNode/SendOrderFragment", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *xingNodeClient) SendResultFragment(ctx context.Context, in *ResultFragment, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := grpc.Invoke(ctx, "/rpc.XingNode/SendResultFragment", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *xingNodeClient) Notifications(ctx context.Context, in *Address, opts ...grpc.CallOption) (XingNode_NotificationsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_XingNode_serviceDesc.Streams[0], c.cc, "/rpc.XingNode/Notifications", opts...)
	if err != nil {
		return nil, err
	}
	x := &xingNodeNotificationsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type XingNode_NotificationsClient interface {
	Recv() (*Result, error)
	grpc.ClientStream
}

type xingNodeNotificationsClient struct {
	grpc.ClientStream
}

func (x *xingNodeNotificationsClient) Recv() (*Result, error) {
	m := new(Result)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *xingNodeClient) GetResults(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Results, error) {
	out := new(Results)
	err := grpc.Invoke(ctx, "/rpc.XingNode/GetResults", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for XingNode service

type XingNodeServer interface {
	// Send an OrderFragment to some target Node.
	SendOrderFragment(context.Context, *OrderFragment) (*Nothing, error)
	// Send a ResultFragment to some target Node, where the ResultFragment is the
	// result of a computation on two OrderFragments.
	SendResultFragment(context.Context, *ResultFragment) (*Nothing, error)
	// Get Result messages for successful order matches that have happened since
	// this procedure was last used, as well as new ones that occur.
	Notifications(*Address, XingNode_NotificationsServer) error
	// Get all relevant results that have been computed in the current epoch.
	GetResults(context.Context, *Address) (*Results, error)
}

func RegisterXingNodeServer(s *grpc.Server, srv XingNodeServer) {
	s.RegisterService(&_XingNode_serviceDesc, srv)
}

func _XingNode_SendOrderFragment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderFragment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XingNodeServer).SendOrderFragment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.XingNode/SendOrderFragment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XingNodeServer).SendOrderFragment(ctx, req.(*OrderFragment))
	}
	return interceptor(ctx, in, info, handler)
}

func _XingNode_SendResultFragment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResultFragment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XingNodeServer).SendResultFragment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.XingNode/SendResultFragment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XingNodeServer).SendResultFragment(ctx, req.(*ResultFragment))
	}
	return interceptor(ctx, in, info, handler)
}

func _XingNode_Notifications_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Address)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(XingNodeServer).Notifications(m, &xingNodeNotificationsServer{stream})
}

type XingNode_NotificationsServer interface {
	Send(*Result) error
	grpc.ServerStream
}

type xingNodeNotificationsServer struct {
	grpc.ServerStream
}

func (x *xingNodeNotificationsServer) Send(m *Result) error {
	return x.ServerStream.SendMsg(m)
}

func _XingNode_GetResults_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Address)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XingNodeServer).GetResults(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.XingNode/GetResults",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XingNodeServer).GetResults(ctx, req.(*Address))
	}
	return interceptor(ctx, in, info, handler)
}

var _XingNode_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.XingNode",
	HandlerType: (*XingNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendOrderFragment",
			Handler:    _XingNode_SendOrderFragment_Handler,
		},
		{
			MethodName: "SendResultFragment",
			Handler:    _XingNode_SendResultFragment_Handler,
		},
		{
			MethodName: "GetResults",
			Handler:    _XingNode_GetResults_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Notifications",
			Handler:       _XingNode_Notifications_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rpc.proto",
}

// Client API for TerminalNode service

type TerminalNodeClient interface {
	// Send an Atom to some target Node.
	SendAtom(ctx context.Context, in *Atom, opts ...grpc.CallOption) (*Atom, error)
}

type terminalNodeClient struct {
	cc *grpc.ClientConn
}

func NewTerminalNodeClient(cc *grpc.ClientConn) TerminalNodeClient {
	return &terminalNodeClient{cc}
}

func (c *terminalNodeClient) SendAtom(ctx context.Context, in *Atom, opts ...grpc.CallOption) (*Atom, error) {
	out := new(Atom)
	err := grpc.Invoke(ctx, "/rpc.TerminalNode/SendAtom", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for TerminalNode service

type TerminalNodeServer interface {
	// Send an Atom to some target Node.
	SendAtom(context.Context, *Atom) (*Atom, error)
}

func RegisterTerminalNodeServer(s *grpc.Server, srv TerminalNodeServer) {
	s.RegisterService(&_TerminalNode_serviceDesc, srv)
}

func _TerminalNode_SendAtom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Atom)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerminalNodeServer).SendAtom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.TerminalNode/SendAtom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerminalNodeServer).SendAtom(ctx, req.(*Atom))
	}
	return interceptor(ctx, in, info, handler)
}

var _TerminalNode_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.TerminalNode",
	HandlerType: (*TerminalNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendAtom",
			Handler:    _TerminalNode_SendAtom_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}

func init() { proto.RegisterFile("rpc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 705 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0xcd, 0x6a, 0xdb, 0x4c,
	0x14, 0x45, 0x92, 0xff, 0x74, 0xad, 0x98, 0x2f, 0x93, 0x8f, 0x22, 0x4c, 0x08, 0x46, 0xe9, 0x8f,
	0x0b, 0xc5, 0x18, 0x87, 0xd2, 0x45, 0xe9, 0x22, 0x04, 0x52, 0xb2, 0x68, 0x92, 0x2a, 0xa1, 0x74,
	0xab, 0x58, 0x13, 0x67, 0x40, 0xd2, 0xb8, 0x33, 0x63, 0x5a, 0x3f, 0x48, 0xdf, 0xa1, 0xaf, 0xd1,
	0x57, 0xe8, 0x13, 0xf4, 0x29, 0xba, 0x2e, 0x73, 0x35, 0x8a, 0x24, 0xdb, 0x69, 0xb2, 0x68, 0x77,
	0x73, 0xcf, 0x3d, 0x77, 0xee, 0xdc, 0x73, 0x46, 0x23, 0x70, 0xc5, 0x7c, 0x3a, 0x9a, 0x0b, 0xae,
	0x38, 0x71, 0xc4, 0x7c, 0x1a, 0xec, 0x43, 0xfb, 0x30, 0x8e, 0x05, 0x95, 0x92, 0xf8, 0xd0, 0x8e,
	0xf2, 0xa5, 0x6f, 0x0d, 0xac, 0xa1, 0x1b, 0x16, 0x61, 0xf0, 0x18, 0xbc, 0x77, 0x8b, 0x44, 0xb1,
	0x82, 0xf9, 0x3f, 0x34, 0x53, 0x1d, 0x1b, 0x5e, 0x1e, 0x04, 0xaf, 0xa1, 0x57, 0x65, 0x51, 0x49,
	0x9e, 0x43, 0x0b, 0x53, 0x7a, 0x43, 0x67, 0xd8, 0x9d, 0x6c, 0x8f, 0x74, 0xf7, 0x2a, 0x29, 0x34,
	0x84, 0xc0, 0x85, 0xf6, 0x29, 0x57, 0x37, 0x2c, 0x9b, 0x05, 0x21, 0x34, 0xdf, 0x2f, 0xa8, 0x58,
	0x92, 0x27, 0xd0, 0xb8, 0x16, 0x3c, 0xc5, 0x2e, 0x1b, 0x8b, 0x31, 0x4d, 0x02, 0x68, 0x7e, 0xd2,
	0x7c, 0xdf, 0x46, 0x9e, 0x87, 0xbc, 0x82, 0x92, 0xa7, 0x82, 0x5f, 0x36, 0x6c, 0x9d, 0x89, 0x98,
	0x8a, 0x63, 0x11, 0xcd, 0x52, 0x9a, 0x29, 0xb2, 0x0b, 0xb6, 0xe2, 0x66, 0xeb, 0x7a, 0x89, 0xad,
	0xf8, 0x6d, 0x6b, 0xfb, 0xcf, 0xad, 0x7b, 0x60, 0xb3, 0xd8, 0x77, 0x06, 0xd6, 0xd0, 0x0b, 0x6d,
	0x16, 0x6b, 0x09, 0xb9, 0xee, 0x72, 0x12, 0xfb, 0x0d, 0x04, 0x8b, 0x90, 0xec, 0x82, 0x8b, 0xcb,
	0xcb, 0xe5, 0x9c, 0xfa, 0xcd, 0x81, 0x35, 0x74, 0xc2, 0x12, 0x20, 0x03, 0xe8, 0x62, 0x70, 0x1e,
	0x09, 0xa6, 0x96, 0x7e, 0x0b, 0xf3, 0x55, 0x88, 0x04, 0xe0, 0x5d, 0x4b, 0x75, 0xc4, 0x63, 0x7a,
	0x71, 0x13, 0x09, 0xea, 0xb7, 0x71, 0xfb, 0x1a, 0xa6, 0x39, 0x32, 0x8b, 0x4b, 0x4e, 0x27, 0xe7,
	0x54, 0x31, 0xb2, 0x07, 0x30, 0x17, 0x6c, 0x6a, 0x18, 0x2e, 0x32, 0x2a, 0x08, 0x79, 0x0a, 0xbd,
	0x34, 0xfa, 0xf2, 0x81, 0x27, 0x8b, 0xd4, 0x70, 0x00, 0x39, 0x2b, 0x28, 0xf2, 0x58, 0x56, 0xe5,
	0x75, 0x0d, 0xaf, 0x86, 0x06, 0xdf, 0x1d, 0xe8, 0x85, 0x54, 0x2e, 0x12, 0xf5, 0x6f, 0x95, 0xdf,
	0x03, 0xb8, 0x5a, 0x2c, 0xcf, 0x6a, 0xe2, 0x57, 0x10, 0xad, 0xb0, 0xa4, 0x49, 0x52, 0x10, 0x9a,
	0x48, 0xa8, 0x42, 0x64, 0x04, 0xa4, 0xe0, 0x17, 0x47, 0x3d, 0x89, 0xd1, 0x0a, 0x2f, 0xdc, 0x90,
	0x21, 0x63, 0xd8, 0xb9, 0x2d, 0xaf, 0x14, 0xe4, 0xc6, 0x6c, 0x4a, 0xad, 0x79, 0xd8, 0x79, 0x80,
	0x87, 0xee, 0xbd, 0x1e, 0xc2, 0x03, 0x3c, 0xec, 0x3e, 0xd0, 0x43, 0x6f, 0xa3, 0x87, 0x3f, 0x2d,
	0x68, 0xe5, 0x1e, 0x1a, 0xd9, 0xad, 0x3b, 0x64, 0xb7, 0xef, 0x93, 0xdd, 0x59, 0x97, 0xdd, 0x87,
	0xb6, 0x11, 0xa0, 0xf8, 0x64, 0x4c, 0xa8, 0x33, 0x66, 0x6c, 0x63, 0x57, 0x11, 0xea, 0xf7, 0x07,
	0xc7, 0x35, 0xee, 0xe4, 0x81, 0xfe, 0xc4, 0x6e, 0x07, 0x34, 0x36, 0x94, 0x00, 0x66, 0x8b, 0xb1,
	0x8c, 0xf2, 0x25, 0x10, 0x8c, 0xa0, 0x9d, 0x4f, 0x28, 0xc9, 0x3e, 0xb4, 0x04, 0x2e, 0xcd, 0xa3,
	0xd5, 0xc5, 0x2b, 0x98, 0x67, 0x43, 0x93, 0x0a, 0xbe, 0x5a, 0xd0, 0x38, 0x54, 0x3c, 0xfd, 0x3b,
	0x97, 0xf9, 0x11, 0xb4, 0x12, 0x1a, 0xcf, 0xa8, 0x40, 0x81, 0x9c, 0xd0, 0x44, 0x84, 0x40, 0x23,
	0x8e, 0x54, 0x64, 0x84, 0xc1, 0xb5, 0x9e, 0x43, 0xb2, 0x59, 0x16, 0xa9, 0x85, 0x28, 0x74, 0x29,
	0x81, 0xc9, 0x37, 0x0b, 0xdc, 0x8b, 0xcf, 0x91, 0x48, 0x4f, 0xb5, 0x4e, 0xcf, 0xa0, 0x71, 0xce,
	0xb2, 0x19, 0x59, 0x6f, 0xdc, 0xcf, 0xcf, 0x6a, 0x9e, 0x5c, 0x72, 0x00, 0xff, 0xe1, 0x93, 0x7b,
	0x94, 0x70, 0x49, 0xc5, 0x39, 0xa5, 0x42, 0x12, 0x40, 0x06, 0xc2, 0xfd, 0x9d, 0xb5, 0x0d, 0xa8,
	0x24, 0x6f, 0xa0, 0xbf, 0x5a, 0x74, 0x96, 0x1d, 0x0b, 0x9e, 0x29, 0x46, 0x45, 0xad, 0x7c, 0xbd,
	0xff, 0xd8, 0x9a, 0xfc, 0xb0, 0xa0, 0xf3, 0x91, 0x65, 0x33, 0x3c, 0xe9, 0x4b, 0xd8, 0xbe, 0xa0,
	0x59, 0x5c, 0x7f, 0xa2, 0x09, 0x96, 0xd5, 0xb0, 0x95, 0x73, 0xbf, 0x02, 0xa2, 0xcb, 0x56, 0x1e,
	0x98, 0x9d, 0x8a, 0x63, 0x77, 0x14, 0xbe, 0x80, 0xad, 0x53, 0xae, 0xd8, 0x35, 0x9b, 0x46, 0x8a,
	0xf1, 0x4c, 0x92, 0x9a, 0x77, 0xfd, 0xaa, 0xe7, 0x63, 0x8b, 0x0c, 0x01, 0xde, 0x52, 0x55, 0x5c,
	0x90, 0x3a, 0xd5, 0xab, 0x50, 0xe5, 0x64, 0x0c, 0xde, 0x25, 0x15, 0x29, 0xcb, 0xa2, 0x04, 0xe7,
	0x1a, 0x40, 0x47, 0x1f, 0x10, 0xaf, 0x8a, 0x9b, 0xd7, 0x29, 0x9e, 0xf6, 0xcb, 0xe5, 0x55, 0x0b,
	0x7f, 0xc6, 0x07, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x81, 0xd2, 0x8d, 0x24, 0x99, 0x07, 0x00,
	0x00,
}
