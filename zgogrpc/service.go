/*
  Grpc 客户端
*/
package zgogrpc

import (
  "context"
  "fmt"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/keepalive"
  "time"
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
  Server(ctx context.Context,opts ...grpc.ServerOption) (*grpc.Server, error)
  // 客户端通过获取到的Grpc server 将客户端接口实现注册到servers 交给Run启动
  Run(ctx context.Context, s *grpc.Server, port int) (int, error)
  // 调用GPC服务
  Client(ctx context.Context, ip string, port int, opts ...grpc.DialOption) (*grpc.ClientConn, error)

  WithInsecure() grpc.DialOption

  WithCallOptions(size int) grpc.DialOption

  WithKeepalive(maxConnectionIdle,maxConnectionAge,maxConnectionAgeGrace,time,timeout time.Duration) grpc.ServerOption
}

func GetGrpc() *zgogrpc {
  return &zgogrpc{
    res: NewGrpcResourcer(),
  }
}

func (e *zgogrpc) Run(ctx context.Context, s *grpc.Server, port int) (int, error) {
  return e.res.Run(ctx, s, port)
}

func (e *zgogrpc) Server(ctx context.Context,opts ...grpc.ServerOption) (*grpc.Server, error) {
  return e.res.Server(ctx,opts...), nil
}

func (e *zgogrpc) Client(ctx context.Context, ip string, port int, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
  return e.res.Client(ctx, fmt.Sprintf("%s:%d", ip, port), opts...)
}

func (e *zgogrpc) WithInsecure() grpc.DialOption {
  return grpc.WithTransportCredentials(insecure.NewCredentials())
}

func (e *zgogrpc) WithCallOptions(size int) grpc.DialOption {
  if size == 0 {
    size = 1024 * 1024 * 16
  }
  return grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(size),grpc.MaxCallSendMsgSize(size))
}

func (e *zgogrpc) WithKeepalive(maxConnectionIdle,maxConnectionAge,maxConnectionAgeGrace,time,timeout time.Duration) grpc.ServerOption {
  keep := keepalive.ServerParameters{
    MaxConnectionIdle:     maxConnectionIdle,
    MaxConnectionAge:      maxConnectionAge,
    MaxConnectionAgeGrace: maxConnectionAgeGrace,
    Time:                  time,
    Timeout:               timeout,
  }
  return grpc.KeepaliveParams(keep)
}

