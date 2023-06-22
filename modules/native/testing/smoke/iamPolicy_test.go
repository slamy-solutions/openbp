package smoke

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

func TestIAMPolicy(t *testing.T) {
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithIAMService())
	err := nativeStub.Connect()
	require.Nil(t, err)
	defer nativeStub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	description := tools.GetRandomString(20)

	resource := tools.GetRandomString(30)
	action := tools.GetRandomString(30)

	createResponse, err := nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 name,
		Description:          description,
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{resource},
		Actions:              []string{action},
	})
	require.Nil(t, err)
	defer nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createResponse.Policy.Uuid})

	assert.Equal(t, name, createResponse.Policy.Name)
	assert.Equal(t, description, createResponse.Policy.Description)
	assert.False(t, createResponse.Policy.NamespaceIndependent)
	require.Len(t, createResponse.Policy.Resources, 1)
	assert.Equal(t, resource, createResponse.Policy.Resources[0])
	require.Len(t, createResponse.Policy.Actions, 1)
	assert.Equal(t, action, createResponse.Policy.Actions[0])

	getResponse, err := nativeStub.Services.IAM.Policy.Get(ctx, &policy.GetPolicyRequest{
		Namespace: "",
		Uuid:      createResponse.Policy.Uuid,
		UseCache:  true,
	})
	require.Nil(t, err)
	assert.Equal(t, name, getResponse.Policy.Name)
	assert.Equal(t, description, getResponse.Policy.Description)
	assert.False(t, getResponse.Policy.NamespaceIndependent)
	require.Len(t, getResponse.Policy.Resources, 1)
	assert.Equal(t, resource, getResponse.Policy.Resources[0])
	require.Len(t, getResponse.Policy.Actions, 1)
	assert.Equal(t, action, getResponse.Policy.Actions[0])

	deleteResponse, err := nativeStub.Services.IAM.Policy.Delete(ctx, &policy.DeletePolicyRequest{Namespace: "", Uuid: createResponse.Policy.Uuid})
	require.Nil(t, err)
	assert.True(t, deleteResponse.Existed)
}
