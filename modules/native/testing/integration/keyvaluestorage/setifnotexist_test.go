package namespace

import (
	"testing"

	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
)

type SetKeyIfNotExistTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *SetKeyIfNotExistTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithKeyValueStorageService().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *SetKeyIfNotExistTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestSetKeyIfNotExistTestSuite(t *testing.T) {
	suite.Run(t, new(SetKeyIfNotExistTestSuite))
}

func (s *SetKeyIfNotExistTestSuite) ValidatesInputs() {
	s.T().Skip("Not implemented")
}

func (s *SetKeyIfNotExistTestSuite) SetsValue() {
	s.T().Skip("Not implemented")
}

func (s *SetKeyIfNotExistTestSuite) SetsValueInNamespace() {
	s.T().Skip("Not implemented")
}

func (s *SetKeyIfNotExistTestSuite) DoesntUpdateValueIfAlreadyExist() {
	s.T().Skip("Not implemented")
}
