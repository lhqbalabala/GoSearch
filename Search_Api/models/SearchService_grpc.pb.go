// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: SearchService.proto

package models

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SearchClient is the client API for Search service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchClient interface {
	GetSearchResult(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResult, error)
	GetSearchPictureResult(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchPictureResult, error)
	GetDocById(ctx context.Context, in *DocRequest, opts ...grpc.CallOption) (*DocResult, error)
}

type searchClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchClient(cc grpc.ClientConnInterface) SearchClient {
	return &searchClient{cc}
}

func (c *searchClient) GetSearchResult(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResult, error) {
	out := new(SearchResult)
	err := c.cc.Invoke(ctx, "/models.Search/GetSearchResult", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) GetSearchPictureResult(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchPictureResult, error) {
	out := new(SearchPictureResult)
	err := c.cc.Invoke(ctx, "/models.Search/GetSearchPictureResult", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) GetDocById(ctx context.Context, in *DocRequest, opts ...grpc.CallOption) (*DocResult, error) {
	out := new(DocResult)
	err := c.cc.Invoke(ctx, "/models.Search/GetDocById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchServer is the server API for Search service.
// All implementations must embed UnimplementedSearchServer
// for forward compatibility
type SearchServer interface {
	GetSearchResult(context.Context, *SearchRequest) (*SearchResult, error)
	GetSearchPictureResult(context.Context, *SearchRequest) (*SearchPictureResult, error)
	GetDocById(context.Context, *DocRequest) (*DocResult, error)
	mustEmbedUnimplementedSearchServer()
}

// UnimplementedSearchServer must be embedded to have forward compatible implementations.
type UnimplementedSearchServer struct {
}

func (UnimplementedSearchServer) GetSearchResult(context.Context, *SearchRequest) (*SearchResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSearchResult not implemented")
}
func (UnimplementedSearchServer) GetSearchPictureResult(context.Context, *SearchRequest) (*SearchPictureResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSearchPictureResult not implemented")
}
func (UnimplementedSearchServer) GetDocById(context.Context, *DocRequest) (*DocResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDocById not implemented")
}
func (UnimplementedSearchServer) mustEmbedUnimplementedSearchServer() {}

// UnsafeSearchServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchServer will
// result in compilation errors.
type UnsafeSearchServer interface {
	mustEmbedUnimplementedSearchServer()
}

func RegisterSearchServer(s grpc.ServiceRegistrar, srv SearchServer) {
	s.RegisterService(&Search_ServiceDesc, srv)
}

func _Search_GetSearchResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).GetSearchResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/models.Search/GetSearchResult",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).GetSearchResult(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_GetSearchPictureResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).GetSearchPictureResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/models.Search/GetSearchPictureResult",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).GetSearchPictureResult(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_GetDocById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DocRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).GetDocById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/models.Search/GetDocById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).GetDocById(ctx, req.(*DocRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Search_ServiceDesc is the grpc.ServiceDesc for Search service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Search_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "models.Search",
	HandlerType: (*SearchServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSearchResult",
			Handler:    _Search_GetSearchResult_Handler,
		},
		{
			MethodName: "GetSearchPictureResult",
			Handler:    _Search_GetSearchPictureResult_Handler,
		},
		{
			MethodName: "GetDocById",
			Handler:    _Search_GetDocById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "SearchService.proto",
}
