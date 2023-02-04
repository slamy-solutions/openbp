package smoke

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

func TestIAMRole(t *testing.T) {
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithIAMPolicyService().WithIAMRoleService())
	err := nativeStub.Connect()
	require.Nil(t, err)
	defer nativeStub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	name := tools.GetRandomString(20)
	description := tools.GetRandomString(20)

	roleCreateResponse, err := nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        name,
		Description: description,
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(t, err)
	defer nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: roleCreateResponse.Role.Uuid})
	assert.Equal(t, name, roleCreateResponse.Role.Name)
	assert.Equal(t, description, roleCreateResponse.Role.Description)

	// Create policy and add it to the role
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

	_, err = nativeStub.Services.IamRole.AddPolicy(ctx, &role.AddPolicyRequest{
		RoleNamespace:   "",
		RoleUUID:        roleCreateResponse.Role.Uuid,
		PolicyNamespace: "",
		PolicyUUID:      policyCreateResponse.Policy.Uuid,
	})
	require.Nil(t, err)

	// Get role and check if information is actual
	roleGetResponse, err := nativeStub.Services.IamRole.Get(ctx, &role.GetRoleRequest{
		Namespace: "",
		Uuid:      roleCreateResponse.Role.Uuid,
		UseCache:  true,
	})
	require.Nil(t, err)
	assert.Equal(t, name, roleGetResponse.Role.Name)
	assert.Equal(t, description, roleGetResponse.Role.Description)
	require.Len(t, roleGetResponse.Role.Policies, 1)
	assert.Equal(t, policyCreateResponse.Policy.Uuid, roleGetResponse.Role.Policies[0].Uuid)
	assert.Equal(t, "", roleGetResponse.Role.Policies[0].Namespace)

	// Delete role
	roleDeleteResponse, err := nativeStub.Services.IamRole.Delete(ctx, &role.DeleteRoleRequest{
		Namespace: "",
		Uuid:      roleCreateResponse.Role.Uuid,
	})
	require.Nil(t, err)
	assert.True(t, roleDeleteResponse.Existed)
}
