package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
)

type CreateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *CreateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMPolicyService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *CreateTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestCreateTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTestSuite))
}

func (s *CreateTestSuite) TestParamsValidation() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	for _, tc := range createParamsValidationTestCases {
		s.Run(tc.testName, func() {
			r, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
				Namespace:            "",
				Name:                 "",
				Description:          "",
				Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
				NamespaceIndependent: false,
				Resources:            []string{},
				Actions:              []string{},
			})
			if err != nil {
				defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: r.Policy.Uuid})
			}

			assert.Equal(s.T(), tc.ok, err == nil)
			if err != nil {
				e, ok := status.FromError(err)
				require.True(s.T(), ok)
				require.Equal(s.T(), e.Code(), codes.InvalidArgument)
			}
		})
	}
}
