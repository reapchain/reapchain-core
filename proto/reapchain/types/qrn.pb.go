// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: reapchain/types/qrn.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	crypto "github.com/reapchain/reapchain-core/proto/reapchain/crypto"
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
	Address []byte           `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	PubKey  crypto.PublicKey `protobuf:"bytes,2,opt,name=pub_key,json=pubKey,proto3" json:"pub_key"`
	Value   uint64           `protobuf:"varint,3,opt,name=value,proto3" json:"value,omitempty"`
	Height  int64            `protobuf:"varint,4,opt,name=height,proto3" json:"height,omitempty"`
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

func (m *Qrn) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *Qrn) GetPubKey() crypto.PublicKey {
	if m != nil {
		return m.PubKey
	}
	return crypto.PublicKey{}
}

func (m *Qrn) GetValue() uint64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Qrn) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type SimpleQrn struct {
	PubKey *crypto.PublicKey `protobuf:"bytes,1,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
	Value  uint64            `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	Height int64             `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *SimpleQrn) Reset()         { *m = SimpleQrn{} }
func (m *SimpleQrn) String() string { return proto.CompactTextString(m) }
func (*SimpleQrn) ProtoMessage()    {}
func (*SimpleQrn) Descriptor() ([]byte, []int) {
	return fileDescriptor_90e57f19e736c95a, []int{2}
}
func (m *SimpleQrn) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SimpleQrn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SimpleQrn.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SimpleQrn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SimpleQrn.Merge(m, src)
}
func (m *SimpleQrn) XXX_Size() int {
	return m.Size()
}
func (m *SimpleQrn) XXX_DiscardUnknown() {
	xxx_messageInfo_SimpleQrn.DiscardUnknown(m)
}

var xxx_messageInfo_SimpleQrn proto.InternalMessageInfo

func (m *SimpleQrn) GetPubKey() *crypto.PublicKey {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *SimpleQrn) GetValue() uint64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *SimpleQrn) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func init() {
	proto.RegisterType((*QrnSet)(nil), "reapchain.types.QrnSet")
	proto.RegisterType((*Qrn)(nil), "reapchain.types.Qrn")
	proto.RegisterType((*SimpleQrn)(nil), "reapchain.types.SimpleQrn")
}

func init() { proto.RegisterFile("reapchain/types/qrn.proto", fileDescriptor_90e57f19e736c95a) }

var fileDescriptor_90e57f19e736c95a = []byte{
	// 320 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xc1, 0x4e, 0xfa, 0x40,
	0x10, 0xc6, 0xbb, 0xb4, 0xff, 0x92, 0xff, 0x62, 0x62, 0xd2, 0x10, 0x53, 0x21, 0xa9, 0x0d, 0xa7,
	0x5e, 0xdc, 0x4d, 0xd0, 0x8b, 0x1e, 0xb9, 0x72, 0x91, 0x72, 0xf3, 0x62, 0xda, 0x32, 0x69, 0x1b,
	0x60, 0x77, 0xd9, 0x6e, 0x4d, 0xf6, 0x11, 0xbc, 0xf9, 0x58, 0x1c, 0x39, 0x7a, 0x32, 0x06, 0x5e,
	0xc4, 0xd0, 0x2a, 0x25, 0x68, 0xe2, 0x6d, 0xbe, 0xfd, 0x66, 0xbf, 0xdf, 0x64, 0x06, 0x5f, 0x4a,
	0x88, 0x44, 0x92, 0x45, 0x39, 0xa3, 0x4a, 0x0b, 0x28, 0xe8, 0x4a, 0x32, 0x22, 0x24, 0x57, 0xdc,
	0x39, 0x3f, 0x58, 0xa4, 0xb2, 0x7a, 0xdd, 0x94, 0xa7, 0xbc, 0xf2, 0xe8, 0xbe, 0xaa, 0xdb, 0x7a,
	0xfd, 0x26, 0x21, 0x91, 0x5a, 0x28, 0x4e, 0xe7, 0xa0, 0x8b, 0xda, 0x1c, 0x0c, 0xb1, 0x3d, 0x91,
	0x6c, 0x0a, 0xca, 0x09, 0xb0, 0xb5, 0x92, 0xac, 0x70, 0x91, 0x6f, 0x06, 0x9d, 0x61, 0x97, 0x9c,
	0x84, 0x93, 0x89, 0x64, 0x61, 0xd5, 0x31, 0x78, 0x41, 0xd8, 0x9c, 0x48, 0xe6, 0xb8, 0xb8, 0x1d,
	0xcd, 0x66, 0x12, 0x8a, 0xfd, 0x27, 0x14, 0x9c, 0x85, 0xdf, 0xd2, 0xb9, 0xc7, 0x6d, 0x51, 0xc6,
	0x4f, 0x73, 0xd0, 0x6e, 0xcb, 0x47, 0x41, 0x67, 0xd8, 0x3f, 0x8a, 0xab, 0x87, 0x20, 0x0f, 0x65,
	0xbc, 0xc8, 0x93, 0x31, 0xe8, 0x91, 0xb5, 0x7e, 0xbf, 0x32, 0x42, 0x5b, 0x94, 0xf1, 0x18, 0xb4,
	0xd3, 0xc5, 0xff, 0x9e, 0xa3, 0x45, 0x09, 0xae, 0xe9, 0xa3, 0xc0, 0x0a, 0x6b, 0xe1, 0x5c, 0x60,
	0x3b, 0x83, 0x3c, 0xcd, 0x94, 0x6b, 0xf9, 0x28, 0x30, 0xc3, 0x2f, 0x35, 0xe0, 0xf8, 0xff, 0x34,
	0x5f, 0x8a, 0x05, 0xec, 0x07, 0xba, 0x6d, 0xb0, 0xe8, 0x4f, 0xec, 0x4f, 0x60, 0xeb, 0x77, 0xa0,
	0x79, 0x0c, 0x1c, 0x4d, 0xd7, 0x5b, 0x0f, 0x6d, 0xb6, 0x1e, 0xfa, 0xd8, 0x7a, 0xe8, 0x75, 0xe7,
	0x19, 0x9b, 0x9d, 0x67, 0xbc, 0xed, 0x3c, 0xe3, 0xf1, 0x2e, 0xcd, 0x55, 0x56, 0xc6, 0x24, 0xe1,
	0x4b, 0xda, 0xac, 0xfc, 0x50, 0x5d, 0x27, 0x5c, 0x02, 0xad, 0xaf, 0x73, 0x72, 0xd3, 0xd8, 0xae,
	0x9e, 0x6f, 0x3e, 0x03, 0x00, 0x00, 0xff, 0xff, 0x59, 0x97, 0x19, 0xf3, 0xed, 0x01, 0x00, 0x00,
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
	if m.Height != 0 {
		i = encodeVarintQrn(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x20
	}
	if m.Value != 0 {
		i = encodeVarintQrn(dAtA, i, uint64(m.Value))
		i--
		dAtA[i] = 0x18
	}
	{
		size, err := m.PubKey.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQrn(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintQrn(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SimpleQrn) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SimpleQrn) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SimpleQrn) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Height != 0 {
		i = encodeVarintQrn(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x18
	}
	if m.Value != 0 {
		i = encodeVarintQrn(dAtA, i, uint64(m.Value))
		i--
		dAtA[i] = 0x10
	}
	if m.PubKey != nil {
		{
			size, err := m.PubKey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQrn(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
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
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovQrn(uint64(l))
	}
	l = m.PubKey.Size()
	n += 1 + l + sovQrn(uint64(l))
	if m.Value != 0 {
		n += 1 + sovQrn(uint64(m.Value))
	}
	if m.Height != 0 {
		n += 1 + sovQrn(uint64(m.Height))
	}
	return n
}

func (m *SimpleQrn) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PubKey != nil {
		l = m.PubKey.Size()
		n += 1 + l + sovQrn(uint64(l))
	}
	if m.Value != 0 {
		n += 1 + sovQrn(uint64(m.Value))
	}
	if m.Height != 0 {
		n += 1 + sovQrn(uint64(m.Height))
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			if err := m.PubKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
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
		case 4:
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
func (m *SimpleQrn) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: SimpleQrn: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SimpleQrn: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
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
			if m.PubKey == nil {
				m.PubKey = &crypto.PublicKey{}
			}
			if err := m.PubKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
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
		case 3:
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
