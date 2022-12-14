// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.9.1
// source: hashlock.proto

package types

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Hashlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HashlockId    []byte `protobuf:"bytes,1,opt,name=hashlockId,proto3" json:"hashlockId,omitempty"`
	Status        int32  `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
	CreateTime    int64  `protobuf:"varint,3,opt,name=CreateTime,proto3" json:"CreateTime,omitempty"`
	ToAddress     string `protobuf:"bytes,4,opt,name=toAddress,proto3" json:"toAddress,omitempty"`
	ReturnAddress string `protobuf:"bytes,5,opt,name=returnAddress,proto3" json:"returnAddress,omitempty"`
	Amount        int64  `protobuf:"varint,6,opt,name=amount,proto3" json:"amount,omitempty"`
	Frozentime    int64  `protobuf:"varint,7,opt,name=frozentime,proto3" json:"frozentime,omitempty"`
}

func (x *Hashlock) Reset() {
	*x = Hashlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hashlock_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Hashlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Hashlock) ProtoMessage() {}

func (x *Hashlock) ProtoReflect() protoreflect.Message {
	mi := &file_hashlock_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Hashlock.ProtoReflect.Descriptor instead.
func (*Hashlock) Descriptor() ([]byte, []int) {
	return file_hashlock_proto_rawDescGZIP(), []int{0}
}

func (x *Hashlock) GetHashlockId() []byte {
	if x != nil {
		return x.HashlockId
	}
	return nil
}

func (x *Hashlock) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Hashlock) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

func (x *Hashlock) GetToAddress() string {
	if x != nil {
		return x.ToAddress
	}
	return ""
}

func (x *Hashlock) GetReturnAddress() string {
	if x != nil {
		return x.ReturnAddress
	}
	return ""
}

func (x *Hashlock) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *Hashlock) GetFrozentime() int64 {
	if x != nil {
		return x.Frozentime
	}
	return 0
}

type HashlockLock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount        int64  `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
	Time          int64  `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	Hash          []byte `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	ToAddress     string `protobuf:"bytes,4,opt,name=toAddress,proto3" json:"toAddress,omitempty"`
	ReturnAddress string `protobuf:"bytes,5,opt,name=returnAddress,proto3" json:"returnAddress,omitempty"`
}

func (x *HashlockLock) Reset() {
	*x = HashlockLock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hashlock_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HashlockLock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HashlockLock) ProtoMessage() {}

func (x *HashlockLock) ProtoReflect() protoreflect.Message {
	mi := &file_hashlock_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HashlockLock.ProtoReflect.Descriptor instead.
func (*HashlockLock) Descriptor() ([]byte, []int) {
	return file_hashlock_proto_rawDescGZIP(), []int{1}
}

func (x *HashlockLock) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *HashlockLock) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *HashlockLock) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *HashlockLock) GetToAddress() string {
	if x != nil {
		return x.ToAddress
	}
	return ""
}

func (x *HashlockLock) GetReturnAddress() string {
	if x != nil {
		return x.ReturnAddress
	}
	return ""
}

type HashlockSend struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Secret []byte `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"` // bytes  hash     = 3;
}

func (x *HashlockSend) Reset() {
	*x = HashlockSend{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hashlock_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HashlockSend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HashlockSend) ProtoMessage() {}

func (x *HashlockSend) ProtoReflect() protoreflect.Message {
	mi := &file_hashlock_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HashlockSend.ProtoReflect.Descriptor instead.
func (*HashlockSend) Descriptor() ([]byte, []int) {
	return file_hashlock_proto_rawDescGZIP(), []int{2}
}

func (x *HashlockSend) GetSecret() []byte {
	if x != nil {
		return x.Secret
	}
	return nil
}

type Hashlockquery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time        int64 `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	Status      int32 `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
	Amount      int64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	CreateTime  int64 `protobuf:"varint,4,opt,name=createTime,proto3" json:"createTime,omitempty"`
	CurrentTime int64 `protobuf:"varint,5,opt,name=currentTime,proto3" json:"currentTime,omitempty"`
}

func (x *Hashlockquery) Reset() {
	*x = Hashlockquery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hashlock_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Hashlockquery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Hashlockquery) ProtoMessage() {}

func (x *Hashlockquery) ProtoReflect() protoreflect.Message {
	mi := &file_hashlock_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Hashlockquery.ProtoReflect.Descriptor instead.
func (*Hashlockquery) Descriptor() ([]byte, []int) {
	return file_hashlock_proto_rawDescGZIP(), []int{3}
}

func (x *Hashlockquery) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Hashlockquery) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Hashlockquery) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *Hashlockquery) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

func (x *Hashlockquery) GetCurrentTime() int64 {
	if x != nil {
		return x.CurrentTime
	}
	return 0
}

type HashRecv struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HashlockId  []byte         `protobuf:"bytes,1,opt,name=HashlockId,proto3" json:"HashlockId,omitempty"`
	Information *Hashlockquery `protobuf:"bytes,2,opt,name=Information,proto3" json:"Information,omitempty"`
}

func (x *HashRecv) Reset() {
	*x = HashRecv{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hashlock_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HashRecv) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HashRecv) ProtoMessage() {}

func (x *HashRecv) ProtoReflect() protoreflect.Message {
	mi := &file_hashlock_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HashRecv.ProtoReflect.Descriptor instead.
func (*HashRecv) Descriptor() ([]byte, []int) {
	return file_hashlock_proto_rawDescGZIP(), []int{4}
}

func (x *HashRecv) GetHashlockId() []byte {
	if x != nil {
		return x.HashlockId
	}
	return nil
}

func (x *HashRecv) GetInformation() *Hashlockquery {
	if x != nil {
		return x.Information
	}
	return nil
}

type HashlockUnlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Secret []byte `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"` // bytes  hash     = 3;
}

func (x *HashlockUnlock) Reset() {
	*x = HashlockUnlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hashlock_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HashlockUnlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HashlockUnlock) ProtoMessage() {}

func (x *HashlockUnlock) ProtoReflect() protoreflect.Message {
	mi := &file_hashlock_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HashlockUnlock.ProtoReflect.Descriptor instead.
func (*HashlockUnlock) Descriptor() ([]byte, []int) {
	return file_hashlock_proto_rawDescGZIP(), []int{5}
}

func (x *HashlockUnlock) GetSecret() []byte {
	if x != nil {
		return x.Secret
	}
	return nil
}

// message for hashlock
type HashlockAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//	*HashlockAction_Hlock
	//	*HashlockAction_Hsend
	//	*HashlockAction_Hunlock
	Value isHashlockAction_Value `protobuf_oneof:"value"`
	Ty    int32                  `protobuf:"varint,4,opt,name=ty,proto3" json:"ty,omitempty"`
}

func (x *HashlockAction) Reset() {
	*x = HashlockAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hashlock_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HashlockAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HashlockAction) ProtoMessage() {}

func (x *HashlockAction) ProtoReflect() protoreflect.Message {
	mi := &file_hashlock_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HashlockAction.ProtoReflect.Descriptor instead.
func (*HashlockAction) Descriptor() ([]byte, []int) {
	return file_hashlock_proto_rawDescGZIP(), []int{6}
}

func (m *HashlockAction) GetValue() isHashlockAction_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *HashlockAction) GetHlock() *HashlockLock {
	if x, ok := x.GetValue().(*HashlockAction_Hlock); ok {
		return x.Hlock
	}
	return nil
}

func (x *HashlockAction) GetHsend() *HashlockSend {
	if x, ok := x.GetValue().(*HashlockAction_Hsend); ok {
		return x.Hsend
	}
	return nil
}

func (x *HashlockAction) GetHunlock() *HashlockUnlock {
	if x, ok := x.GetValue().(*HashlockAction_Hunlock); ok {
		return x.Hunlock
	}
	return nil
}

func (x *HashlockAction) GetTy() int32 {
	if x != nil {
		return x.Ty
	}
	return 0
}

type isHashlockAction_Value interface {
	isHashlockAction_Value()
}

type HashlockAction_Hlock struct {
	Hlock *HashlockLock `protobuf:"bytes,1,opt,name=hlock,proto3,oneof"`
}

type HashlockAction_Hsend struct {
	Hsend *HashlockSend `protobuf:"bytes,2,opt,name=hsend,proto3,oneof"`
}

type HashlockAction_Hunlock struct {
	Hunlock *HashlockUnlock `protobuf:"bytes,3,opt,name=hunlock,proto3,oneof"`
}

func (*HashlockAction_Hlock) isHashlockAction_Value() {}

func (*HashlockAction_Hsend) isHashlockAction_Value() {}

func (*HashlockAction_Hunlock) isHashlockAction_Value() {}

var File_hashlock_proto protoreflect.FileDescriptor

var file_hashlock_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x68, 0x61, 0x73, 0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73, 0x22, 0xde, 0x01, 0x0a, 0x08, 0x48, 0x61, 0x73, 0x68,
	0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x1e, 0x0a, 0x0a, 0x68, 0x61, 0x73, 0x68, 0x6c, 0x6f, 0x63, 0x6b,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x68, 0x61, 0x73, 0x68, 0x6c, 0x6f,
	0x63, 0x6b, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a, 0x0a,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x74, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65,
	0x74, 0x75, 0x72, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x66, 0x72, 0x6f, 0x7a,
	0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x66, 0x72,
	0x6f, 0x7a, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x92, 0x01, 0x0a, 0x0c, 0x48, 0x61, 0x73,
	0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x4c, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x6f, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x6f,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65, 0x74, 0x75, 0x72,
	0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x26, 0x0a,
	0x0c, 0x48, 0x61, 0x73, 0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x95, 0x01, 0x0a, 0x0d, 0x48, 0x61, 0x73, 0x68, 0x6c, 0x6f,
	0x63, 0x6b, 0x71, 0x75, 0x65, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x63,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0b, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x62, 0x0a,
	0x08, 0x48, 0x61, 0x73, 0x68, 0x52, 0x65, 0x63, 0x76, 0x12, 0x1e, 0x0a, 0x0a, 0x48, 0x61, 0x73,
	0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x48,
	0x61, 0x73, 0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x0b, 0x49, 0x6e, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x52, 0x0b, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x28, 0x0a, 0x0e, 0x48, 0x61, 0x73, 0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x55, 0x6e, 0x6c,
	0x6f, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0xb6, 0x01, 0x0a, 0x0e,
	0x48, 0x61, 0x73, 0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2b,
	0x0a, 0x05, 0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x4c, 0x6f,
	0x63, 0x6b, 0x48, 0x00, 0x52, 0x05, 0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x2b, 0x0a, 0x05, 0x68,
	0x73, 0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x65, 0x6e, 0x64, 0x48,
	0x00, 0x52, 0x05, 0x68, 0x73, 0x65, 0x6e, 0x64, 0x12, 0x31, 0x0a, 0x07, 0x68, 0x75, 0x6e, 0x6c,
	0x6f, 0x63, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x6c, 0x6f, 0x63, 0x6b, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b,
	0x48, 0x00, 0x52, 0x07, 0x68, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x74,
	0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x74, 0x79, 0x42, 0x07, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2e, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_hashlock_proto_rawDescOnce sync.Once
	file_hashlock_proto_rawDescData = file_hashlock_proto_rawDesc
)

func file_hashlock_proto_rawDescGZIP() []byte {
	file_hashlock_proto_rawDescOnce.Do(func() {
		file_hashlock_proto_rawDescData = protoimpl.X.CompressGZIP(file_hashlock_proto_rawDescData)
	})
	return file_hashlock_proto_rawDescData
}

var file_hashlock_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_hashlock_proto_goTypes = []interface{}{
	(*Hashlock)(nil),       // 0: types.Hashlock
	(*HashlockLock)(nil),   // 1: types.HashlockLock
	(*HashlockSend)(nil),   // 2: types.HashlockSend
	(*Hashlockquery)(nil),  // 3: types.Hashlockquery
	(*HashRecv)(nil),       // 4: types.HashRecv
	(*HashlockUnlock)(nil), // 5: types.HashlockUnlock
	(*HashlockAction)(nil), // 6: types.HashlockAction
}
var file_hashlock_proto_depIdxs = []int32{
	3, // 0: types.HashRecv.Information:type_name -> types.Hashlockquery
	1, // 1: types.HashlockAction.hlock:type_name -> types.HashlockLock
	2, // 2: types.HashlockAction.hsend:type_name -> types.HashlockSend
	5, // 3: types.HashlockAction.hunlock:type_name -> types.HashlockUnlock
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_hashlock_proto_init() }
func file_hashlock_proto_init() {
	if File_hashlock_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_hashlock_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Hashlock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hashlock_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HashlockLock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hashlock_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HashlockSend); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hashlock_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Hashlockquery); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hashlock_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HashRecv); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hashlock_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HashlockUnlock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hashlock_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HashlockAction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_hashlock_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*HashlockAction_Hlock)(nil),
		(*HashlockAction_Hsend)(nil),
		(*HashlockAction_Hunlock)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_hashlock_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_hashlock_proto_goTypes,
		DependencyIndexes: file_hashlock_proto_depIdxs,
		MessageInfos:      file_hashlock_proto_msgTypes,
	}.Build()
	File_hashlock_proto = out.File
	file_hashlock_proto_rawDesc = nil
	file_hashlock_proto_goTypes = nil
	file_hashlock_proto_depIdxs = nil
}
