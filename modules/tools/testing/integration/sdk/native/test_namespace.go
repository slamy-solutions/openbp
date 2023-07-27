package native

/*
import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	tools "github.com/slamy-solutions/openbp/modules/tools/testing/tools"

	sdk "github.com/slamy-solutions/openbp/modules/tools/libs/golang"
	namespace_grpc "github.com/slamy-solutions/openbp/modules/tools/libs/golang/native/namespace"
)

type NamespaceTestSuite struct {
	suite.Suite

	sdk *sdk.OpenBPStub
}

func (suite *NamespaceTestSuite) SetupTest() {

}

func (suite *NamespaceTestSuite) TearDownSuite() {

}

func (s *NamespaceTestSuite) TestEnsureParameters() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	testCases := []struct {
		testName    string
		name        string
		fullName    string
		description string
		ok          bool
	}{
		{"All OK", tools.GetRandomString(20), "Some big name", "123", true},
		{"Max name length is 32", tools.GetRandomString(33), "Some big name", "123", false},
		{"Invalid characters in name", tools.GetRandomString(16) + "%", "Some big name", "123", false},
		{"Name must not be empty", "", "Some big name", "123", false},
		{"Empty fullName is ok", tools.GetRandomString(20), "", "123", true},
		{"Max fullName length is 128", tools.GetRandomString(20), tools.GetRandomString(130), "123", false},
		{"Empty description is ok", tools.GetRandomString(20), "alalal", "", true},
		{"Max description length is 512", tools.GetRandomString(20), tools.GetRandomString(10), tools.GetRandomString(530), false},
	}

	for _, tc := range testCases {
		s.Run(tc.testName, func() {
			response, err := s.sdk.Native.Namespace.Ensure(ctx, &namespace_grpc.EnsureNamespaceRequest{
				Name:        tc.name,
				FullName:    tc.fullName,
				Description: tc.description,
			})
			defer s.sdk.Native.Namespace.Delete(context.Background(), &namespace_grpc.DeleteNamespaceRequest{Name: tc.name})

		})
	}

}

func TestNamespace(t *testing.T) {
	suite.Run(t, new(NamespaceTestSuite))
}
*/
