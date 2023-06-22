package policy

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type GetMultipleTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *GetMultipleTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *GetMultipleTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestGetMultipleTestSuite(t *testing.T) {
	suite.Run(t, new(GetMultipleTestSuite))
}

func (s *GetMultipleTestSuite) TestGetMultiple() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName1 := tools.GetRandomString(20)
	namespaceName2 := tools.GetRandomString(20)

	var policies = []*struct {
		namespace string
		create    bool
		search    bool
		uuid      string
	}{
		{"", true, true, primitive.NewObjectID().Hex()},
		{"", true, true, primitive.NewObjectID().Hex()},
		{"", true, false, primitive.NewObjectID().Hex()},
		{"", true, false, primitive.NewObjectID().Hex()},
		{"", false, true, primitive.NewObjectID().Hex()},
		{"", false, true, primitive.NewObjectID().Hex()},

		{namespaceName1, true, true, primitive.NewObjectID().Hex()},
		{namespaceName1, true, true, primitive.NewObjectID().Hex()},
		{namespaceName1, true, false, primitive.NewObjectID().Hex()},
		{namespaceName1, true, false, primitive.NewObjectID().Hex()},
		{namespaceName1, false, true, primitive.NewObjectID().Hex()},
		{namespaceName1, false, true, primitive.NewObjectID().Hex()},

		{namespaceName2, true, true, primitive.NewObjectID().Hex()},
		{namespaceName2, true, true, primitive.NewObjectID().Hex()},
		{namespaceName2, true, false, primitive.NewObjectID().Hex()},
		{namespaceName2, true, false, primitive.NewObjectID().Hex()},
		{namespaceName2, false, true, primitive.NewObjectID().Hex()},
		{namespaceName2, false, true, primitive.NewObjectID().Hex()},

		{tools.GetRandomString(20), false, true, primitive.NewObjectID().Hex()},
		{tools.GetRandomString(20), false, true, primitive.NewObjectID().Hex()},
	}

	//Create required namespaces
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName1,
		FullName:    "",
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(ctx, &namespace.DeleteNamespaceRequest{Name: namespaceName1})

	_, err = s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName2,
		FullName:    "",
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(ctx, &namespace.DeleteNamespaceRequest{Name: namespaceName2})

	// Make sure at the end all policies deleted
	defer func() {
		for _, p := range policies {
			if p.create {
				s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{
					Namespace: p.namespace,
					Uuid:      p.uuid,
				})
			}
		}
	}()

	// Create all the policies
	for _, p := range policies {
		if p.create {
			r, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
				Namespace:            p.namespace,
				Name:                 tools.GetRandomString(20),
				Description:          tools.GetRandomString(20),
				Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
				NamespaceIndependent: false,
				Resources:            []string{},
				Actions:              []string{},
			})
			require.Nil(s.T(), err)
			p.uuid = r.Policy.Uuid
		}
	}

	policiesToSearch := make([]*policy.GetMultiplePoliciesRequest_RequestedPolicy, 0, 12)
	for _, p := range policies {
		if p.search {
			policiesToSearch = append(policiesToSearch, &policy.GetMultiplePoliciesRequest_RequestedPolicy{
				Namespace: p.namespace,
				Uuid:      p.uuid,
			})
		}
	}
	require.Len(s.T(), policiesToSearch, 14)

	r, err := s.nativeStub.Services.IAM.Policy.GetMultiple(ctx, &policy.GetMultiplePoliciesRequest{
		Policies: policiesToSearch,
	})
	require.Nil(s.T(), err)

	receivedPolicies := make([]*policy.Policy, 0, 14)
	for {
		chunk, err := r.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}

		receivedPolicies = append(receivedPolicies, chunk.Policy)
	}
	require.Len(s.T(), receivedPolicies, 6)

	existInResponse := func(namespaceName string, uuid string) bool {
		for _, p := range receivedPolicies {
			if p.Namespace == namespaceName && p.Uuid == uuid {
				return true
			}
		}
		return false
	}

	existInRequest := func(namespaceName string, uuid string) bool {
		for _, p := range policiesToSearch {
			if p.Namespace == namespaceName && p.Uuid == uuid {
				return true
			}
		}
		return false
	}

	// Validating if request and result are same
	for _, p := range receivedPolicies {
		require.True(s.T(), existInRequest(p.Namespace, p.Uuid))
	}
	for _, p := range policies {
		if p.create && p.search {
			require.True(s.T(), existInResponse(p.namespace, p.uuid))
		}
	}
}
