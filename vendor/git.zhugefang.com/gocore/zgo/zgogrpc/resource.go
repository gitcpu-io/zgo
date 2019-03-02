/*
@Time : 2019-02-27 14:38
@Author : zhangjianguo
@File : resource
@Software: GoLand
*/
package zgogrpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GrpcResourcer interface {
	Server(ctx context.Context) (*grpc.Server, error)
	Run(ctx context.Context, s *grpc.Server, port string) error
}
type grpcResource struct {
	res GrpcResourcer //使用resource另外的一个接口
}

func NewGrpcResourcer() GrpcResourcer {
	return &grpcResource{}
}

func (e *grpcResource) Server(ctx context.Context) (*grpc.Server, error) {
	return grpc.NewServer(), nil
}

func (e *grpcResource) Run(ctx context.Context, s *grpc.Server, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}
