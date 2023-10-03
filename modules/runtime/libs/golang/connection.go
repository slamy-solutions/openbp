package iot

import (
	"errors"

	"github.com/slamy-solutions/openbp/modules/runtime/libs/golang/manager/rpc"
	"github.com/slamy-solutions/openbp/modules/runtime/libs/golang/manager/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func makeGrpcDial(address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(
		[]grpc.DialOption{
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
			// Disabled because module may not be enabled
			// grpc.WithBlock(),
			// grpc.WithTimeout(time.Second * 5),
		},
		opts...,
	)

	dial, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, errors.New("failed to establish grpc dial: " + err.Error())
	}

	return dial, nil
}

// Connect to Runtime Manager service
func NewManagerConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, *ManagerService, error) {
	dial, err := makeGrpcDial(address, opts...)
	if err != nil {
		return nil, nil, errors.New("failed to connect to service: " + err.Error())
	}

	return dial, &ManagerService{
		Runtime: runtime.NewRuntimeServiceClient(dial),
		RPC:     rpc.NewRPCServiceClient(dial),
	}, nil
}
