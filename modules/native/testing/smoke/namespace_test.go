package smoke

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

func TestNamespace(t *testing.T) {
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService())
	err := nativeStub.Connect()
	require.Nil(t, err)
	defer nativeStub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	fullName := tools.GetRandomString(30)
	description := tools.GetRandomString(30)

	_, err = nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        name,
		FullName:    fullName,
		Description: description,
	})
	defer nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})

	require.Nil(t, err)

	r, err := nativeStub.Services.Namespace.Get(ctx, &namespace.GetNamespaceRequest{Name: name})
	require.Nil(t, err)

	assert.Equal(t, r.Namespace.Name, name)
	assert.Equal(t, r.Namespace.FullName, fullName)
	assert.Equal(t, r.Namespace.Description, description)

	r2, err := nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})
	require.Nil(t, err)
	require.True(t, r2.Existed)
}
