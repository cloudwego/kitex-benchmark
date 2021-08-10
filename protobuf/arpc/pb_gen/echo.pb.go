// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pb_gen/echo.proto

package arpc

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

type ArpcMsg struct {
	Msg    string `protobuf:"bytes,1,opt,name=Msg,proto3" json:"Msg,omitempty"`
	Finish bool   `protobuf:"varint,2,opt,name=Finish,proto3" json:"Finish,omitempty"`
}

func (m *ArpcMsg) Reset()         { *m = ArpcMsg{} }
func (m *ArpcMsg) String() string { return proto.CompactTextString(m) }
func (*ArpcMsg) ProtoMessage()    {}
func (*ArpcMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf7d8c8c134b619d, []int{0}
}
func (m *ArpcMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ArpcMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ArpcMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ArpcMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ArpcMsg.Merge(m, src)
}
func (m *ArpcMsg) XXX_Size() int {
	return m.Size()
}
func (m *ArpcMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ArpcMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ArpcMsg proto.InternalMessageInfo

func (m *ArpcMsg) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *ArpcMsg) GetFinish() bool {
	if m != nil {
		return m.Finish
	}
	return false
}

func init() {
	proto.RegisterType((*ArpcMsg)(nil), "rpcx.ArpcMsg")
}

func init() { proto.RegisterFile("pb_gen/echo.proto", fileDescriptor_bf7d8c8c134b619d) }

var fileDescriptor_bf7d8c8c134b619d = []byte{
	// 158 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x48, 0x8a, 0x4f,
	0x4f, 0xcd, 0xd3, 0x4f, 0x4d, 0xce, 0xc8, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29,
	0x2a, 0x48, 0xae, 0x50, 0x32, 0xe6, 0x62, 0x0f, 0x2a, 0x48, 0xae, 0xf0, 0x2d, 0x4e, 0x17, 0x12,
	0xe0, 0x62, 0xf6, 0x2d, 0x4e, 0x97, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0x31, 0x85, 0xc4,
	0xb8, 0xd8, 0xdc, 0x32, 0xf3, 0x32, 0x8b, 0x33, 0x24, 0x98, 0x14, 0x18, 0x35, 0x38, 0x82, 0xa0,
	0x3c, 0x23, 0x53, 0x2e, 0x0e, 0x90, 0x26, 0xd7, 0xe4, 0x8c, 0x7c, 0x21, 0x4d, 0x2e, 0x76, 0x90,
	0xa1, 0x20, 0xe5, 0xbc, 0x7a, 0x20, 0x23, 0xf5, 0xa0, 0xe6, 0x49, 0xa1, 0x72, 0x95, 0x18, 0x9c,
	0x24, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f,
	0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x89, 0x0d, 0xec, 0x24, 0x63,
	0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2e, 0xfd, 0xf7, 0xdc, 0xa7, 0x00, 0x00, 0x00,
}

func (m *ArpcMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ArpcMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ArpcMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Finish {
		i--
		if m.Finish {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if len(m.Msg) > 0 {
		i -= len(m.Msg)
		copy(dAtA[i:], m.Msg)
		i = encodeVarintEcho(dAtA, i, uint64(len(m.Msg)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEcho(dAtA []byte, offset int, v uint64) int {
	offset -= sovEcho(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ArpcMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Msg)
	if l > 0 {
		n += 1 + l + sovEcho(uint64(l))
	}
	if m.Finish {
		n += 2
	}
	return n
}

func sovEcho(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEcho(x uint64) (n int) {
	return sovEcho(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ArpcMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEcho
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
			return fmt.Errorf("proto: ArpcMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ArpcMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEcho
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEcho
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEcho
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Finish", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEcho
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
			m.Finish = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipEcho(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEcho
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
func skipEcho(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEcho
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
					return 0, ErrIntOverflowEcho
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
					return 0, ErrIntOverflowEcho
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
				return 0, ErrInvalidLengthEcho
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEcho
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEcho
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEcho        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEcho          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEcho = fmt.Errorf("proto: unexpected end of group")
)
