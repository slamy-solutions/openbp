package namespace

import (
	"testing"

	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
)

type SetKeyTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *SetKeyTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithKeyValueStorageService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *SetKeyTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestSetKeyTestSuite(t *testing.T) {
	suite.Run(t, new(SetKeyTestSuite))
}

func (s *SetKeyTestSuite) TestSetsKey() {

}
