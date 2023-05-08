// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: podc/types/vrf.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	crypto "github.com/reapchain/reapchain-core/proto/podc/crypto"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type VrfSet struct {
	Height int64  `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Vrfs   []*Vrf `protobuf:"bytes,2,rep,name=vrfs,proto3" json:"vrfs,omitempty"`
}

func (m *VrfSet) Reset()         { *m = VrfSet{} }
func (m *VrfSet) String() string { return proto.CompactTextString(m) }
func (*VrfSet) ProtoMessage()    {}
func (*VrfSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_a2d47c3ad37760bc, []int{0}
}
func (m *VrfSet) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *VrfSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_VrfSet.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *VrfSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VrfSet.Merge(m, src)
}
func (m *VrfSet) XXX_Size() int {
	return m.Size()
}
func (m *VrfSet) XXX_DiscardUnknown() {
	xxx_messageInfo_VrfSet.DiscardUnknown(m)
}

var xxx_messageInfo_VrfSet proto.InternalMessageInfo

func (m *VrfSet) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *VrfSet) GetVrfs() []*Vrf {
	if m != nil {
		return m.Vrfs
	}
	return nil
}

type Vrf struct {
	Height                        int64            `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Timestamp                     time.Time        `protobuf:"bytes,2,opt,name=timestamp,proto3,stdtime" json:"timestamp"`
	SteeringMemberCandidatePubKey crypto.PublicKey `protobuf:"bytes,3,opt,name=steering_member_candidate_pub_key,json=steeringMemberCandidatePubKey,proto3" json:"steering_member_candidate_pub_key"`
	Value                         []byte           `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	Proof                         []byte           `protobuf:"bytes,5,opt,name=proof,proto3" json:"proof,omitempty"`
	VrfIndex                      int32            `protobuf:"varint,6,opt,name=vrf_index,json=vrfIndex,proto3" json:"vrf_index,omitempty"`
	Seed                          []byte           `protobuf:"bytes,7,opt,name=seed,proto3" json:"seed,omitempty"`
}

func (m *Vrf) Reset()         { *m = Vrf{} }
func (m *Vrf) String() string { return proto.CompactTextString(m) }
func (*Vrf) ProtoMessage()    {}
func (*Vrf) Descriptor() ([]byte, []int) {
	return fileDescriptor_a2d47c3ad37760bc, []int{1}
}
func (m *Vrf) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Vrf) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Vrf.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Vrf) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vrf.Merge(m, src)
}
func (m *Vrf) XXX_Size() int {
	return m.Size()
}
func (m *Vrf) XXX_DiscardUnknown() {
	xxx_messageInfo_Vrf.DiscardUnknown(m)
}

var xxx_messageInfo_Vrf proto.InternalMessageInfo

func (m *Vrf) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Vrf) GetTimestamp() time.Time {
	if m != nil {
		return m.Timestamp
	}
	return time.Time{}
}

func (m *Vrf) GetSteeringMemberCandidatePubKey() crypto.PublicKey {
	if m != nil {
		return m.SteeringMemberCandidatePubKey
	}
	return crypto.PublicKey{}
}

func (m *Vrf) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Vrf) GetProof() []byte {
	if m != nil {
		return m.Proof
	}
	return nil
}

func (m *Vrf) GetVrfIndex() int32 {
	if m != nil {
		return m.VrfIndex
	}
	return 0
}

func (m *Vrf) GetSeed() []byte {
	if m != nil {
		return m.Seed
	}
	return nil
}

func init() {
	proto.RegisterType((*VrfSet)(nil), "reapchain.types.VrfSet")
	proto.RegisterType((*Vrf)(nil), "reapchain.types.Vrf")
}

func init() { proto.RegisterFile("podc/types/vrf.proto", fileDescriptor_a2d47c3ad37760bc) }

var fileDescriptor_a2d47c3ad37760bc = []byte{
	// 403 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0xcd, 0x8e, 0xd3, 0x30,
	0x18, 0xac, 0xfb, 0xc7, 0xae, 0x17, 0x09, 0xc9, 0xaa, 0x56, 0x56, 0x57, 0xa4, 0x61, 0x4f, 0xb9,
	0x60, 0x4b, 0x0b, 0x4f, 0x50, 0x4e, 0xb0, 0x02, 0xad, 0x02, 0xea, 0x81, 0x4b, 0x94, 0x9f, 0xcf,
	0xa9, 0xd5, 0xa6, 0xb6, 0x1c, 0x27, 0x22, 0x6f, 0xb1, 0xef, 0xc1, 0x8b, 0xec, 0x71, 0x8f, 0x9c,
	0x00, 0xb5, 0x2f, 0x82, 0xe2, 0x34, 0x5b, 0x84, 0xc4, 0xed, 0x9b, 0xf9, 0x66, 0x26, 0x13, 0xdb,
	0x78, 0xa6, 0x55, 0x96, 0x72, 0xdb, 0x68, 0x28, 0x79, 0x6d, 0x04, 0xd3, 0x46, 0x59, 0x45, 0x5e,
	0x18, 0x88, 0x75, 0xba, 0x8e, 0xe5, 0x8e, 0xb9, 0xd5, 0x7c, 0x96, 0xab, 0x5c, 0xb9, 0x1d, 0x6f,
	0xa7, 0x4e, 0x36, 0xbf, 0x74, 0xe6, 0xd4, 0x34, 0xda, 0x2a, 0xbe, 0x81, 0xa6, 0x3c, 0xf2, 0x8b,
	0x5c, 0xa9, 0x7c, 0x0b, 0xdc, 0xa1, 0xa4, 0x12, 0xdc, 0xca, 0x02, 0x4a, 0x1b, 0x17, 0xba, 0x13,
	0x5c, 0x7f, 0xc0, 0xd3, 0x95, 0x11, 0x9f, 0xc1, 0x92, 0x4b, 0x3c, 0x5d, 0x83, 0xcc, 0xd7, 0x96,
	0x22, 0x1f, 0x05, 0xa3, 0xf0, 0x88, 0x48, 0x80, 0xc7, 0xb5, 0x11, 0x25, 0x1d, 0xfa, 0xa3, 0xe0,
	0xe2, 0x66, 0xc6, 0xfe, 0x29, 0xc4, 0x56, 0x46, 0x84, 0x4e, 0x71, 0xfd, 0x7d, 0x88, 0x47, 0x2b,
	0x23, 0xfe, 0x9b, 0xb4, 0xc4, 0xe7, 0x4f, 0x9f, 0xa7, 0x43, 0x1f, 0x05, 0x17, 0x37, 0x73, 0xd6,
	0x15, 0x64, 0x7d, 0x41, 0xf6, 0xa5, 0x57, 0x2c, 0xcf, 0x1e, 0x7e, 0x2e, 0x06, 0xf7, 0xbf, 0x16,
	0x28, 0x3c, 0xd9, 0xc8, 0x06, 0xbf, 0x2a, 0x2d, 0x80, 0x91, 0xbb, 0x3c, 0x2a, 0xa0, 0x48, 0xc0,
	0x44, 0x69, 0xbc, 0xcb, 0x64, 0x16, 0x5b, 0x88, 0x74, 0x95, 0x44, 0x1b, 0x68, 0xe8, 0xc8, 0x65,
	0x5f, 0xfd, 0x55, 0xb5, 0x3b, 0x19, 0x76, 0x57, 0x25, 0x5b, 0x99, 0xde, 0x42, 0xb3, 0x1c, 0xb7,
	0xe1, 0xe1, 0xcb, 0x3e, 0xeb, 0xa3, 0x8b, 0x7a, 0xd7, 0x27, 0xdd, 0x55, 0xc9, 0x2d, 0x34, 0x64,
	0x86, 0x27, 0x75, 0xbc, 0xad, 0x80, 0x8e, 0x7d, 0x14, 0x3c, 0x0f, 0x3b, 0xd0, 0xb2, 0xda, 0x28,
	0x25, 0xe8, 0xa4, 0x63, 0x1d, 0x20, 0x57, 0xf8, 0xbc, 0x36, 0x22, 0x92, 0xbb, 0x0c, 0xbe, 0xd1,
	0xa9, 0x8f, 0x82, 0x49, 0x78, 0x56, 0x1b, 0xf1, 0xbe, 0xc5, 0x84, 0xe0, 0x71, 0x09, 0x90, 0xd1,
	0x67, 0xce, 0xe1, 0xe6, 0xe5, 0xa7, 0x87, 0xbd, 0x87, 0x1e, 0xf7, 0x1e, 0xfa, 0xbd, 0xf7, 0xd0,
	0xfd, 0xc1, 0x1b, 0x3c, 0x1e, 0xbc, 0xc1, 0x8f, 0x83, 0x37, 0xf8, 0xfa, 0x36, 0x97, 0x76, 0x5d,
	0x25, 0x2c, 0x55, 0x05, 0x7f, 0xfa, 0x85, 0xd3, 0xf4, 0x3a, 0x55, 0xe6, 0x78, 0xa3, 0xfc, 0xf4,
	0x66, 0x92, 0xa9, 0x63, 0xde, 0xfc, 0x09, 0x00, 0x00, 0xff, 0xff, 0xf1, 0x9a, 0x0a, 0xa0, 0x48,
	0x02, 0x00, 0x00,
}

func (m *VrfSet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *VrfSet) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *VrfSet) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Vrfs) > 0 {
		for iNdEx := len(m.Vrfs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Vrfs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintVrf(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Height != 0 {
		i = encodeVarintVrf(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Vrf) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Vrf) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Vrf) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Seed) > 0 {
		i -= len(m.Seed)
		copy(dAtA[i:], m.Seed)
		i = encodeVarintVrf(dAtA, i, uint64(len(m.Seed)))
		i--
		dAtA[i] = 0x3a
	}
	if m.VrfIndex != 0 {
		i = encodeVarintVrf(dAtA, i, uint64(m.VrfIndex))
		i--
		dAtA[i] = 0x30
	}
	if len(m.Proof) > 0 {
		i -= len(m.Proof)
		copy(dAtA[i:], m.Proof)
		i = encodeVarintVrf(dAtA, i, uint64(len(m.Proof)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintVrf(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x22
	}
	{
		size, err := m.SteeringMemberCandidatePubKey.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintVrf(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	n2, err2 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.Timestamp, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.Timestamp):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintVrf(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x12
	if m.Height != 0 {
		i = encodeVarintVrf(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintVrf(dAtA []byte, offset int, v uint64) int {
	offset -= sovVrf(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *VrfSet) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovVrf(uint64(m.Height))
	}
	if len(m.Vrfs) > 0 {
		for _, e := range m.Vrfs {
			l = e.Size()
			n += 1 + l + sovVrf(uint64(l))
		}
	}
	return n
}

func (m *Vrf) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovVrf(uint64(m.Height))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.Timestamp)
	n += 1 + l + sovVrf(uint64(l))
	l = m.SteeringMemberCandidatePubKey.Size()
	n += 1 + l + sovVrf(uint64(l))
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovVrf(uint64(l))
	}
	l = len(m.Proof)
	if l > 0 {
		n += 1 + l + sovVrf(uint64(l))
	}
	if m.VrfIndex != 0 {
		n += 1 + sovVrf(uint64(m.VrfIndex))
	}
	l = len(m.Seed)
	if l > 0 {
		n += 1 + l + sovVrf(uint64(l))
	}
	return n
}

func sovVrf(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozVrf(x uint64) (n int) {
	return sovVrf(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *VrfSet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVrf
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
			return fmt.Errorf("proto: VrfSet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: VrfSet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrf
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vrfs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrf
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
				return ErrInvalidLengthVrf
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVrf
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vrfs = append(m.Vrfs, &Vrf{})
			if err := m.Vrfs[len(m.Vrfs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVrf(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVrf
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
func (m *Vrf) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVrf
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
			return fmt.Errorf("proto: Vrf: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Vrf: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrf
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrf
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
				return ErrInvalidLengthVrf
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVrf
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.Timestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SteeringMemberCandidatePubKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrf
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
				return ErrInvalidLengthVrf
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVrf
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SteeringMemberCandidatePubKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrf
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
				return ErrInvalidLengthVrf
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthVrf
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = append(m.Value[:0], dAtA[iNdEx:postIndex]...)
			if m.Value == nil {
				m.Value = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proof", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrf
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
				return ErrInvalidLengthVrf
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthVrf
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Proof = append(m.Proof[:0], dAtA[iNdEx:postIndex]...)
			if m.Proof == nil {
				m.Proof = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VrfIndex", wireType)
			}
			m.VrfIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrf
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VrfIndex |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Seed", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrf
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
				return ErrInvalidLengthVrf
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthVrf
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Seed = append(m.Seed[:0], dAtA[iNdEx:postIndex]...)
			if m.Seed == nil {
				m.Seed = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVrf(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVrf
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
func skipVrf(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowVrf
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
					return 0, ErrIntOverflowVrf
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
					return 0, ErrIntOverflowVrf
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
				return 0, ErrInvalidLengthVrf
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupVrf
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthVrf
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthVrf        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowVrf          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupVrf = fmt.Errorf("proto: unexpected end of group")
)
