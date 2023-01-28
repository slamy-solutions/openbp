package services

import (
	"google.golang.org/grpc"

	tools "github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/tools"

	native "github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/services/native"
)

func RegisterGRPCServices(services *tools.ModulesStub, grpcServer *grpc.Server) error {
	return native.RegisterNativeServices(services, grpcServer)
}
