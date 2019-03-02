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

type zgogrpc struct {
	res GrpcResourcer
}

type Grpcer interface {
	NewGrpc() (*zgogrpc, error) //初始化方法
	Run(ctx context.Context, s *grpc.Server, port string) error
	Server(ctx context.Context) (*grpc.Server, error)
}

func GetGrpc() (*zgogrpc, error) {
	return &zgogrpc{
		res: NewGrpcResourcer(),
	}, nil
}

func Grpc() Grpcer {
	return &zgogrpc{
		res: NewGrpcResourcer(),
	}
}

func (e *zgogrpc) NewGrpc() (*zgogrpc, error) {
	return GetGrpc()
}

func (e *zgogrpc) Run(ctx context.Context, s *grpc.Server, port string) error {
	return e.res.Run(ctx, s, port)
}

func (e *zgogrpc) Server(ctx context.Context) (*grpc.Server, error) {
	return e.res.Server(ctx)
}
