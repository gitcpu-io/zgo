/*
  Grpc 客户端
*/
package zgogrpc

import (
	"context"
	"google.golang.org/grpc"
)

const (
	rpc_port = ":50051"
)

type zgogrpc struct {
	res GrpcResourcer
}

/*
 Grpcer 对外使用接口
*/
type Grpcer interface {
	// 获取Grpc 服务端
	// 默认端口50051
	Server(ctx context.Context) (*grpc.Server, error)
	// 客户端通过获取到的Grpc server 将客户端接口实现注册到servers 交给Run启动
	Run(ctx context.Context, s *grpc.Server) (string, error)
	// 调用GPC服务
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
