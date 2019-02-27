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
	Server(ctx context.Context, port string, sd *grpc.ServiceDesc, ss interface{}) (interface{}, error)
	Client(ctx context.Context, args map[string]interface{}) (interface{}, error)
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

func (e *zgogrpc) Server(ctx context.Context, port string, sd *grpc.ServiceDesc, ss interface{}) (interface{}, error) {
	return e.res.Server(ctx, port, sd, ss)
}
func (e *zgogrpc) Client(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return e.res.Client(ctx, args)
}
