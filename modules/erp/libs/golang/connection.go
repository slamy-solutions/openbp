package erp

import (
	"errors"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	catalogGRPC "github.com/slamy-solutions/openbp/modules/erp/libs/golang/core/catalog"
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

// Connect to IOT Core service
func NewCoreConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, *CoreService, error) {
	dial, err := makeGrpcDial(address, opts...)
	if err != nil {
		return nil, nil, errors.New("failed to connect to service: " + err.Error())
	}

	return dial, &CoreService{
		Catalog: CatalogServices{
			Catalog: catalogGRPC.NewCatalogServiceClient(dial),
			Entry:   catalogGRPC.NewCatalogEntryServiceClient(dial),
		},
	}, nil
}
