// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: reapchain/statesync/types.proto

package statesync

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Message struct {
	// Types that are valid to be assigned to Sum:
	//	*Message_SnapshotsRequest
	//	*Message_SnapshotsResponse
	//	*Message_ChunkRequest
	//	*Message_ChunkResponse
	Sum isMessage_Sum `protobuf_oneof:"sum"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbd2da69750eb828, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Message.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return m.Size()
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

type isMessage_Sum interface {
	isMessage_Sum()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Message_SnapshotsRequest struct {
	SnapshotsRequest *SnapshotsRequest `protobuf:"bytes,1,opt,name=snapshots_request,json=snapshotsRequest,proto3,oneof" json:"snapshots_request,omitempty"`
}
type Message_SnapshotsResponse struct {
	SnapshotsResponse *SnapshotsResponse `protobuf:"bytes,2,opt,name=snapshots_response,json=snapshotsResponse,proto3,oneof" json:"snapshots_response,omitempty"`
}
type Message_ChunkRequest struct {
	ChunkRequest *ChunkRequest `protobuf:"bytes,3,opt,name=chunk_request,json=chunkRequest,proto3,oneof" json:"chunk_request,omitempty"`
}
type Message_ChunkResponse struct {
	ChunkResponse *ChunkResponse `protobuf:"bytes,4,opt,name=chunk_response,json=chunkResponse,proto3,oneof" json:"chunk_response,omitempty"`
}

func (*Message_SnapshotsRequest) isMessage_Sum()  {}
func (*Message_SnapshotsResponse) isMessage_Sum() {}
func (*Message_ChunkRequest) isMessage_Sum()      {}
func (*Message_ChunkResponse) isMessage_Sum()     {}

func (m *Message) GetSum() isMessage_Sum {
	if m != nil {
		return m.Sum
	}
	return nil
}

func (m *Message) GetSnapshotsRequest() *SnapshotsRequest {
	if x, ok := m.GetSum().(*Message_SnapshotsRequest); ok {
		return x.SnapshotsRequest
	}
	return nil
}

func (m *Message) GetSnapshotsResponse() *SnapshotsResponse {
	if x, ok := m.GetSum().(*Message_SnapshotsResponse); ok {
		return x.SnapshotsResponse
	}
	return nil
}

func (m *Message) GetChunkRequest() *ChunkRequest {
	if x, ok := m.GetSum().(*Message_ChunkRequest); ok {
		return x.ChunkRequest
	}
	return nil
}

func (m *Message) GetChunkResponse() *ChunkResponse {
	if x, ok := m.GetSum().(*Message_ChunkResponse); ok {
		return x.ChunkResponse
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Message) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Message_SnapshotsRequest)(nil),
		(*Message_SnapshotsResponse)(nil),
		(*Message_ChunkRequest)(nil),
		(*Message_ChunkResponse)(nil),
	}
}

type SnapshotsRequest struct {
}

func (m *SnapshotsRequest) Reset()         { *m = SnapshotsRequest{} }
func (m *SnapshotsRequest) String() string { return proto.CompactTextString(m) }
func (*SnapshotsRequest) ProtoMessage()    {}
func (*SnapshotsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbd2da69750eb828, []int{1}
}
func (m *SnapshotsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SnapshotsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SnapshotsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SnapshotsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotsRequest.Merge(m, src)
}
func (m *SnapshotsRequest) XXX_Size() int {
	return m.Size()
}
func (m *SnapshotsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotsRequest proto.InternalMessageInfo

type SnapshotsResponse struct {
	Height   uint64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Format   uint32 `protobuf:"varint,2,opt,name=format,proto3" json:"format,omitempty"`
	Chunks   uint32 `protobuf:"varint,3,opt,name=chunks,proto3" json:"chunks,omitempty"`
	Hash     []byte `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
	Metadata []byte `protobuf:"bytes,5,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (m *SnapshotsResponse) Reset()         { *m = SnapshotsResponse{} }
func (m *SnapshotsResponse) String() string { return proto.CompactTextString(m) }
func (*SnapshotsResponse) ProtoMessage()    {}
func (*SnapshotsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbd2da69750eb828, []int{2}
}
func (m *SnapshotsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SnapshotsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SnapshotsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SnapshotsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotsResponse.Merge(m, src)
}
func (m *SnapshotsResponse) XXX_Size() int {
	return m.Size()
}
func (m *SnapshotsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotsResponse proto.InternalMessageInfo

func (m *SnapshotsResponse) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *SnapshotsResponse) GetFormat() uint32 {
	if m != nil {
		return m.Format
	}
	return 0
}

func (m *SnapshotsResponse) GetChunks() uint32 {
	if m != nil {
		return m.Chunks
	}
	return 0
}

func (m *SnapshotsResponse) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *SnapshotsResponse) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type ChunkRequest struct {
	Height uint64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Format uint32 `protobuf:"varint,2,opt,name=format,proto3" json:"format,omitempty"`
	Index  uint32 `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
}

func (m *ChunkRequest) Reset()         { *m = ChunkRequest{} }
func (m *ChunkRequest) String() string { return proto.CompactTextString(m) }
func (*ChunkRequest) ProtoMessage()    {}
func (*ChunkRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbd2da69750eb828, []int{3}
}
func (m *ChunkRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChunkRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChunkRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChunkRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChunkRequest.Merge(m, src)
}
func (m *ChunkRequest) XXX_Size() int {
	return m.Size()
}
func (m *ChunkRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ChunkRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ChunkRequest proto.InternalMessageInfo

func (m *ChunkRequest) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *ChunkRequest) GetFormat() uint32 {
	if m != nil {
		return m.Format
	}
	return 0
}

func (m *ChunkRequest) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type ChunkResponse struct {
	Height  uint64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Format  uint32 `protobuf:"varint,2,opt,name=format,proto3" json:"format,omitempty"`
	Index   uint32 `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
	Chunk   []byte `protobuf:"bytes,4,opt,name=chunk,proto3" json:"chunk,omitempty"`
	Missing bool   `protobuf:"varint,5,opt,name=missing,proto3" json:"missing,omitempty"`
}

func (m *ChunkResponse) Reset()         { *m = ChunkResponse{} }
func (m *ChunkResponse) String() string { return proto.CompactTextString(m) }
func (*ChunkResponse) ProtoMessage()    {}
func (*ChunkResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbd2da69750eb828, []int{4}
}
func (m *ChunkResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChunkResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChunkResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChunkResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChunkResponse.Merge(m, src)
}
func (m *ChunkResponse) XXX_Size() int {
	return m.Size()
}
func (m *ChunkResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ChunkResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ChunkResponse proto.InternalMessageInfo

func (m *ChunkResponse) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *ChunkResponse) GetFormat() uint32 {
	if m != nil {
		return m.Format
	}
	return 0
}

func (m *ChunkResponse) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *ChunkResponse) GetChunk() []byte {
	if m != nil {
		return m.Chunk
	}
	return nil
}

func (m *ChunkResponse) GetMissing() bool {
	if m != nil {
		return m.Missing
	}
	return false
}

func init() {
	proto.RegisterType((*Message)(nil), "reapchain.statesync.Message")
	proto.RegisterType((*SnapshotsRequest)(nil), "reapchain.statesync.SnapshotsRequest")
	proto.RegisterType((*SnapshotsResponse)(nil), "reapchain.statesync.SnapshotsResponse")
	proto.RegisterType((*ChunkRequest)(nil), "reapchain.statesync.ChunkRequest")
	proto.RegisterType((*ChunkResponse)(nil), "reapchain.statesync.ChunkResponse")
}

func init() { proto.RegisterFile("reapchain/statesync/types.proto", fileDescriptor_dbd2da69750eb828) }

var fileDescriptor_dbd2da69750eb828 = []byte{
	// 403 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x53, 0xcd, 0x6a, 0xdb, 0x40,
	0x18, 0x94, 0xfc, 0xcf, 0x57, 0xab, 0xd8, 0xdb, 0x52, 0x44, 0x0f, 0x6a, 0x2b, 0x68, 0xe9, 0xa5,
	0x12, 0xb4, 0x2f, 0x50, 0xdc, 0x8b, 0xa1, 0xf4, 0xb2, 0x35, 0x04, 0x7c, 0x09, 0x6b, 0x79, 0x23,
	0x89, 0x58, 0x2b, 0x45, 0xdf, 0x0a, 0xe2, 0x07, 0xc8, 0x29, 0x97, 0x3c, 0x56, 0x8e, 0x3e, 0x85,
	0x1c, 0x83, 0xfd, 0x22, 0x41, 0x2b, 0x59, 0x56, 0x1c, 0x93, 0x10, 0xc8, 0xcd, 0x33, 0x1e, 0x8d,
	0x66, 0x46, 0x7c, 0xf0, 0x29, 0xe5, 0x2c, 0xf1, 0x02, 0x16, 0x0a, 0x17, 0x25, 0x93, 0x1c, 0x97,
	0xc2, 0x73, 0xe5, 0x32, 0xe1, 0xe8, 0x24, 0x69, 0x2c, 0x63, 0xf2, 0xae, 0x12, 0x38, 0x95, 0xc0,
	0xbe, 0x69, 0x40, 0xf7, 0x1f, 0x47, 0x64, 0x3e, 0x27, 0x13, 0x18, 0xa2, 0x60, 0x09, 0x06, 0xb1,
	0xc4, 0xe3, 0x94, 0x9f, 0x65, 0x1c, 0xa5, 0xa9, 0x7f, 0xd6, 0xbf, 0xbf, 0xf9, 0xf9, 0xd5, 0x39,
	0xf0, 0xb0, 0xf3, 0x7f, 0xab, 0xa6, 0x85, 0x78, 0xac, 0xd1, 0x01, 0xee, 0x71, 0xe4, 0x08, 0x48,
	0xdd, 0x15, 0x93, 0x58, 0x20, 0x37, 0x1b, 0xca, 0xf6, 0xdb, 0x73, 0xb6, 0x85, 0x7a, 0xac, 0xd1,
	0x21, 0xee, 0x93, 0x64, 0x0c, 0x86, 0x17, 0x64, 0xe2, 0xb4, 0x8a, 0xda, 0x54, 0x9e, 0x5f, 0x0e,
	0x7a, 0xfe, 0xc9, 0x95, 0xbb, 0x98, 0x7d, 0xaf, 0x86, 0xc9, 0x5f, 0x78, 0xbb, 0x75, 0x2a, 0xe3,
	0xb5, 0x94, 0x95, 0xfd, 0x94, 0x55, 0x15, 0xcd, 0xf0, 0xea, 0xc4, 0xa8, 0x0d, 0x4d, 0xcc, 0x22,
	0x9b, 0xc0, 0x60, 0x7f, 0x1e, 0xfb, 0x52, 0x87, 0xe1, 0xa3, 0x72, 0xe4, 0x03, 0x74, 0x02, 0x1e,
	0xfa, 0x41, 0xb1, 0x75, 0x8b, 0x96, 0x28, 0xe7, 0x4f, 0xe2, 0x34, 0x62, 0x52, 0x8d, 0x65, 0xd0,
	0x12, 0xe5, 0xbc, 0x7a, 0x23, 0xaa, 0xc2, 0x06, 0x2d, 0x11, 0x21, 0xd0, 0x0a, 0x18, 0x06, 0x2a,
	0x7b, 0x9f, 0xaa, 0xdf, 0xe4, 0x23, 0xf4, 0x22, 0x2e, 0xd9, 0x9c, 0x49, 0x66, 0xb6, 0x15, 0x5f,
	0x61, 0x7b, 0x02, 0xfd, 0xfa, 0x2a, 0x2f, 0xce, 0xf1, 0x1e, 0xda, 0xa1, 0x98, 0xf3, 0xf3, 0x32,
	0x46, 0x01, 0xec, 0x0b, 0x1d, 0x8c, 0x07, 0x0b, 0xbd, 0x8e, 0x6f, 0xce, 0xaa, 0x9e, 0x65, 0xbd,
	0x02, 0x10, 0x13, 0xba, 0x51, 0x88, 0x18, 0x0a, 0x5f, 0xd5, 0xeb, 0xd1, 0x2d, 0x1c, 0x4d, 0xaf,
	0xd7, 0x96, 0xbe, 0x5a, 0x5b, 0xfa, 0xdd, 0xda, 0xd2, 0xaf, 0x36, 0x96, 0xb6, 0xda, 0x58, 0xda,
	0xed, 0xc6, 0xd2, 0xa6, 0xbf, 0xfd, 0x50, 0x2e, 0xd8, 0x4c, 0x7d, 0xdb, 0x84, 0x2d, 0x1d, 0xc1,
	0xa5, 0x8b, 0x99, 0x87, 0x3f, 0x16, 0x6c, 0xe6, 0xba, 0xbb, 0x2b, 0x52, 0x57, 0xe3, 0x1e, 0xb8,
	0xaa, 0x59, 0x47, 0xfd, 0xf5, 0xeb, 0x3e, 0x00, 0x00, 0xff, 0xff, 0xf7, 0xe0, 0xe5, 0xe3, 0x73,
	0x03, 0x00, 0x00,
}

func (m *Message) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Sum != nil {
		{
			size := m.Sum.Size()
			i -= size
			if _, err := m.Sum.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *Message_SnapshotsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message_SnapshotsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.SnapshotsRequest != nil {
		{
			size, err := m.SnapshotsRequest.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *Message_SnapshotsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message_SnapshotsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.SnapshotsResponse != nil {
		{
			size, err := m.SnapshotsResponse.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func (m *Message_ChunkRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message_ChunkRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.ChunkRequest != nil {
		{
			size, err := m.ChunkRequest.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	return len(dAtA) - i, nil
}
func (m *Message_ChunkResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message_ChunkResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.ChunkResponse != nil {
		{
			size, err := m.ChunkResponse.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	return len(dAtA) - i, nil
}
func (m *SnapshotsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SnapshotsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SnapshotsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *SnapshotsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SnapshotsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SnapshotsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Metadata) > 0 {
		i -= len(m.Metadata)
		copy(dAtA[i:], m.Metadata)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Metadata)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x22
	}
	if m.Chunks != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Chunks))
		i--
		dAtA[i] = 0x18
	}
	if m.Format != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Format))
		i--
		dAtA[i] = 0x10
	}
	if m.Height != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ChunkRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChunkRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChunkRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Index != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x18
	}
	if m.Format != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Format))
		i--
		dAtA[i] = 0x10
	}
	if m.Height != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ChunkResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChunkResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChunkResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Missing {
		i--
		if m.Missing {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if len(m.Chunk) > 0 {
		i -= len(m.Chunk)
		copy(dAtA[i:], m.Chunk)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Chunk)))
		i--
		dAtA[i] = 0x22
	}
	if m.Index != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x18
	}
	if m.Format != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Format))
		i--
		dAtA[i] = 0x10
	}
	if m.Height != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Message) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Sum != nil {
		n += m.Sum.Size()
	}
	return n
}

func (m *Message_SnapshotsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SnapshotsRequest != nil {
		l = m.SnapshotsRequest.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}
func (m *Message_SnapshotsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SnapshotsResponse != nil {
		l = m.SnapshotsResponse.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}
func (m *Message_ChunkRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ChunkRequest != nil {
		l = m.ChunkRequest.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}
func (m *Message_ChunkResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ChunkResponse != nil {
		l = m.ChunkResponse.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}
func (m *SnapshotsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *SnapshotsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovTypes(uint64(m.Height))
	}
	if m.Format != 0 {
		n += 1 + sovTypes(uint64(m.Format))
	}
	if m.Chunks != 0 {
		n += 1 + sovTypes(uint64(m.Chunks))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Metadata)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *ChunkRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovTypes(uint64(m.Height))
	}
	if m.Format != 0 {
		n += 1 + sovTypes(uint64(m.Format))
	}
	if m.Index != 0 {
		n += 1 + sovTypes(uint64(m.Index))
	}
	return n
}

func (m *ChunkResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovTypes(uint64(m.Height))
	}
	if m.Format != 0 {
		n += 1 + sovTypes(uint64(m.Format))
	}
	if m.Index != 0 {
		n += 1 + sovTypes(uint64(m.Index))
	}
	l = len(m.Chunk)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.Missing {
		n += 2
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Message) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Message: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Message: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SnapshotsRequest", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &SnapshotsRequest{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Sum = &Message_SnapshotsRequest{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SnapshotsResponse", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &SnapshotsResponse{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Sum = &Message_SnapshotsResponse{v}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChunkRequest", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &ChunkRequest{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Sum = &Message_ChunkRequest{v}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChunkResponse", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &ChunkResponse{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Sum = &Message_ChunkResponse{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SnapshotsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SnapshotsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SnapshotsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SnapshotsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SnapshotsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SnapshotsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Format", wireType)
			}
			m.Format = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Format |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chunks", wireType)
			}
			m.Chunks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Chunks |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append(m.Hash[:0], dAtA[iNdEx:postIndex]...)
			if m.Hash == nil {
				m.Hash = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metadata", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Metadata = append(m.Metadata[:0], dAtA[iNdEx:postIndex]...)
			if m.Metadata == nil {
				m.Metadata = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ChunkRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ChunkRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChunkRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Format", wireType)
			}
			m.Format = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Format |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ChunkResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ChunkResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChunkResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Format", wireType)
			}
			m.Format = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Format |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chunk", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Chunk = append(m.Chunk[:0], dAtA[iNdEx:postIndex]...)
			if m.Chunk == nil {
				m.Chunk = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Missing", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Missing = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
