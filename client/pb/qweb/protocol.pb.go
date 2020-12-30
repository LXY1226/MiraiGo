// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: protocol.proto

package qweb

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

type QWebReq struct {
	Seq        int64  `protobuf:"varint,1,opt,name=seq,proto3" json:"seq,omitempty"`
	Qua        string `protobuf:"bytes,2,opt,name=qua,proto3" json:"qua,omitempty"`
	DeviceInfo string `protobuf:"bytes,3,opt,name=deviceInfo,proto3" json:"deviceInfo,omitempty"`
	BusiBuff   []byte `protobuf:"bytes,4,opt,name=busiBuff,proto3" json:"busiBuff,omitempty"`
	TraceId    string `protobuf:"bytes,5,opt,name=traceId,proto3" json:"traceId,omitempty"`
}

func (m *QWebReq) Reset()         { *m = QWebReq{} }
func (m *QWebReq) String() string { return proto.CompactTextString(m) }
func (*QWebReq) ProtoMessage()    {}
func (*QWebReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{0}
}
func (m *QWebReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QWebReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QWebReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QWebReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QWebReq.Merge(m, src)
}
func (m *QWebReq) XXX_Size() int {
	return m.Size()
}
func (m *QWebReq) XXX_DiscardUnknown() {
	xxx_messageInfo_QWebReq.DiscardUnknown(m)
}

var xxx_messageInfo_QWebReq proto.InternalMessageInfo

func (m *QWebReq) GetSeq() int64 {
	if m != nil {
		return m.Seq
	}
	return 0
}

func (m *QWebReq) GetQua() string {
	if m != nil {
		return m.Qua
	}
	return ""
}

func (m *QWebReq) GetDeviceInfo() string {
	if m != nil {
		return m.DeviceInfo
	}
	return ""
}

func (m *QWebReq) GetBusiBuff() []byte {
	if m != nil {
		return m.BusiBuff
	}
	return nil
}

func (m *QWebReq) GetTraceId() string {
	if m != nil {
		return m.TraceId
	}
	return ""
}

type QWebRsp struct {
	Seq      int64  `protobuf:"varint,1,opt,name=seq,proto3" json:"seq,omitempty"`
	RetCode  int64  `protobuf:"varint,2,opt,name=retCode,proto3" json:"retCode,omitempty"`
	ErrMsg   string `protobuf:"bytes,3,opt,name=errMsg,proto3" json:"errMsg,omitempty"`
	BusiBuff []byte `protobuf:"bytes,4,opt,name=busiBuff,proto3" json:"busiBuff,omitempty"`
}

func (m *QWebRsp) Reset()         { *m = QWebRsp{} }
func (m *QWebRsp) String() string { return proto.CompactTextString(m) }
func (*QWebRsp) ProtoMessage()    {}
func (*QWebRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc2336598a3f7e0, []int{1}
}
func (m *QWebRsp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QWebRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QWebRsp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QWebRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QWebRsp.Merge(m, src)
}
func (m *QWebRsp) XXX_Size() int {
	return m.Size()
}
func (m *QWebRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_QWebRsp.DiscardUnknown(m)
}

var xxx_messageInfo_QWebRsp proto.InternalMessageInfo

func (m *QWebRsp) GetSeq() int64 {
	if m != nil {
		return m.Seq
	}
	return 0
}

func (m *QWebRsp) GetRetCode() int64 {
	if m != nil {
		return m.RetCode
	}
	return 0
}

func (m *QWebRsp) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

func (m *QWebRsp) GetBusiBuff() []byte {
	if m != nil {
		return m.BusiBuff
	}
	return nil
}

func init() {
	proto.RegisterType((*QWebReq)(nil), "QWebReq")
	proto.RegisterType((*QWebRsp)(nil), "QWebRsp")
}

func init() { proto.RegisterFile("protocol.proto", fileDescriptor_2bc2336598a3f7e0) }

var fileDescriptor_2bc2336598a3f7e0 = []byte{
	// 218 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x4f, 0xce, 0xcf, 0xd1, 0x03, 0x33, 0x94, 0x9a, 0x19, 0xb9, 0xd8, 0x03, 0xc3, 0x53, 0x93,
	0x82, 0x52, 0x0b, 0x85, 0x04, 0xb8, 0x98, 0x8b, 0x53, 0x0b, 0x25, 0x18, 0x15, 0x18, 0x35, 0x98,
	0x83, 0x40, 0x4c, 0x90, 0x48, 0x61, 0x69, 0xa2, 0x04, 0x93, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x88,
	0x29, 0x24, 0xc7, 0xc5, 0x95, 0x92, 0x5a, 0x96, 0x99, 0x9c, 0xea, 0x99, 0x97, 0x96, 0x2f, 0xc1,
	0x0c, 0x96, 0x40, 0x12, 0x11, 0x92, 0xe2, 0xe2, 0x48, 0x2a, 0x2d, 0xce, 0x74, 0x2a, 0x4d, 0x4b,
	0x93, 0x60, 0x51, 0x60, 0xd4, 0xe0, 0x09, 0x82, 0xf3, 0x85, 0x24, 0xb8, 0xd8, 0x4b, 0x8a, 0x12,
	0x93, 0x53, 0x3d, 0x53, 0x24, 0x58, 0xc1, 0x1a, 0x61, 0x5c, 0xa5, 0x4c, 0xa8, 0x23, 0x8a, 0x0b,
	0xb0, 0x38, 0x42, 0x82, 0x8b, 0xbd, 0x28, 0xb5, 0xc4, 0x39, 0x3f, 0x25, 0x15, 0xec, 0x10, 0xe6,
	0x20, 0x18, 0x57, 0x48, 0x8c, 0x8b, 0x2d, 0xb5, 0xa8, 0xc8, 0xb7, 0x38, 0x1d, 0xea, 0x10, 0x28,
	0x0f, 0x9f, 0x23, 0x9c, 0x14, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23,
	0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x8a,
	0x4d, 0xcf, 0xba, 0xb0, 0x3c, 0x35, 0x29, 0x89, 0x0d, 0x1c, 0x32, 0xc6, 0x80, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xb6, 0x10, 0xed, 0xa1, 0x2b, 0x01, 0x00, 0x00,
}

func (m *QWebReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QWebReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QWebReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TraceId) > 0 {
		i -= len(m.TraceId)
		copy(dAtA[i:], m.TraceId)
		i = encodeVarintProtocol(dAtA, i, uint64(len(m.TraceId)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.BusiBuff) > 0 {
		i -= len(m.BusiBuff)
		copy(dAtA[i:], m.BusiBuff)
		i = encodeVarintProtocol(dAtA, i, uint64(len(m.BusiBuff)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.DeviceInfo) > 0 {
		i -= len(m.DeviceInfo)
		copy(dAtA[i:], m.DeviceInfo)
		i = encodeVarintProtocol(dAtA, i, uint64(len(m.DeviceInfo)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Qua) > 0 {
		i -= len(m.Qua)
		copy(dAtA[i:], m.Qua)
		i = encodeVarintProtocol(dAtA, i, uint64(len(m.Qua)))
		i--
		dAtA[i] = 0x12
	}
	if m.Seq != 0 {
		i = encodeVarintProtocol(dAtA, i, uint64(m.Seq))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QWebRsp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QWebRsp) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QWebRsp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.BusiBuff) > 0 {
		i -= len(m.BusiBuff)
		copy(dAtA[i:], m.BusiBuff)
		i = encodeVarintProtocol(dAtA, i, uint64(len(m.BusiBuff)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ErrMsg) > 0 {
		i -= len(m.ErrMsg)
		copy(dAtA[i:], m.ErrMsg)
		i = encodeVarintProtocol(dAtA, i, uint64(len(m.ErrMsg)))
		i--
		dAtA[i] = 0x1a
	}
	if m.RetCode != 0 {
		i = encodeVarintProtocol(dAtA, i, uint64(m.RetCode))
		i--
		dAtA[i] = 0x10
	}
	if m.Seq != 0 {
		i = encodeVarintProtocol(dAtA, i, uint64(m.Seq))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintProtocol(dAtA []byte, offset int, v uint64) int {
	offset -= sovProtocol(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QWebReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Seq != 0 {
		n += 1 + sovProtocol(uint64(m.Seq))
	}
	l = len(m.Qua)
	if l > 0 {
		n += 1 + l + sovProtocol(uint64(l))
	}
	l = len(m.DeviceInfo)
	if l > 0 {
		n += 1 + l + sovProtocol(uint64(l))
	}
	l = len(m.BusiBuff)
	if l > 0 {
		n += 1 + l + sovProtocol(uint64(l))
	}
	l = len(m.TraceId)
	if l > 0 {
		n += 1 + l + sovProtocol(uint64(l))
	}
	return n
}

func (m *QWebRsp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Seq != 0 {
		n += 1 + sovProtocol(uint64(m.Seq))
	}
	if m.RetCode != 0 {
		n += 1 + sovProtocol(uint64(m.RetCode))
	}
	l = len(m.ErrMsg)
	if l > 0 {
		n += 1 + l + sovProtocol(uint64(l))
	}
	l = len(m.BusiBuff)
	if l > 0 {
		n += 1 + l + sovProtocol(uint64(l))
	}
	return n
}

func sovProtocol(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProtocol(x uint64) (n int) {
	return sovProtocol(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QWebReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProtocol
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
			return fmt.Errorf("proto: QWebReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QWebReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Seq", wireType)
			}
			m.Seq = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtocol
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Seq |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Qua", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtocol
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
				return ErrInvalidLengthProtocol
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProtocol
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Qua = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeviceInfo", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtocol
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
				return ErrInvalidLengthProtocol
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProtocol
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DeviceInfo = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BusiBuff", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtocol
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
				return ErrInvalidLengthProtocol
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProtocol
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BusiBuff = append(m.BusiBuff[:0], dAtA[iNdEx:postIndex]...)
			if m.BusiBuff == nil {
				m.BusiBuff = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TraceId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtocol
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
				return ErrInvalidLengthProtocol
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProtocol
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TraceId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProtocol(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthProtocol
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthProtocol
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
func (m *QWebRsp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProtocol
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
			return fmt.Errorf("proto: QWebRsp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QWebRsp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Seq", wireType)
			}
			m.Seq = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtocol
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Seq |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RetCode", wireType)
			}
			m.RetCode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtocol
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RetCode |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrMsg", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtocol
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
				return ErrInvalidLengthProtocol
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProtocol
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ErrMsg = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BusiBuff", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtocol
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
				return ErrInvalidLengthProtocol
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProtocol
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BusiBuff = append(m.BusiBuff[:0], dAtA[iNdEx:postIndex]...)
			if m.BusiBuff == nil {
				m.BusiBuff = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProtocol(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthProtocol
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthProtocol
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
func skipProtocol(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProtocol
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
					return 0, ErrIntOverflowProtocol
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
					return 0, ErrIntOverflowProtocol
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
				return 0, ErrInvalidLengthProtocol
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProtocol
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProtocol
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProtocol        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProtocol          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProtocol = fmt.Errorf("proto: unexpected end of group")
)
