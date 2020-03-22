// Code generated by protoc-gen-go. DO NOT EDIT.
// source: PeerResponse.proto

package protobuf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PeerResponse struct {
	Peers                []*PeerResponse_Peer `protobuf:"bytes,1,rep,name=Peers,proto3" json:"Peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *PeerResponse) Reset()         { *m = PeerResponse{} }
func (m *PeerResponse) String() string { return proto.CompactTextString(m) }
func (*PeerResponse) ProtoMessage()    {}
func (*PeerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c15da134fa15a697, []int{0}
}

func (m *PeerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerResponse.Unmarshal(m, b)
}
func (m *PeerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerResponse.Marshal(b, m, deterministic)
}
func (m *PeerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerResponse.Merge(m, src)
}
func (m *PeerResponse) XXX_Size() int {
	return xxx_messageInfo_PeerResponse.Size(m)
}
func (m *PeerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PeerResponse proto.InternalMessageInfo

func (m *PeerResponse) GetPeers() []*PeerResponse_Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}

type PeerResponse_Peer struct {
	Ip                   []byte   `protobuf:"bytes,1,opt,name=Ip,proto3" json:"Ip,omitempty"`
	Port                 uint32   `protobuf:"varint,2,opt,name=Port,proto3" json:"Port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PeerResponse_Peer) Reset()         { *m = PeerResponse_Peer{} }
func (m *PeerResponse_Peer) String() string { return proto.CompactTextString(m) }
func (*PeerResponse_Peer) ProtoMessage()    {}
func (*PeerResponse_Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_c15da134fa15a697, []int{0, 0}
}

func (m *PeerResponse_Peer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerResponse_Peer.Unmarshal(m, b)
}
func (m *PeerResponse_Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerResponse_Peer.Marshal(b, m, deterministic)
}
func (m *PeerResponse_Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerResponse_Peer.Merge(m, src)
}
func (m *PeerResponse_Peer) XXX_Size() int {
	return xxx_messageInfo_PeerResponse_Peer.Size(m)
}
func (m *PeerResponse_Peer) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerResponse_Peer.DiscardUnknown(m)
}

var xxx_messageInfo_PeerResponse_Peer proto.InternalMessageInfo

func (m *PeerResponse_Peer) GetIp() []byte {
	if m != nil {
		return m.Ip
	}
	return nil
}

func (m *PeerResponse_Peer) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func init() {
	proto.RegisterType((*PeerResponse)(nil), "DeMed.PeerResponse")
	proto.RegisterType((*PeerResponse_Peer)(nil), "DeMed.PeerResponse.Peer")
}

func init() {
	proto.RegisterFile("PeerResponse.proto", fileDescriptor_c15da134fa15a697)
}

var fileDescriptor_c15da134fa15a697 = []byte{
	// 134 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x0a, 0x48, 0x4d, 0x2d,
	0x0a, 0x4a, 0x2d, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x75, 0x49, 0xf5, 0x4d, 0x4d, 0x51, 0xca, 0xe2, 0xe2, 0x41, 0x96, 0x14, 0xd2, 0xe3, 0x62, 0x05,
	0xf1, 0x8b, 0x25, 0x18, 0x15, 0x98, 0x35, 0xb8, 0x8d, 0x24, 0xf4, 0xc0, 0xca, 0xf4, 0x50, 0x0c,
	0x00, 0x73, 0x20, 0xca, 0xa4, 0xb4, 0xb8, 0x58, 0x40, 0x0c, 0x21, 0x3e, 0x2e, 0x26, 0xcf, 0x02,
	0x09, 0x46, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x26, 0xcf, 0x02, 0x21, 0x21, 0x2e, 0x96, 0x80, 0xfc,
	0xa2, 0x12, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xde, 0x20, 0x30, 0xdb, 0x89, 0x2b, 0x8a, 0x03, 0x6c,
	0x77, 0x52, 0x69, 0x5a, 0x12, 0x1b, 0x98, 0x65, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xd6, 0x45,
	0xec, 0x29, 0x9b, 0x00, 0x00, 0x00,
}
