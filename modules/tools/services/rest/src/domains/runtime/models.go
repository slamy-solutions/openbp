package runtime

import (
	runtimeGRPC "github.com/slamy-solutions/openbp/modules/runtime/libs/golang/manager/runtime"
)

type formatedRuntime struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Run       bool   `json:"run"`
}

func formatedRuntimeFromGRPC(grpcRuntime *runtimeGRPC.Runtime) formatedRuntime {
	return formatedRuntime{
		Namespace: grpcRuntime.Namespace,
		Name:      grpcRuntime.Name,
		Run:       grpcRuntime.Run,
	}
}
