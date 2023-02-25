package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type UpdateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *UpdateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMPolicyService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *UpdateTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateTestSuite))
}

func (s *UpdateTestSuite) TestUpdateInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	description := tools.GetRandomString(20)
	resources := []string{tools.GetRandomString(20)}
	actions := []string{tools.GetRandomString(20)}

	createResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 name,
		Description:          description,
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            resources,
		Actions:              actions,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createResponse.Policy.Uuid})

	newName := tools.GetRandomString(20)
	newDescription := tools.GetRandomString(20)
	newResources := []string{tools.GetRandomString(20)}
	newActions := []string{tools.GetRandomString(20)}

	updateResponse, err := s.nativeStub.Services.IamPolicy.Update(ctx, &policy.UpdatePolicyRequest{
		Namespace:            "",
		Uuid:                 createResponse.Policy.Uuid,
		Name:                 newName,
		Description:          newDescription,
		NamespaceIndependent: true,
		Resources:            newResources,
		Actions:              newActions,
	})

	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, updateResponse.Policy.Name)
	require.Equal(s.T(), newDescription, updateResponse.Policy.Description)
	require.Len(s.T(), updateResponse.Policy.Actions, 1)
	require.Equal(s.T(), newActions[0], updateResponse.Policy.Actions[0])
	require.Len(s.T(), updateResponse.Policy.Resources, 1)
	require.Equal(s.T(), newResources[0], updateResponse.Policy.Resources[0])

	getResponse, err := s.nativeStub.Services.IamPolicy.Get(ctx, &policy.GetPolicyRequest{
		Namespace: "",
		Uuid:      createResponse.Policy.Uuid,
		UseCache:  true,
	})

	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, getResponse.Policy.Name)
	require.Equal(s.T(), newDescription, getResponse.Policy.Description)
	require.Len(s.T(), getResponse.Policy.Actions, 1)
	require.Equal(s.T(), newActions[0], getResponse.Policy.Actions[0])
	require.Len(s.T(), getResponse.Policy.Resources, 1)
	require.Equal(s.T(), newResources[0], getResponse.Policy.Resources[0])
}

func (s *UpdateTestSuite) TestUpdateInNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName,
		FullName:    "",
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	name := tools.GetRandomString(20)
	description := tools.GetRandomString(20)
	resources := []string{tools.GetRandomString(20)}
	actions := []string{tools.GetRandomString(20)}

	createResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 name,
		Description:          description,
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            resources,
		Actions:              actions,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createResponse.Policy.Uuid})

	newName := tools.GetRandomString(20)
	newDescription := tools.GetRandomString(20)
	newResources := []string{tools.GetRandomString(20)}
	newActions := []string{tools.GetRandomString(20)}

	updateResponse, err := s.nativeStub.Services.IamPolicy.Update(ctx, &policy.UpdatePolicyRequest{
		Namespace:            namespaceName,
		Uuid:                 createResponse.Policy.Uuid,
		Name:                 newName,
		Description:          newDescription,
		NamespaceIndependent: true,
		Resources:            newResources,
		Actions:              newActions,
	})

	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, updateResponse.Policy.Name)
	require.Equal(s.T(), newDescription, updateResponse.Policy.Description)
	require.Len(s.T(), updateResponse.Policy.Actions, 1)
	require.Equal(s.T(), newActions[0], updateResponse.Policy.Actions[0])
	require.Len(s.T(), updateResponse.Policy.Resources, 1)
	require.Equal(s.T(), newResources[0], updateResponse.Policy.Resources[0])
	require.True(s.T(), updateResponse.Policy.NamespaceIndependent)

	getResponse, err := s.nativeStub.Services.IamPolicy.Get(ctx, &policy.GetPolicyRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Policy.Uuid,
		UseCache:  true,
	})

	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, getResponse.Policy.Name)
	require.Equal(s.T(), newDescription, getResponse.Policy.Description)
	require.Len(s.T(), getResponse.Policy.Actions, 1)
	require.Equal(s.T(), newActions[0], getResponse.Policy.Actions[0])
	require.Len(s.T(), getResponse.Policy.Resources, 1)
	require.Equal(s.T(), newResources[0], getResponse.Policy.Resources[0])
	require.True(s.T(), getResponse.Policy.NamespaceIndependent)
}

func (s *UpdateTestSuite) TestFailsWithNotFoundErrorWhenPolicyDoesntExistInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	description := tools.GetRandomString(20)
	resources := []string{tools.GetRandomString(20)}
	actions := []string{tools.GetRandomString(20)}

	createResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 name,
		Description:          description,
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            resources,
		Actions:              actions,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamPolicy.Update(ctx, &policy.UpdatePolicyRequest{
		Namespace:            "",
		Uuid:                 primitive.NewObjectID().Hex(),
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		NamespaceIndependent: true,
		Resources:            []string{tools.GetRandomString(20)},
		Actions:              []string{tools.GetRandomString(20)},
	})

	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *UpdateTestSuite) TestFailsWithNotFoundErrorWhenPolicyDoesntExistInNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName,
		FullName:    "",
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	name := tools.GetRandomString(20)
	description := tools.GetRandomString(20)
	resources := []string{tools.GetRandomString(20)}
	actions := []string{tools.GetRandomString(20)}

	createResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 name,
		Description:          description,
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            resources,
		Actions:              actions,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamPolicy.Update(ctx, &policy.UpdatePolicyRequest{
		Namespace:            namespaceName,
		Uuid:                 primitive.NewObjectID().Hex(),
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		NamespaceIndependent: true,
		Resources:            []string{tools.GetRandomString(20)},
		Actions:              []string{tools.GetRandomString(20)},
	})

	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *UpdateTestSuite) TestFailsWithNotFoundErrorWhenNamespaceDoesntExist() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName,
		FullName:    "",
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	name := tools.GetRandomString(20)
	description := tools.GetRandomString(20)
	resources := []string{tools.GetRandomString(20)}
	actions := []string{tools.GetRandomString(20)}

	createResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 name,
		Description:          description,
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            resources,
		Actions:              actions,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamPolicy.Update(ctx, &policy.UpdatePolicyRequest{
		Namespace:            tools.GetRandomString(20),
		Uuid:                 primitive.NewObjectID().Hex(),
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		NamespaceIndependent: true,
		Resources:            []string{tools.GetRandomString(20)},
		Actions:              []string{tools.GetRandomString(20)},
	})

	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}
