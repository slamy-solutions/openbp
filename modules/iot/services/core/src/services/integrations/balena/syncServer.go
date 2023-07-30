package balena

import (
	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
)

type SyncServer struct {
	balena.UnimplementedBalenaSyncServiceServer

	logger *logrus.Entry
}

func NewSyncServer(logger *logrus.Entry) *SyncServer {
	return &SyncServer{
		logger: logger,
	}
}
