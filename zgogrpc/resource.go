package zgogrpc

import (
  "context"
  "fmt"
  "google.golang.org/grpc"
  "log"
  "net"
)

type GrpcResourcer interface {
  Server(ctx context.Context) *grpc.Server
  Client(ctx context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
  Run(ctx context.Context, s *grpc.Server, port string) (string, error)
}
type grpcResource struct {
  res GrpcResourcer //使用resource另外的一个接口
}

func NewGrpcResourcer() GrpcResourcer {
  return &grpcResource{}
}

func (e *grpcResource) Server(ctx context.Context) *grpc.Server {
  return grpc.NewServer()
}
func (e *grpcResource) Client(ctx context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
  return grpc.Dial(target, opts...)

}
func (e *grpcResource) Run(ctx context.Context, s *grpc.Server, port string) (string, error) {
  lis, err := net.Listen("tcp", fmt.Sprintf("%s%s", ":", port))
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }
  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
  return port, err
}
