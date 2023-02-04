package smoke

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

func TestKeyValueStorage(t *testing.T) {
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithKeyValueStorageService())
	err := nativeStub.Connect()
	require.Nil(t, err)
	defer nativeStub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := tools.GetRandomString(20)
	value := tools.GetRandomString(20)

	_, err = nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: "",
		Key:       key,
		Value:     []byte(value),
	})
	require.Nil(t, err)

	defer nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{
		Namespace: "",
		Key:       key,
	})

	getResponse, err := nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: "",
		Key:       key,
		UseCache:  true,
	})
	require.Nil(t, err)
	require.Equal(t, value, string(getResponse.Value))

	removeResponse, err := nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{
		Namespace: "",
		Key:       key,
	})
	require.Nil(t, err)
	require.True(t, removeResponse.Removed)
}
