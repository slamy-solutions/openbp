package smoke

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

func TestIAMIdentity(t *testing.T) {
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithIAMPolicyService().WithIAMRoleService().WithIAMIdentityService())
	err := nativeStub.Connect()
	require.Nil(t, err)
	defer nativeStub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	name := tools.GetRandomString(20)
	initiallyActive := true

	identityCreateResponse, err := nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            name,
		InitiallyActive: initiallyActive,
	})
	require.Nil(t, err)
	defer nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	assert.Equal(t, name, identityCreateResponse.Identity.Name)
	assert.Equal(t, initiallyActive, identityCreateResponse.Identity.Active)

	// Create role and add it to the identity
	roleCreateResponse, err := nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(t, err)
	defer nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: roleCreateResponse.Role.Uuid})
	_, err = nativeStub.Services.IamIdentity.AddRole(ctx, &identity.AddRoleRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		RoleNamespace:     "",
		RoleUUID:          roleCreateResponse.Role.Uuid,
	})
	require.Nil(t, err)

	// Create policy and add it to the identity
	policyCreateResponse, err := nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(t, err)
	defer nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: policyCreateResponse.Policy.Uuid})
	_, err = nativeStub.Services.IamIdentity.AddPolicy(ctx, &identity.AddPolicyRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		PolicyNamespace:   "",
		PolicyUUID:        policyCreateResponse.Policy.Uuid,
	})
	require.Nil(t, err)

	// Get the identity and validate if information is actual
	identityGetResponse, err := nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  true,
	})
	require.Nil(t, err)
	assert.Equal(t, name, identityGetResponse.Identity.Name)
	assert.Equal(t, initiallyActive, identityGetResponse.Identity.Active)
	require.Len(t, identityGetResponse.Identity.Policies, 1)
	assert.Equal(t, "", identityGetResponse.Identity.Policies[0].Namespace)
	assert.Equal(t, policyCreateResponse.Policy.Uuid, identityGetResponse.Identity.Policies[0].Uuid)
	require.Len(t, identityGetResponse.Identity.Roles, 1)
	assert.Equal(t, "", identityGetResponse.Identity.Roles[0].Namespace)
	assert.Equal(t, roleCreateResponse.Role.Uuid, identityGetResponse.Identity.Roles[0].Uuid)

	// Delete identity
	identityDeleteResponse, err := nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
	})
	require.Nil(t, err)
	assert.True(t, identityDeleteResponse.Existed)
}
