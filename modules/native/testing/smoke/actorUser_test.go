package smoke

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

func TestActorUser(t *testing.T) {
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithActorUserService())
	err := nativeStub.Connect()
	require.Nil(t, err)
	defer nativeStub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	login := tools.GetRandomString(20)
	fullName := tools.GetRandomString(20)
	avatar := tools.GetRandomString(20)
	email := tools.GetRandomString(20)

	userCreateResponse, err := nativeStub.Services.ActorUser.Create(ctx, &user.CreateRequest{
		Namespace: "",
		Login:     login,
		FullName:  fullName,
		Avatar:    avatar,
		Email:     email,
	})
	require.Nil(t, err)
	defer nativeStub.Services.ActorUser.Delete(context.Background(), &user.DeleteRequest{Namespace: "", Uuid: userCreateResponse.User.Uuid})

	assert.Equal(t, login, userCreateResponse.User.Login)
	assert.Equal(t, fullName, userCreateResponse.User.FullName)
	assert.Equal(t, avatar, userCreateResponse.User.Avatar)
	assert.Equal(t, email, userCreateResponse.User.Email)

	userGetResponse, err := nativeStub.Services.ActorUser.Get(ctx, &user.GetRequest{
		Namespace: "",
		Uuid:      userCreateResponse.User.Uuid,
		UseCache:  true,
	})
	require.Nil(t, err)

	assert.Equal(t, login, userGetResponse.User.Login)
	assert.Equal(t, fullName, userGetResponse.User.FullName)
	assert.Equal(t, avatar, userGetResponse.User.Avatar)
	assert.Equal(t, email, userGetResponse.User.Email)

	userDeleteResponse, err := nativeStub.Services.ActorUser.Delete(ctx, &user.DeleteRequest{
		Namespace: "",
		Uuid:      userCreateResponse.User.Uuid,
	})
	require.Nil(t, err)
	assert.True(t, userDeleteResponse.Existed)
}
