// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: reapchain/types/qrn.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	crypto "github.com/reapchain/reapchain-core/proto/reapchain/crypto"
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

type QrnSet struct {
	Qrns []*Qrn `protobuf:"bytes,1,rep,name=qrns,proto3" json:"qrns,omitempty"`
}

func (m *QrnSet) Reset()         { *m = QrnSet{} }
func (m *QrnSet) String() string { return proto.CompactTextString(m) }
func (*QrnSet) ProtoMessage()    {}
func (*QrnSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_90e57f19e736c95a, []int{0}
}
func (m *QrnSet) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QrnSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QrnSet.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QrnSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QrnSet.Merge(m, src)
}
func (m *QrnSet) XXX_Size() int {
	return m.Size()
}
func (m *QrnSet) XXX_DiscardUnknown() {
	xxx_messageInfo_QrnSet.DiscardUnknown(m)
}

var xxx_messageInfo_QrnSet proto.InternalMessageInfo

func (m *QrnSet) GetQrns() []*Qrn {
	if m != nil {
		return m.Qrns
	}
	return nil
}

type Qrn struct {
	Height               int64            `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Timestamp            time.Time        `protobuf:"bytes,2,opt,name=timestamp,proto3,stdtime" json:"timestamp"`
	StandingMemberPubKey crypto.PublicKey `protobuf:"bytes,3,opt,name=standing_member_pub_key,json=standingMemberPubKey,proto3" json:"standing_member_pub_key"`
	Value                uint64           `protobuf:"varint,4,opt,name=value,proto3" json:"value,omitempty"`
	Signature            []byte           `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *Qrn) Reset()         { *m = Qrn{} }
func (m *Qrn) String() string { return proto.CompactTextString(m) }
func (*Qrn) ProtoMessage()    {}
func (*Qrn) Descriptor() ([]byte, []int) {
	return fileDescriptor_90e57f19e736c95a, []int{1}
}
func (m *Qrn) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Qrn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Qrn.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Qrn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Qrn.Merge(m, src)
}
func (m *Qrn) XXX_Size() int {
	return m.Size()
}
func (m *Qrn) XXX_DiscardUnknown() {
	xxx_messageInfo_Qrn.DiscardUnknown(m)
}

var xxx_messageInfo_Qrn proto.InternalMessageInfo

func (m *Qrn) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Qrn) GetTimestamp() time.Time {
	if m != nil {
		return m.Timestamp
	}
	return time.Time{}
}

func (m *Qrn) GetStandingMemberPubKey() crypto.PublicKey {
	if m != nil {
		return m.StandingMemberPubKey
	}
	return crypto.PublicKey{}
}

func (m *Qrn) GetValue() uint64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Qrn) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*QrnSet)(nil), "reapchain.types.QrnSet")
	proto.RegisterType((*Qrn)(nil), "reapchain.types.Qrn")
}

func init() { proto.RegisterFile("reapchain/types/qrn.proto", fileDescriptor_90e57f19e736c95a) }

var fileDescriptor_90e57f19e736c95a = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x51, 0xb1, 0x6e, 0xea, 0x30,
	0x14, 0x8d, 0x5f, 0x02, 0x7a, 0x98, 0x27, 0x3d, 0x29, 0x8a, 0xde, 0x4b, 0xa1, 0x0a, 0x11, 0x53,
	0x96, 0xda, 0x12, 0x9d, 0xba, 0x66, 0x45, 0x95, 0x20, 0x74, 0xa8, 0xba, 0xa0, 0x24, 0x75, 0x9d,
	0x08, 0x62, 0x07, 0xc7, 0xa9, 0x94, 0xbf, 0xe0, 0xb3, 0x18, 0x19, 0x3b, 0xb5, 0x15, 0x7c, 0x40,
	0x7f, 0xa1, 0xc2, 0x01, 0x52, 0xb1, 0xdd, 0x73, 0xcf, 0xf1, 0xb9, 0xc7, 0xf7, 0xc2, 0x2b, 0x41,
	0xc2, 0x3c, 0x4e, 0xc2, 0x94, 0x61, 0x59, 0xe5, 0xa4, 0xc0, 0x2b, 0xc1, 0x50, 0x2e, 0xb8, 0xe4,
	0xe6, 0xdf, 0x33, 0x85, 0x14, 0xd5, 0xb3, 0x28, 0xa7, 0x5c, 0x71, 0xf8, 0x50, 0xd5, 0xb2, 0x5e,
	0xbf, 0x71, 0x88, 0x45, 0x95, 0x4b, 0x8e, 0x17, 0xa4, 0x2a, 0x8e, 0xe4, 0x80, 0x72, 0x4e, 0x97,
	0x04, 0x2b, 0x14, 0x95, 0x2f, 0x58, 0xa6, 0x19, 0x29, 0x64, 0x98, 0xe5, 0xb5, 0x60, 0x38, 0x82,
	0xed, 0xa9, 0x60, 0x33, 0x22, 0x4d, 0x0f, 0x1a, 0x2b, 0xc1, 0x0a, 0x1b, 0xb8, 0xba, 0xd7, 0x1d,
	0x59, 0xe8, 0x62, 0x3a, 0x9a, 0x0a, 0x16, 0x28, 0xc5, 0xf0, 0x0b, 0x40, 0x7d, 0x2a, 0x98, 0xf9,
	0x0f, 0xb6, 0x13, 0x92, 0xd2, 0x44, 0xda, 0xc0, 0x05, 0x9e, 0x1e, 0x1c, 0x91, 0xe9, 0xc3, 0xce,
	0x79, 0x8c, 0xfd, 0xcb, 0x05, 0x5e, 0x77, 0xd4, 0x43, 0x75, 0x10, 0x74, 0x0a, 0x82, 0x1e, 0x4e,
	0x0a, 0xff, 0xf7, 0xe6, 0x7d, 0xa0, 0xad, 0x3f, 0x06, 0x20, 0x68, 0x9e, 0x99, 0x8f, 0xf0, 0x7f,
	0x21, 0x43, 0xf6, 0x9c, 0x32, 0x3a, 0xcf, 0x48, 0x16, 0x11, 0x31, 0xcf, 0xcb, 0x68, 0xbe, 0x20,
	0x95, 0xad, 0x2b, 0xc7, 0xfe, 0x8f, 0x80, 0xf5, 0xbf, 0xd1, 0xa4, 0x8c, 0x96, 0x69, 0x3c, 0x26,
	0x95, 0x6f, 0x1c, 0x2c, 0x03, 0xeb, 0xe4, 0x70, 0xaf, 0x0c, 0x26, 0x65, 0x34, 0x26, 0x95, 0x69,
	0xc1, 0xd6, 0x6b, 0xb8, 0x2c, 0x89, 0x6d, 0xb8, 0xc0, 0x33, 0x82, 0x1a, 0x98, 0xd7, 0xb0, 0x53,
	0xa4, 0x94, 0x85, 0xb2, 0x14, 0xc4, 0x6e, 0xb9, 0xc0, 0xfb, 0x13, 0x34, 0x0d, 0x7f, 0xb6, 0xd9,
	0x39, 0x60, 0xbb, 0x73, 0xc0, 0xe7, 0xce, 0x01, 0xeb, 0xbd, 0xa3, 0x6d, 0xf7, 0x8e, 0xf6, 0xb6,
	0x77, 0xb4, 0xa7, 0x3b, 0x9a, 0xca, 0xa4, 0x8c, 0x50, 0xcc, 0x33, 0xdc, 0x1c, 0xe2, 0x5c, 0xdd,
	0xc4, 0x5c, 0x1c, 0xb7, 0x8f, 0x2f, 0x2e, 0x1d, 0xb5, 0x55, 0xfb, 0xf6, 0x3b, 0x00, 0x00, 0xff,
	0xff, 0x3b, 0x44, 0x73, 0xb7, 0x03, 0x02, 0x00, 0x00,
}

func (m *QrnSet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QrnSet) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QrnSet) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Qrns) > 0 {
		for iNdEx := len(m.Qrns) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Qrns[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQrn(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Qrn) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Qrn) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Qrn) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintQrn(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x2a
	}
	if m.Value != 0 {
		i = encodeVarintQrn(dAtA, i, uint64(m.Value))
		i--
		dAtA[i] = 0x20
	}
	{
		size, err := m.StandingMemberPubKey.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQrn(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	n2, err2 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.Timestamp, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.Timestamp):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintQrn(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x12
	if m.Height != 0 {
		i = encodeVarintQrn(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintQrn(dAtA []byte, offset int, v uint64) int {
	offset -= sovQrn(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QrnSet) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Qrns) > 0 {
		for _, e := range m.Qrns {
			l = e.Size()
			n += 1 + l + sovQrn(uint64(l))
		}
	}
	return n
}

func (m *Qrn) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovQrn(uint64(m.Height))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.Timestamp)
	n += 1 + l + sovQrn(uint64(l))
	l = m.StandingMemberPubKey.Size()
	n += 1 + l + sovQrn(uint64(l))
	if m.Value != 0 {
		n += 1 + sovQrn(uint64(m.Value))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovQrn(uint64(l))
	}
	return n
}

func sovQrn(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQrn(x uint64) (n int) {
	return sovQrn(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QrnSet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQrn
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
			return fmt.Errorf("proto: QrnSet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QrnSet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Qrns", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQrn
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
				return ErrInvalidLengthQrn
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQrn
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Qrns = append(m.Qrns, &Qrn{})
			if err := m.Qrns[len(m.Qrns)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQrn(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQrn
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
func (m *Qrn) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQrn
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
			return fmt.Errorf("proto: Qrn: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Qrn: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQrn
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
					return ErrIntOverflowQrn
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
				return ErrInvalidLengthQrn
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQrn
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
				return fmt.Errorf("proto: wrong wireType = %d for field StandingMemberPubKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQrn
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
				return ErrInvalidLengthQrn
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQrn
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StandingMemberPubKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			m.Value = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQrn
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Value |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQrn
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
				return ErrInvalidLengthQrn
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQrn
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQrn(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQrn
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
func skipQrn(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQrn
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
					return 0, ErrIntOverflowQrn
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
					return 0, ErrIntOverflowQrn
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
				return 0, ErrInvalidLengthQrn
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQrn
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQrn
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQrn        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQrn          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQrn = fmt.Errorf("proto: unexpected end of group")
)
