package tools

import (
	"sync"
	"time"

	"google.golang.org/grpc"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
)

var adminStubMutex *sync.Mutex = &sync.Mutex{}
var adminStub *native.NativeStub = nil

func CreateConnection() *grpc.ClientConn {
	opts := append(
		[]grpc.DialOption{
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
			grpc.WithBlock(),
			grpc.WithTimeout(time.Second * 15),
		},
	)

	return grpc.Dial(address, opts...)
}

func GetAdminStub() {
	adminStubMutex.Lock()
	defer adminStubMutex.Unlock()
}
