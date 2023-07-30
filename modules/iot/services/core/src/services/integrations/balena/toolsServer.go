package balena

import (
	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
)

type ToolsServer struct {
	balena.UnimplementedBalenaToolsServiceServer

	logger *logrus.Entry
}

func NewToolsServer(logger *logrus.Entry) *ToolsServer {
	return &ToolsServer{
		logger: logger,
	}
}
