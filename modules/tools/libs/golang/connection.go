package sdk

import (
	"time"

	"google.golang.org/grpc"
)

type Connector interface {
	MakeConnection() (*grpc.ClientConn, error)
}

type simpleTokenConnector struct {
}

func makeGrpcConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(
		[]grpc.DialOption{
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
			grpc.WithBlock(),
			grpc.WithTimeout(time.Second * 15),
		},
		opts...,
	)

	return grpc.Dial(address, opts...)
}
