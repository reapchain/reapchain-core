// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: reapchain/types/steering_member_candidate.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	crypto "gitlab.reappay.net/reapchain/reapchain-core/proto/reapchain/crypto"
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

type SteeringMemberCandidateSet struct {
	SteeringMemberCandidates []*SteeringMemberCandidate `protobuf:"bytes,1,rep,name=steering_member_candidates,json=steeringMemberCandidates,proto3" json:"steering_member_candidates,omitempty"`
}

func (m *SteeringMemberCandidateSet) Reset()         { *m = SteeringMemberCandidateSet{} }
func (m *SteeringMemberCandidateSet) String() string { return proto.CompactTextString(m) }
func (*SteeringMemberCandidateSet) ProtoMessage()    {}
func (*SteeringMemberCandidateSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_56738ebc1a97ea7b, []int{0}
}
func (m *SteeringMemberCandidateSet) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SteeringMemberCandidateSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SteeringMemberCandidateSet.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SteeringMemberCandidateSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SteeringMemberCandidateSet.Merge(m, src)
}
func (m *SteeringMemberCandidateSet) XXX_Size() int {
	return m.Size()
}
func (m *SteeringMemberCandidateSet) XXX_DiscardUnknown() {
	xxx_messageInfo_SteeringMemberCandidateSet.DiscardUnknown(m)
}

var xxx_messageInfo_SteeringMemberCandidateSet proto.InternalMessageInfo

func (m *SteeringMemberCandidateSet) GetSteeringMemberCandidates() []*SteeringMemberCandidate {
	if m != nil {
		return m.SteeringMemberCandidates
	}
	return nil
}

type SteeringMemberCandidate struct {
	Address []byte           `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	PubKey  crypto.PublicKey `protobuf:"bytes,2,opt,name=pub_key,json=pubKey,proto3" json:"pub_key"`
}

func (m *SteeringMemberCandidate) Reset()         { *m = SteeringMemberCandidate{} }
func (m *SteeringMemberCandidate) String() string { return proto.CompactTextString(m) }
func (*SteeringMemberCandidate) ProtoMessage()    {}
func (*SteeringMemberCandidate) Descriptor() ([]byte, []int) {
	return fileDescriptor_56738ebc1a97ea7b, []int{1}
}
func (m *SteeringMemberCandidate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SteeringMemberCandidate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SteeringMemberCandidate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SteeringMemberCandidate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SteeringMemberCandidate.Merge(m, src)
}
func (m *SteeringMemberCandidate) XXX_Size() int {
	return m.Size()
}
func (m *SteeringMemberCandidate) XXX_DiscardUnknown() {
	xxx_messageInfo_SteeringMemberCandidate.DiscardUnknown(m)
}

var xxx_messageInfo_SteeringMemberCandidate proto.InternalMessageInfo

func (m *SteeringMemberCandidate) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *SteeringMemberCandidate) GetPubKey() crypto.PublicKey {
	if m != nil {
		return m.PubKey
	}
	return crypto.PublicKey{}
}

type SimpleSteeringMemberCandidate struct {
	PubKey *crypto.PublicKey `protobuf:"bytes,1,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
}

func (m *SimpleSteeringMemberCandidate) Reset()         { *m = SimpleSteeringMemberCandidate{} }
func (m *SimpleSteeringMemberCandidate) String() string { return proto.CompactTextString(m) }
func (*SimpleSteeringMemberCandidate) ProtoMessage()    {}
func (*SimpleSteeringMemberCandidate) Descriptor() ([]byte, []int) {
	return fileDescriptor_56738ebc1a97ea7b, []int{2}
}
func (m *SimpleSteeringMemberCandidate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SimpleSteeringMemberCandidate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SimpleSteeringMemberCandidate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SimpleSteeringMemberCandidate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SimpleSteeringMemberCandidate.Merge(m, src)
}
func (m *SimpleSteeringMemberCandidate) XXX_Size() int {
	return m.Size()
}
func (m *SimpleSteeringMemberCandidate) XXX_DiscardUnknown() {
	xxx_messageInfo_SimpleSteeringMemberCandidate.DiscardUnknown(m)
}

var xxx_messageInfo_SimpleSteeringMemberCandidate proto.InternalMessageInfo

func (m *SimpleSteeringMemberCandidate) GetPubKey() *crypto.PublicKey {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func init() {
	proto.RegisterType((*SteeringMemberCandidateSet)(nil), "reapchain.types.SteeringMemberCandidateSet")
	proto.RegisterType((*SteeringMemberCandidate)(nil), "reapchain.types.SteeringMemberCandidate")
	proto.RegisterType((*SimpleSteeringMemberCandidate)(nil), "reapchain.types.SimpleSteeringMemberCandidate")
}

func init() {
	proto.RegisterFile("reapchain/types/steering_member_candidate.proto", fileDescriptor_56738ebc1a97ea7b)
}

var fileDescriptor_56738ebc1a97ea7b = []byte{
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0x41, 0x4b, 0xc3, 0x30,
	0x1c, 0xc5, 0x1b, 0x95, 0x0d, 0x32, 0x41, 0x28, 0x82, 0xa5, 0xc3, 0x3a, 0x76, 0xea, 0xc5, 0x06,
	0xa6, 0x27, 0xf1, 0x34, 0x8f, 0x43, 0x90, 0x0d, 0x3d, 0x78, 0x19, 0x49, 0xfb, 0xb7, 0x86, 0x75,
	0x4d, 0x48, 0xd2, 0x43, 0xee, 0x7e, 0x00, 0x3f, 0xd6, 0x8e, 0x3b, 0x7a, 0x12, 0xd9, 0xbe, 0x88,
	0xac, 0xb5, 0x6e, 0x0c, 0x8b, 0xb7, 0x24, 0x2f, 0xef, 0xf7, 0x5e, 0xf8, 0x07, 0x13, 0x05, 0x54,
	0xc6, 0xaf, 0x94, 0xe7, 0xc4, 0x58, 0x09, 0x9a, 0x68, 0x03, 0xa0, 0x78, 0x9e, 0x4e, 0xe7, 0x30,
	0x67, 0xa0, 0xa6, 0x31, 0xcd, 0x13, 0x9e, 0x50, 0x03, 0x91, 0x54, 0xc2, 0x08, 0xf7, 0xe4, 0xd7,
	0x10, 0x95, 0x06, 0xff, 0x34, 0x15, 0xa9, 0x28, 0x35, 0xb2, 0x59, 0x55, 0xd7, 0xfc, 0xee, 0x96,
	0x1b, 0x2b, 0x2b, 0x8d, 0x20, 0x33, 0xb0, 0xba, 0x12, 0xfb, 0x6f, 0x08, 0xfb, 0x93, 0x9f, 0x9c,
	0xfb, 0x32, 0xe6, 0xae, 0x4e, 0x99, 0x80, 0x71, 0x5f, 0xb0, 0xdf, 0xd8, 0x42, 0x7b, 0xa8, 0x77,
	0x18, 0x76, 0x06, 0x61, 0xb4, 0xd7, 0x23, 0x6a, 0x00, 0x8e, 0x3d, 0xfd, 0xb7, 0xa0, 0xfb, 0x02,
	0x9f, 0x35, 0x98, 0x5c, 0x0f, 0xb7, 0x69, 0x92, 0x28, 0xd0, 0x9b, 0x3c, 0x14, 0x1e, 0x8f, 0xeb,
	0xad, 0x7b, 0x83, 0xdb, 0xb2, 0x60, 0xd3, 0x19, 0x58, 0xef, 0xa0, 0x87, 0xc2, 0xce, 0xa0, 0xbb,
	0xd3, 0xa4, 0x7a, 0x6a, 0xf4, 0x50, 0xb0, 0x8c, 0xc7, 0x23, 0xb0, 0xc3, 0xa3, 0xc5, 0xe7, 0x85,
	0x33, 0x6e, 0xc9, 0x82, 0x8d, 0xc0, 0xf6, 0x1f, 0xf1, 0xf9, 0x84, 0xcf, 0x65, 0x06, 0x4d, 0xb1,
	0xd7, 0x5b, 0x38, 0xfa, 0x17, 0x5e, 0x63, 0x87, 0x4f, 0x8b, 0x55, 0x80, 0x96, 0xab, 0x00, 0x7d,
	0xad, 0x02, 0xf4, 0xbe, 0x0e, 0x9c, 0xe5, 0x3a, 0x70, 0x3e, 0xd6, 0x81, 0xf3, 0x7c, 0x9b, 0x72,
	0x93, 0x51, 0x56, 0x42, 0x24, 0xb5, 0x51, 0x0e, 0x86, 0xe8, 0x22, 0xd6, 0x97, 0x19, 0x65, 0x64,
	0x67, 0xf4, 0xd5, 0xf8, 0xf6, 0xbe, 0x02, 0x6b, 0x95, 0xc7, 0x57, 0xdf, 0x01, 0x00, 0x00, 0xff,
	0xff, 0xaf, 0xb3, 0xa2, 0x86, 0x24, 0x02, 0x00, 0x00,
}

func (m *SteeringMemberCandidateSet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SteeringMemberCandidateSet) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SteeringMemberCandidateSet) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SteeringMemberCandidates) > 0 {
		for iNdEx := len(m.SteeringMemberCandidates) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SteeringMemberCandidates[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSteeringMemberCandidate(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *SteeringMemberCandidate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SteeringMemberCandidate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SteeringMemberCandidate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.PubKey.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintSteeringMemberCandidate(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintSteeringMemberCandidate(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SimpleSteeringMemberCandidate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SimpleSteeringMemberCandidate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SimpleSteeringMemberCandidate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.PubKey != nil {
		{
			size, err := m.PubKey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSteeringMemberCandidate(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSteeringMemberCandidate(dAtA []byte, offset int, v uint64) int {
	offset -= sovSteeringMemberCandidate(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SteeringMemberCandidateSet) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.SteeringMemberCandidates) > 0 {
		for _, e := range m.SteeringMemberCandidates {
			l = e.Size()
			n += 1 + l + sovSteeringMemberCandidate(uint64(l))
		}
	}
	return n
}

func (m *SteeringMemberCandidate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovSteeringMemberCandidate(uint64(l))
	}
	l = m.PubKey.Size()
	n += 1 + l + sovSteeringMemberCandidate(uint64(l))
	return n
}

func (m *SimpleSteeringMemberCandidate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PubKey != nil {
		l = m.PubKey.Size()
		n += 1 + l + sovSteeringMemberCandidate(uint64(l))
	}
	return n
}

func sovSteeringMemberCandidate(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSteeringMemberCandidate(x uint64) (n int) {
	return sovSteeringMemberCandidate(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SteeringMemberCandidateSet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSteeringMemberCandidate
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
			return fmt.Errorf("proto: SteeringMemberCandidateSet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SteeringMemberCandidateSet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SteeringMemberCandidates", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSteeringMemberCandidate
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
				return ErrInvalidLengthSteeringMemberCandidate
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSteeringMemberCandidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SteeringMemberCandidates = append(m.SteeringMemberCandidates, &SteeringMemberCandidate{})
			if err := m.SteeringMemberCandidates[len(m.SteeringMemberCandidates)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSteeringMemberCandidate(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSteeringMemberCandidate
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
func (m *SteeringMemberCandidate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSteeringMemberCandidate
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
			return fmt.Errorf("proto: SteeringMemberCandidate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SteeringMemberCandidate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSteeringMemberCandidate
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
				return ErrInvalidLengthSteeringMemberCandidate
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSteeringMemberCandidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = append(m.Address[:0], dAtA[iNdEx:postIndex]...)
			if m.Address == nil {
				m.Address = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSteeringMemberCandidate
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
				return ErrInvalidLengthSteeringMemberCandidate
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSteeringMemberCandidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PubKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSteeringMemberCandidate(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSteeringMemberCandidate
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
func (m *SimpleSteeringMemberCandidate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSteeringMemberCandidate
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
			return fmt.Errorf("proto: SimpleSteeringMemberCandidate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SimpleSteeringMemberCandidate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSteeringMemberCandidate
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
				return ErrInvalidLengthSteeringMemberCandidate
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSteeringMemberCandidate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PubKey == nil {
				m.PubKey = &crypto.PublicKey{}
			}
			if err := m.PubKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSteeringMemberCandidate(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSteeringMemberCandidate
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
func skipSteeringMemberCandidate(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSteeringMemberCandidate
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
					return 0, ErrIntOverflowSteeringMemberCandidate
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
					return 0, ErrIntOverflowSteeringMemberCandidate
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
				return 0, ErrInvalidLengthSteeringMemberCandidate
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSteeringMemberCandidate
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSteeringMemberCandidate
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSteeringMemberCandidate        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSteeringMemberCandidate          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSteeringMemberCandidate = fmt.Errorf("proto: unexpected end of group")
)
