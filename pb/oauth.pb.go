// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/oauth.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type CheckTokenRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CheckTokenRequest) Reset()         { *m = CheckTokenRequest{} }
func (m *CheckTokenRequest) String() string { return proto.CompactTextString(m) }
func (*CheckTokenRequest) ProtoMessage()    {}
func (*CheckTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a064a29c17b5838, []int{0}
}

func (m *CheckTokenRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckTokenRequest.Unmarshal(m, b)
}
func (m *CheckTokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckTokenRequest.Marshal(b, m, deterministic)
}
func (m *CheckTokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckTokenRequest.Merge(m, src)
}
func (m *CheckTokenRequest) XXX_Size() int {
	return xxx_messageInfo_CheckTokenRequest.Size(m)
}
func (m *CheckTokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckTokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CheckTokenRequest proto.InternalMessageInfo

func (m *CheckTokenRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type CheckTokenResponse struct {
	UserDetails          *UserDetails   `protobuf:"bytes,1,opt,name=userDetails,proto3" json:"userDetails,omitempty"`
	ClientDetails        *ClientDetails `protobuf:"bytes,2,opt,name=clientDetails,proto3" json:"clientDetails,omitempty"`
	IsValidToken         bool           `protobuf:"varint,3,opt,name=isValidToken,proto3" json:"isValidToken,omitempty"`
	Err                  string         `protobuf:"bytes,4,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *CheckTokenResponse) Reset()         { *m = CheckTokenResponse{} }
func (m *CheckTokenResponse) String() string { return proto.CompactTextString(m) }
func (*CheckTokenResponse) ProtoMessage()    {}
func (*CheckTokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a064a29c17b5838, []int{1}
}

func (m *CheckTokenResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckTokenResponse.Unmarshal(m, b)
}
func (m *CheckTokenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckTokenResponse.Marshal(b, m, deterministic)
}
func (m *CheckTokenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckTokenResponse.Merge(m, src)
}
func (m *CheckTokenResponse) XXX_Size() int {
	return xxx_messageInfo_CheckTokenResponse.Size(m)
}
func (m *CheckTokenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckTokenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CheckTokenResponse proto.InternalMessageInfo

func (m *CheckTokenResponse) GetUserDetails() *UserDetails {
	if m != nil {
		return m.UserDetails
	}
	return nil
}

func (m *CheckTokenResponse) GetClientDetails() *ClientDetails {
	if m != nil {
		return m.ClientDetails
	}
	return nil
}

func (m *CheckTokenResponse) GetIsValidToken() bool {
	if m != nil {
		return m.IsValidToken
	}
	return false
}

func (m *CheckTokenResponse) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type ClientDetails struct {
	ClientId                    string   `protobuf:"bytes,1,opt,name=clientId,proto3" json:"clientId,omitempty"`
	AccessTokenValiditySeconds  int32    `protobuf:"varint,2,opt,name=accessTokenValiditySeconds,proto3" json:"accessTokenValiditySeconds,omitempty"`
	RefreshTokenValiditySeconds int32    `protobuf:"varint,3,opt,name=refreshTokenValiditySeconds,proto3" json:"refreshTokenValiditySeconds,omitempty"`
	AuthorizedGrantTypes        []string `protobuf:"bytes,4,rep,name=authorizedGrantTypes,proto3" json:"authorizedGrantTypes,omitempty"`
	XXX_NoUnkeyedLiteral        struct{} `json:"-"`
	XXX_unrecognized            []byte   `json:"-"`
	XXX_sizecache               int32    `json:"-"`
}

func (m *ClientDetails) Reset()         { *m = ClientDetails{} }
func (m *ClientDetails) String() string { return proto.CompactTextString(m) }
func (*ClientDetails) ProtoMessage()    {}
func (*ClientDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a064a29c17b5838, []int{2}
}

func (m *ClientDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientDetails.Unmarshal(m, b)
}
func (m *ClientDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientDetails.Marshal(b, m, deterministic)
}
func (m *ClientDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientDetails.Merge(m, src)
}
func (m *ClientDetails) XXX_Size() int {
	return xxx_messageInfo_ClientDetails.Size(m)
}
func (m *ClientDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientDetails.DiscardUnknown(m)
}

var xxx_messageInfo_ClientDetails proto.InternalMessageInfo

func (m *ClientDetails) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *ClientDetails) GetAccessTokenValiditySeconds() int32 {
	if m != nil {
		return m.AccessTokenValiditySeconds
	}
	return 0
}

func (m *ClientDetails) GetRefreshTokenValiditySeconds() int32 {
	if m != nil {
		return m.RefreshTokenValiditySeconds
	}
	return 0
}

func (m *ClientDetails) GetAuthorizedGrantTypes() []string {
	if m != nil {
		return m.AuthorizedGrantTypes
	}
	return nil
}

type UserDetails struct {
	UserId               int64    `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Username             string   `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Authorities          []string `protobuf:"bytes,3,rep,name=authorities,proto3" json:"authorities,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserDetails) Reset()         { *m = UserDetails{} }
func (m *UserDetails) String() string { return proto.CompactTextString(m) }
func (*UserDetails) ProtoMessage()    {}
func (*UserDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a064a29c17b5838, []int{3}
}

func (m *UserDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserDetails.Unmarshal(m, b)
}
func (m *UserDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserDetails.Marshal(b, m, deterministic)
}
func (m *UserDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserDetails.Merge(m, src)
}
func (m *UserDetails) XXX_Size() int {
	return xxx_messageInfo_UserDetails.Size(m)
}
func (m *UserDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_UserDetails.DiscardUnknown(m)
}

var xxx_messageInfo_UserDetails proto.InternalMessageInfo

func (m *UserDetails) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *UserDetails) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *UserDetails) GetAuthorities() []string {
	if m != nil {
		return m.Authorities
	}
	return nil
}

func init() {
	proto.RegisterType((*CheckTokenRequest)(nil), "pb.CheckTokenRequest")
	proto.RegisterType((*CheckTokenResponse)(nil), "pb.CheckTokenResponse")
	proto.RegisterType((*ClientDetails)(nil), "pb.ClientDetails")
	proto.RegisterType((*UserDetails)(nil), "pb.UserDetails")
}

func init() { proto.RegisterFile("pb/oauth.proto", fileDescriptor_8a064a29c17b5838) }

var fileDescriptor_8a064a29c17b5838 = []byte{
	// 362 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0x4f, 0x4f, 0xfa, 0x40,
	0x10, 0x4d, 0x29, 0x10, 0x98, 0xc2, 0xef, 0x27, 0x1b, 0x24, 0x0d, 0x5e, 0x9a, 0x9e, 0xf0, 0x52,
	0x23, 0x1e, 0x3c, 0x90, 0x18, 0x15, 0x13, 0x63, 0x3c, 0x98, 0x2c, 0xe8, 0xc1, 0x5b, 0xff, 0x8c,
	0xe9, 0x06, 0x6c, 0xeb, 0xee, 0xd6, 0x04, 0x3f, 0x98, 0xdf, 0xc7, 0x6f, 0x62, 0x76, 0x5b, 0xa1,
	0x10, 0xc2, 0x6d, 0xdf, 0xec, 0x9b, 0xb7, 0xf3, 0x76, 0x1e, 0xfc, 0xcb, 0x82, 0xb3, 0xd4, 0xcf,
	0x65, 0xec, 0x65, 0x3c, 0x95, 0x29, 0xa9, 0x65, 0x81, 0x7b, 0x0a, 0xbd, 0x69, 0x8c, 0xe1, 0x62,
	0x9e, 0x2e, 0x30, 0xa1, 0xf8, 0x91, 0xa3, 0x90, 0xa4, 0x0f, 0x0d, 0xa9, 0xb0, 0x6d, 0x38, 0xc6,
	0xa8, 0x4d, 0x0b, 0xe0, 0x7e, 0x1b, 0x40, 0xaa, 0x5c, 0x91, 0xa5, 0x89, 0x40, 0x72, 0x0e, 0x56,
	0x2e, 0x90, 0xdf, 0xa1, 0xf4, 0xd9, 0x52, 0xe8, 0x16, 0x6b, 0xfc, 0xdf, 0xcb, 0x02, 0xef, 0x79,
	0x53, 0xa6, 0x55, 0x0e, 0xb9, 0x84, 0x6e, 0xb8, 0x64, 0x98, 0xc8, 0xbf, 0xa6, 0x9a, 0x6e, 0xea,
	0xa9, 0xa6, 0x69, 0xf5, 0x82, 0x6e, 0xf3, 0x88, 0x0b, 0x1d, 0x26, 0x5e, 0xfc, 0x25, 0x8b, 0xf4,
	0x0c, 0xb6, 0xe9, 0x18, 0xa3, 0x16, 0xdd, 0xaa, 0x91, 0x23, 0x30, 0x91, 0x73, 0xbb, 0xae, 0x47,
	0x57, 0x47, 0xf7, 0xc7, 0x80, 0xee, 0x96, 0x2c, 0x19, 0x42, 0xab, 0x10, 0x7e, 0x88, 0x4a, 0x8f,
	0x6b, 0x4c, 0xae, 0x60, 0xe8, 0x87, 0x21, 0x0a, 0xa1, 0xe5, 0xb4, 0x30, 0x93, 0xab, 0x19, 0x86,
	0x69, 0x12, 0x15, 0x93, 0x36, 0xe8, 0x01, 0x06, 0xb9, 0x86, 0x13, 0x8e, 0x6f, 0x1c, 0x45, 0xbc,
	0x57, 0xc0, 0xd4, 0x02, 0x87, 0x28, 0x64, 0x0c, 0x7d, 0xb5, 0xa5, 0x94, 0xb3, 0x2f, 0x8c, 0xee,
	0xb9, 0x9f, 0xc8, 0xf9, 0x2a, 0x43, 0x61, 0xd7, 0x1d, 0x73, 0xd4, 0xa6, 0x7b, 0xef, 0xdc, 0x10,
	0xac, 0xca, 0x77, 0x93, 0x01, 0x34, 0xd5, 0x87, 0x97, 0xf6, 0x4c, 0x5a, 0x22, 0x65, 0x5c, 0x9d,
	0x12, 0xff, 0x1d, 0xb5, 0x95, 0x36, 0x5d, 0x63, 0xe2, 0x80, 0x55, 0x4a, 0x4b, 0x86, 0x6a, 0x50,
	0xf5, 0x5a, 0xb5, 0x34, 0x7e, 0x84, 0xce, 0xd3, 0x4d, 0x2e, 0xe3, 0x19, 0xf2, 0x4f, 0x16, 0x22,
	0x99, 0x00, 0x6c, 0x02, 0x41, 0x8e, 0xf5, 0xfa, 0x76, 0xc3, 0x34, 0x1c, 0xec, 0x96, 0x8b, 0xdc,
	0xdc, 0x36, 0x5f, 0xeb, 0xde, 0x24, 0x0b, 0x82, 0xa6, 0x0e, 0xe3, 0xc5, 0x6f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x5c, 0x07, 0x51, 0x14, 0x9e, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// OAuthServiceClient is the client API for OAuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OAuthServiceClient interface {
	CheckToken(ctx context.Context, in *CheckTokenRequest, opts ...grpc.CallOption) (*CheckTokenResponse, error)
}

type oAuthServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOAuthServiceClient(cc grpc.ClientConnInterface) OAuthServiceClient {
	return &oAuthServiceClient{cc}
}

func (c *oAuthServiceClient) CheckToken(ctx context.Context, in *CheckTokenRequest, opts ...grpc.CallOption) (*CheckTokenResponse, error) {
	out := new(CheckTokenResponse)
	err := c.cc.Invoke(ctx, "/pb.OAuthService/CheckToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OAuthServiceServer is the server API for OAuthService service.
type OAuthServiceServer interface {
	CheckToken(context.Context, *CheckTokenRequest) (*CheckTokenResponse, error)
}

// UnimplementedOAuthServiceServer can be embedded to have forward compatible implementations.
type UnimplementedOAuthServiceServer struct {
}

func (*UnimplementedOAuthServiceServer) CheckToken(ctx context.Context, req *CheckTokenRequest) (*CheckTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckToken not implemented")
}

func RegisterOAuthServiceServer(s *grpc.Server, srv OAuthServiceServer) {
	s.RegisterService(&_OAuthService_serviceDesc, srv)
}

func _OAuthService_CheckToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OAuthServiceServer).CheckToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.OAuthService/CheckToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OAuthServiceServer).CheckToken(ctx, req.(*CheckTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OAuthService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.OAuthService",
	HandlerType: (*OAuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckToken",
			Handler:    _OAuthService_CheckToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/oauth.proto",
}
