package models

import (
	"errors"

	settingsGRPC "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/settings"
)

type BackendType string

const (
	BackendTypeNative BackendType = "native"
	BackendType1C     BackendType = "1C"
)

var ErrUnknownBackendType = errors.New("unknown backend type")
var ErrBackendUnavailable = errors.New("backend unavailable")
var ErrBackendMissBehaviour = errors.New("backend miss behaviour")

type Backend interface {
	// Get type of the current backend
	GetType() BackendType

	ClientRepository() ClientRepository
	PerformerRepository() PerformerRepository
}

func BackendTypeToGRPC(backendType BackendType) settingsGRPC.BackendType {
	switch backendType {
	case BackendTypeNative:
		return settingsGRPC.BackendType_NATIVE
	case BackendType1C:
		return settingsGRPC.BackendType_ONE_C
	default:
		panic(ErrUnknownBackendType)
	}
}

func BackendTypeFromGRPC(backendType settingsGRPC.BackendType) BackendType {
	switch backendType {
	case settingsGRPC.BackendType_NATIVE:
		return BackendTypeNative
	case settingsGRPC.BackendType_ONE_C:
		return BackendType1C
	default:
		panic(ErrUnknownBackendType)
	}
}
