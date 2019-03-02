/*
@Time : 2019-02-27 14:38
@Author : zhangjianguo
@File : resource
@Software: GoLand
*/
package zgogrpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type GrpcResourcer interface {
	Server(ctx context.Context, port string, sd *grpc.ServiceDesc, ss interface{}) (interface{}, error)
	Client(ctx context.Context, args map[string]interface{}) (interface{}, error)
}
type grpcResource struct {
	res GrpcResourcer //使用resource另外的一个接口
}

func NewGrpcResourcer() GrpcResourcer {

	return &grpcResource{}
}

func (e *grpcResource) Server(ctx context.Context, port string, sd *grpc.ServiceDesc, ss interface{}) (interface{}, error) {

	list, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Print(err)
	}
	s := grpc.NewServer()
	s.RegisterService(sd, ss)

	err = s.Serve(list)
	return port, err
}
func (e *grpcResource) Client(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return nil, nil
}
