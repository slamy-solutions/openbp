package balena

import (
	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
)

type DevicesServer struct {
	balena.UnimplementedBalenaDevicesServiceServer

	logger *logrus.Entry
}

func NewDevicesServer(logger *logrus.Entry) *DevicesServer {
	return &DevicesServer{
		logger: logger,
	}
}
