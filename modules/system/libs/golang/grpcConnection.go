package system

import (
	"time"

	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func makeGrpcClient[T interface{}](clientFunction func(grpc.ClientConnInterface) T, address string, opts ...grpc.DialOption) (*grpc.ClientConn, T, error) {
	opts = append(
		[]grpc.DialOption{
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
			grpc.WithBlock(),
			grpc.WithTimeout(time.Second * 5),
		},
		opts...,
	)

	dial, err := grpc.Dial(address, opts...)
	if err != nil {
		var result T
		return nil, result, err
	}

	client := clientFunction(dial)
	return dial, client, nil
}

// Connect to Vault service
func NewVaultConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, vault.VaultServiceClient, error) {
	return makeGrpcClient(vault.NewVaultServiceClient, address, opts...)
}
