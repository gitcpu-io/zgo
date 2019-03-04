/*
@Time : 2019-02-27 14:38
@Author : zhangjianguo
@File : service
@Software: GoLand
*/
package zgogrpc

import (
	"context"
	"google.golang.org/grpc"
)

const (
	rpc_port = ":80"
)

type zgogrpc struct {
	res GrpcResourcer
}

type Grpcer interface {
	Run(ctx context.Context, s *grpc.Server) (string, error)
	Server(ctx context.Context) (*grpc.Server, error)
	Client(ctx context.Context, ip string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
}

func GetGrpc() *zgogrpc {
	return &zgogrpc{
		res: NewGrpcResourcer(),
	}
}

func (e *zgogrpc) Run(ctx context.Context, s *grpc.Server) (string, error) {
	return e.res.Run(ctx, s, rpc_port)
}

func (e *zgogrpc) Server(ctx context.Context) (*grpc.Server, error) {
	return e.res.Server(ctx), nil
}

func (e *zgogrpc) Client(ctx context.Context, ip string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return e.res.Client(ctx, ip+rpc_port, opts...)
}
