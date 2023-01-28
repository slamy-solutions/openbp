package native

import (
	"google.golang.org/grpc"

	tools "github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/tools"
)

func RegisterNativeServices(modules *tools.ModulesStub, grpcServer *grpc.Server) error {
	RegisterNamespaceService(modules, grpcServer)
	return nil
}
