// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sql/sem/semenumpb/constraint.proto

package semenumpb

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// ForeignKeyAction describes the action which should be taken when a foreign
// key constraint reference is acted upon.
type ForeignKeyAction int32

const (
	ForeignKeyAction_NO_ACTION   ForeignKeyAction = 0
	ForeignKeyAction_RESTRICT    ForeignKeyAction = 1
	ForeignKeyAction_SET_NULL    ForeignKeyAction = 2
	ForeignKeyAction_SET_DEFAULT ForeignKeyAction = 3
	ForeignKeyAction_CASCADE     ForeignKeyAction = 4
)

var ForeignKeyAction_name = map[int32]string{
	0: "NO_ACTION",
	1: "RESTRICT",
	2: "SET_NULL",
	3: "SET_DEFAULT",
	4: "CASCADE",
}

var ForeignKeyAction_value = map[string]int32{
	"NO_ACTION":   0,
	"RESTRICT":    1,
	"SET_NULL":    2,
	"SET_DEFAULT": 3,
	"CASCADE":     4,
}

func (x ForeignKeyAction) String() string {
	return proto.EnumName(ForeignKeyAction_name, int32(x))
}

func (ForeignKeyAction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8e3e948e264df145, []int{0}
}

// Match is the algorithm used to compare composite keys.
type Match int32

const (
	Match_SIMPLE  Match = 0
	Match_FULL    Match = 1
	Match_PARTIAL Match = 2
)

var Match_name = map[int32]string{
	0: "SIMPLE",
	1: "FULL",
	2: "PARTIAL",
}

var Match_value = map[string]int32{
	"SIMPLE":  0,
	"FULL":    1,
	"PARTIAL": 2,
}

func (x Match) String() string {
	return proto.EnumName(Match_name, int32(x))
}

func (Match) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8e3e948e264df145, []int{1}
}

func init() {
	proto.RegisterEnum("cockroach.sql.sem.semenumpb.ForeignKeyAction", ForeignKeyAction_name, ForeignKeyAction_value)
	proto.RegisterEnum("cockroach.sql.sem.semenumpb.Match", Match_name, Match_value)
}

func init() {
	proto.RegisterFile("sql/sem/semenumpb/constraint.proto", fileDescriptor_8e3e948e264df145)
}

var fileDescriptor_8e3e948e264df145 = []byte{
	// 289 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2a, 0x2e, 0xcc, 0xd1,
	0x2f, 0x4e, 0xcd, 0x05, 0xe1, 0xd4, 0xbc, 0xd2, 0xdc, 0x82, 0x24, 0xfd, 0xe4, 0xfc, 0xbc, 0xe2,
	0x92, 0xa2, 0xc4, 0xcc, 0xbc, 0x12, 0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0xe9, 0xe4, 0xfc,
	0xe4, 0xec, 0xa2, 0xfc, 0xc4, 0xe4, 0x0c, 0xbd, 0xe2, 0xc2, 0x1c, 0xbd, 0xe2, 0xd4, 0x5c, 0x3d,
	0xb8, 0x6a, 0x29, 0x91, 0xf4, 0xfc, 0xf4, 0x7c, 0xb0, 0x3a, 0x7d, 0x10, 0x0b, 0xa2, 0x45, 0x2b,
	0x9a, 0x4b, 0xc0, 0x2d, 0xbf, 0x28, 0x35, 0x33, 0x3d, 0xcf, 0x3b, 0xb5, 0xd2, 0x31, 0xb9, 0x24,
	0x33, 0x3f, 0x4f, 0x88, 0x97, 0x8b, 0xd3, 0xcf, 0x3f, 0xde, 0xd1, 0x39, 0xc4, 0xd3, 0xdf, 0x4f,
	0x80, 0x41, 0x88, 0x87, 0x8b, 0x23, 0xc8, 0x35, 0x38, 0x24, 0xc8, 0xd3, 0x39, 0x44, 0x80, 0x11,
	0xc4, 0x0b, 0x76, 0x0d, 0x89, 0xf7, 0x0b, 0xf5, 0xf1, 0x11, 0x60, 0x12, 0xe2, 0xe7, 0xe2, 0x06,
	0xf1, 0x5c, 0x5c, 0xdd, 0x1c, 0x43, 0x7d, 0x42, 0x04, 0x98, 0x85, 0xb8, 0xb9, 0xd8, 0x9d, 0x1d,
	0x83, 0x9d, 0x1d, 0x5d, 0x5c, 0x05, 0x58, 0xb4, 0xb4, 0xb8, 0x58, 0x7d, 0x13, 0x4b, 0x92, 0x33,
	0x84, 0xb8, 0xb8, 0xd8, 0x82, 0x3d, 0x7d, 0x03, 0x7c, 0x5c, 0x05, 0x18, 0x84, 0x38, 0xb8, 0x58,
	0xdc, 0x40, 0x9a, 0x19, 0x41, 0x6a, 0x03, 0x1c, 0x83, 0x42, 0x3c, 0x1d, 0x7d, 0x04, 0x98, 0x9c,
	0x22, 0x4e, 0x3c, 0x94, 0x63, 0x38, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x1b, 0x8f,
	0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b,
	0x8f, 0xe5, 0x18, 0xa2, 0xcc, 0xd2, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5,
	0xe1, 0x1e, 0x4d, 0x49, 0x42, 0xb0, 0xf5, 0x0b, 0xb2, 0xd3, 0xf5, 0x31, 0x82, 0x29, 0x89, 0x0d,
	0xec, 0x53, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x07, 0x88, 0xb6, 0x42, 0x01, 0x00,
	0x00,
}
