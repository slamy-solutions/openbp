package storage

import (
	"github.com/sirupsen/logrus"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
)

type fsRouter struct {
	nativeStub *native.NativeStub

	logger *logrus.Entry
}
