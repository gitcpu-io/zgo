/*
  Grpc 客户端
*/
package zgogrpc

import (
  "context"
  "fmt"
  "google.golang.org/grpc"
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
  Run(ctx context.Context, s *grpc.Server, port string) (string, error)
  // 调用GPC服务
  Client(ctx context.Context, ip, port string, opts ...grpc.DialOption) (*grpc.ClientConn, error)

  WithInsecure() grpc.DialOption
}

func GetGrpc() *zgogrpc {
  return &zgogrpc{
    res: NewGrpcResourcer(),
  }
}

func (e *zgogrpc) Run(ctx context.Context, s *grpc.Server, port string) (string, error) {
  return e.res.Run(ctx, s, port)
}

func (e *zgogrpc) Server(ctx context.Context) (*grpc.Server, error) {
  return e.res.Server(ctx), nil
}

func (e *zgogrpc) Client(ctx context.Context, ip, port string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
  return e.res.Client(ctx, fmt.Sprintf("%s:%s", ip, port), opts...)
}

func (e *zgogrpc) WithInsecure() grpc.DialOption {
  return grpc.WithInsecure()
}
