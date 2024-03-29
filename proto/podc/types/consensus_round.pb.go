// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: podc/types/consensus_round.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
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

// ConsensusRound
type ConsensusRound struct {
	ConsensusStartBlockHeight int64  `protobuf:"varint,1,opt,name=consensus_start_block_height,json=consensusStartBlockHeight,proto3" json:"consensus_start_block_height,omitempty"`
	QrnPeriod                 uint64 `protobuf:"varint,2,opt,name=qrn_period,json=qrnPeriod,proto3" json:"qrn_period,omitempty"`
	VrfPeriod                 uint64 `protobuf:"varint,3,opt,name=vrf_period,json=vrfPeriod,proto3" json:"vrf_period,omitempty"`
	ValidatorPeriod           uint64 `protobuf:"varint,4,opt,name=validator_period,json=validatorPeriod,proto3" json:"validator_period,omitempty"`
	Period                    uint64 `protobuf:"varint,5,opt,name=period,proto3" json:"period,omitempty"`
}

func (m *ConsensusRound) Reset()         { *m = ConsensusRound{} }
func (m *ConsensusRound) String() string { return proto.CompactTextString(m) }
func (*ConsensusRound) ProtoMessage()    {}
func (*ConsensusRound) Descriptor() ([]byte, []int) {
	return fileDescriptor_a69a0c5a13b4ff18, []int{0}
}
func (m *ConsensusRound) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ConsensusRound) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ConsensusRound.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ConsensusRound) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsensusRound.Merge(m, src)
}
func (m *ConsensusRound) XXX_Size() int {
	return m.Size()
}
func (m *ConsensusRound) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsensusRound.DiscardUnknown(m)
}

var xxx_messageInfo_ConsensusRound proto.InternalMessageInfo

func (m *ConsensusRound) GetConsensusStartBlockHeight() int64 {
	if m != nil {
		return m.ConsensusStartBlockHeight
	}
	return 0
}

func (m *ConsensusRound) GetQrnPeriod() uint64 {
	if m != nil {
		return m.QrnPeriod
	}
	return 0
}

func (m *ConsensusRound) GetVrfPeriod() uint64 {
	if m != nil {
		return m.VrfPeriod
	}
	return 0
}

func (m *ConsensusRound) GetValidatorPeriod() uint64 {
	if m != nil {
		return m.ValidatorPeriod
	}
	return 0
}

func (m *ConsensusRound) GetPeriod() uint64 {
	if m != nil {
		return m.Period
	}
	return 0
}

func init() {
	proto.RegisterType((*ConsensusRound)(nil), "podc.types.ConsensusRound")
}

func init() { proto.RegisterFile("podc/types/consensus_round.proto", fileDescriptor_a69a0c5a13b4ff18) }

var fileDescriptor_a69a0c5a13b4ff18 = []byte{
	// 273 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0xd0, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x06, 0xe0, 0x98, 0x96, 0x4a, 0x78, 0x00, 0x14, 0x21, 0x54, 0x10, 0x58, 0x11, 0x53, 0x19,
	0x48, 0x06, 0xd8, 0x91, 0xca, 0xc2, 0x88, 0xc2, 0xc6, 0x12, 0x25, 0x8e, 0x9b, 0x44, 0x14, 0x5f,
	0x7a, 0x76, 0x22, 0xf1, 0x16, 0x3c, 0x06, 0x8f, 0xc2, 0x58, 0x31, 0x31, 0xa2, 0xe4, 0x45, 0x2a,
	0x3b, 0x49, 0xb3, 0x9d, 0xff, 0xff, 0x3b, 0x0f, 0x47, 0xbd, 0x12, 0x52, 0x1e, 0xe8, 0xcf, 0x52,
	0xa8, 0x80, 0x83, 0x54, 0x42, 0xaa, 0x4a, 0x45, 0x08, 0x95, 0x4c, 0xfd, 0x12, 0x41, 0x83, 0x4b,
	0x8d, 0xf0, 0xad, 0xb8, 0x3c, 0xcb, 0x20, 0x03, 0x1b, 0x07, 0x66, 0xea, 0xc4, 0xcd, 0x2f, 0xa1,
	0xc7, 0x4f, 0xc3, 0x6e, 0x68, 0x56, 0xdd, 0x47, 0x7a, 0x35, 0xfe, 0xa6, 0x74, 0x8c, 0x3a, 0x4a,
	0xd6, 0xc0, 0xdf, 0xa3, 0x5c, 0x14, 0x59, 0xae, 0xe7, 0xc4, 0x23, 0x8b, 0x49, 0x78, 0xb1, 0x37,
	0xaf, 0x86, 0x2c, 0x8d, 0x78, 0xb6, 0xc0, 0xbd, 0xa6, 0x74, 0x83, 0x32, 0x2a, 0x05, 0x16, 0x90,
	0xce, 0x0f, 0x3c, 0xb2, 0x98, 0x86, 0x47, 0x1b, 0x94, 0x2f, 0x36, 0x30, 0x75, 0x8d, 0xab, 0xa1,
	0x9e, 0x74, 0x75, 0x8d, 0xab, 0xbe, 0xbe, 0xa5, 0xa7, 0x75, 0xbc, 0x2e, 0xd2, 0x58, 0x03, 0x0e,
	0x68, 0x6a, 0xd1, 0xc9, 0x3e, 0xef, 0xe9, 0x39, 0x9d, 0xf5, 0xe0, 0xd0, 0x82, 0xfe, 0xb5, 0x0c,
	0xbf, 0x1b, 0x46, 0x7e, 0x1a, 0x46, 0xb6, 0x0d, 0x23, 0xff, 0x0d, 0x23, 0x5f, 0x2d, 0x73, 0xb6,
	0x2d, 0x73, 0xfe, 0x5a, 0xe6, 0xbc, 0x3d, 0x64, 0x85, 0xce, 0xab, 0xc4, 0xe7, 0xf0, 0x11, 0xa0,
	0x88, 0x4b, 0x9e, 0xc7, 0x85, 0x1c, 0xa7, 0x3b, 0x0e, 0x28, 0x82, 0xee, 0x50, 0xe3, 0x81, 0x93,
	0x99, 0x4d, 0xee, 0x77, 0x01, 0x00, 0x00, 0xff, 0xff, 0x24, 0xd1, 0xd2, 0x4f, 0x75, 0x01, 0x00,
	0x00,
}

func (this *ConsensusRound) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ConsensusRound)
	if !ok {
		that2, ok := that.(ConsensusRound)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.ConsensusStartBlockHeight != that1.ConsensusStartBlockHeight {
		return false
	}
	if this.QrnPeriod != that1.QrnPeriod {
		return false
	}
	if this.VrfPeriod != that1.VrfPeriod {
		return false
	}
	if this.ValidatorPeriod != that1.ValidatorPeriod {
		return false
	}
	if this.Period != that1.Period {
		return false
	}
	return true
}
func (m *ConsensusRound) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ConsensusRound) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ConsensusRound) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Period != 0 {
		i = encodeVarintConsensusRound(dAtA, i, uint64(m.Period))
		i--
		dAtA[i] = 0x28
	}
	if m.ValidatorPeriod != 0 {
		i = encodeVarintConsensusRound(dAtA, i, uint64(m.ValidatorPeriod))
		i--
		dAtA[i] = 0x20
	}
	if m.VrfPeriod != 0 {
		i = encodeVarintConsensusRound(dAtA, i, uint64(m.VrfPeriod))
		i--
		dAtA[i] = 0x18
	}
	if m.QrnPeriod != 0 {
		i = encodeVarintConsensusRound(dAtA, i, uint64(m.QrnPeriod))
		i--
		dAtA[i] = 0x10
	}
	if m.ConsensusStartBlockHeight != 0 {
		i = encodeVarintConsensusRound(dAtA, i, uint64(m.ConsensusStartBlockHeight))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintConsensusRound(dAtA []byte, offset int, v uint64) int {
	offset -= sovConsensusRound(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ConsensusRound) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ConsensusStartBlockHeight != 0 {
		n += 1 + sovConsensusRound(uint64(m.ConsensusStartBlockHeight))
	}
	if m.QrnPeriod != 0 {
		n += 1 + sovConsensusRound(uint64(m.QrnPeriod))
	}
	if m.VrfPeriod != 0 {
		n += 1 + sovConsensusRound(uint64(m.VrfPeriod))
	}
	if m.ValidatorPeriod != 0 {
		n += 1 + sovConsensusRound(uint64(m.ValidatorPeriod))
	}
	if m.Period != 0 {
		n += 1 + sovConsensusRound(uint64(m.Period))
	}
	return n
}

func sovConsensusRound(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozConsensusRound(x uint64) (n int) {
	return sovConsensusRound(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ConsensusRound) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowConsensusRound
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
			return fmt.Errorf("proto: ConsensusRound: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ConsensusRound: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConsensusStartBlockHeight", wireType)
			}
			m.ConsensusStartBlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConsensusRound
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ConsensusStartBlockHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field QrnPeriod", wireType)
			}
			m.QrnPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConsensusRound
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.QrnPeriod |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VrfPeriod", wireType)
			}
			m.VrfPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConsensusRound
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VrfPeriod |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorPeriod", wireType)
			}
			m.ValidatorPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConsensusRound
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ValidatorPeriod |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Period", wireType)
			}
			m.Period = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConsensusRound
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Period |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipConsensusRound(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthConsensusRound
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
func skipConsensusRound(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowConsensusRound
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
					return 0, ErrIntOverflowConsensusRound
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
					return 0, ErrIntOverflowConsensusRound
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
				return 0, ErrInvalidLengthConsensusRound
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupConsensusRound
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthConsensusRound
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthConsensusRound        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowConsensusRound          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupConsensusRound = fmt.Errorf("proto: unexpected end of group")
)
