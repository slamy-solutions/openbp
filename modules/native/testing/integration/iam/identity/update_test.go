package identity

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
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type UpdateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *UpdateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
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

	createResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace: "",
		Name:      name,
		Managed:   &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: createResponse.Identity.Uuid})

	newName := tools.GetRandomString(20)

	updateResponse, err := s.nativeStub.Services.IAM.Identity.Update(ctx, &identity.UpdateIdentityRequest{
		Namespace: "",
		Uuid:      createResponse.Identity.Uuid,
		NewName:   newName,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, updateResponse.Identity.Name)

	getResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      createResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, getResponse.Identity.Name)
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

	createResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace: namespaceName,
		Name:      name,
		Managed:   &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: createResponse.Identity.Uuid})

	newName := tools.GetRandomString(20)

	updateResponse, err := s.nativeStub.Services.IAM.Identity.Update(ctx, &identity.UpdateIdentityRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Identity.Uuid,
		NewName:   newName,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, updateResponse.Identity.Name)

	getResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, getResponse.Identity.Name)
}

func (s *UpdateTestSuite) TestUpdateNonExistingIdentityInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace: "",
		Name:      tools.GetRandomString(20),
		Managed:   &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: createResponse.Identity.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.Update(ctx, &identity.UpdateIdentityRequest{
		Namespace: "",
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}

func (s *UpdateTestSuite) TestUpdateNonExistingIdentityInNamespace() {
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

	createResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace: namespaceName,
		Name:      tools.GetRandomString(20),
		Managed:   &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: createResponse.Identity.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.Update(ctx, &identity.UpdateIdentityRequest{
		Namespace: "",
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}

func (s *UpdateTestSuite) TestUpdateNonExistingIdentityInNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.nativeStub.Services.IAM.Identity.Update(ctx, &identity.UpdateIdentityRequest{
		Namespace: tools.GetRandomString(20),
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}
