// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: message.proto

/*
	Package mpb is a generated protocol buffer package.

	It is generated from these files:
		message.proto

	It has these top-level messages:
		Message
*/
package mpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Identify int32

const (
	Identify_MSG_ZERO   Identify = 0
	Identify_MSG_STRING Identify = 1
	Identify_MSG_NUMBER Identify = 2
)

var Identify_name = map[int32]string{
	0: "MSG_ZERO",
	1: "MSG_STRING",
	2: "MSG_NUMBER",
}
var Identify_value = map[string]int32{
	"MSG_ZERO":   0,
	"MSG_STRING": 1,
	"MSG_NUMBER": 2,
}

func (x Identify) String() string {
	return proto.EnumName(Identify_name, int32(x))
}
func (Identify) EnumDescriptor() ([]byte, []int) { return fileDescriptorMessage, []int{0} }

type Message struct {
	Identify Identify `protobuf:"varint,1,opt,name=Identify,proto3,enum=mpb.Identify" json:"Identify,omitempty"`
	Payload  []byte   `protobuf:"bytes,2,opt,name=Payload,proto3" json:"Payload,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptorMessage, []int{0} }

func (m *Message) GetIdentify() Identify {
	if m != nil {
		return m.Identify
	}
	return Identify_MSG_ZERO
}

func (m *Message) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "mpb.Message")
	proto.RegisterEnum("mpb.Identify", Identify_name, Identify_value)
}
func (m *Message) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Identify != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintMessage(dAtA, i, uint64(m.Identify))
	}
	if len(m.Payload) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Payload)))
		i += copy(dAtA[i:], m.Payload)
	}
	return i, nil
}

func encodeVarintMessage(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Message) Size() (n int) {
	var l int
	_ = l
	if m.Identify != 0 {
		n += 1 + sovMessage(uint64(m.Identify))
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	return n
}

func sovMessage(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozMessage(x uint64) (n int) {
	return sovMessage(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Message) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Identify", wireType)
			}
			m.Identify = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Identify |= (Identify(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessage
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
func skipMessage(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMessage
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
					return 0, ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessage
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthMessage
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowMessage
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipMessage(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthMessage = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMessage   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("message.proto", fileDescriptorMessage) }

var fileDescriptorMessage = []byte{
	// 166 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xce, 0x2d, 0x48, 0x52, 0xf2,
	0xe3, 0x62, 0xf7, 0x85, 0x88, 0x0a, 0x69, 0x72, 0x71, 0x78, 0xa6, 0xa4, 0xe6, 0x95, 0x64, 0xa6,
	0x55, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x19, 0xf1, 0xea, 0xe5, 0x16, 0x24, 0xe9, 0xc1, 0x04,
	0x83, 0xe0, 0xd2, 0x42, 0x12, 0x5c, 0xec, 0x01, 0x89, 0x95, 0x39, 0xf9, 0x89, 0x29, 0x12, 0x4c,
	0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x30, 0xae, 0x96, 0x05, 0xc2, 0x10, 0x21, 0x1e, 0x2e, 0x0e, 0xdf,
	0x60, 0xf7, 0xf8, 0x28, 0xd7, 0x20, 0x7f, 0x01, 0x06, 0x21, 0x3e, 0x2e, 0x2e, 0x10, 0x2f, 0x38,
	0x24, 0xc8, 0xd3, 0xcf, 0x5d, 0x80, 0x11, 0xc6, 0xf7, 0x0b, 0xf5, 0x75, 0x72, 0x0d, 0x12, 0x60,
	0x72, 0x12, 0x38, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x67,
	0x3c, 0x96, 0x63, 0x48, 0x62, 0x03, 0xbb, 0xd3, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xfb, 0x2a,
	0xd3, 0x4f, 0xb8, 0x00, 0x00, 0x00,
}