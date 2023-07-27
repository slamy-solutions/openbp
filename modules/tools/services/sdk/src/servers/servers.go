package servers

import (
	"github.com/sirupsen/logrus"
	nativeModule "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	systemModule "github.com/slamy-solutions/openbp/modules/system/libs/golang"

	sdkTools "github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/tools"

	sdkGRPC "github.com/slamy-solutions/openbp/modules/tools/libs/golang/sdk/sdk"
	sdkServer "github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/servers/sdk"

	"google.golang.org/grpc"
)

func RegisterGRPCServers(grpcServer *grpc.Server, authHandler sdkTools.AuthHandler, logger *logrus.Entry, systemStub *systemModule.SystemStub, nativeStub *nativeModule.NativeStub) error {
	sdkServer := sdkServer.NewSDKServer(authHandler, logger.WithField("server", "sdk"), nativeStub)
	sdkGRPC.RegisterSDKServiceServer(grpcServer, sdkServer)

	return nil
}
