package sdk

import (
	"google.golang.org/grpc"

	nativeNamespace "github.com/slamy-solutions/openbp/modules/tools/libs/golang/native/namespace"
)

type NativeStub struct {
	Namespace nativeNamespace.NativeNamespaceServiceClient
}

type OpenBPStub struct {
	Native NativeStub
}

func NewStubForConnection(grpcDial *grpc.ClientConn) *OpenBPStub {
	return &OpenBPStub{
		Native: NativeStub{
			Namespace: nativeNamespace.NewNativeNamespaceServiceClient(grpcDial),
		},
	}
}
