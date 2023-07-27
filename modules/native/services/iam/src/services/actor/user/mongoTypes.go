package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	nativeActorUserGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/actor/user"
)

type UserInMongo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Login    string             `bson:"login"`
	Identity string             `bson:"identity"`

	FullName string `bson:"fullName"`
	Avatar   string `bson:"avatar"`
	Email    string `bson:"email"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version uint64    `bson:"version"`
}

func (u *UserInMongo) ToGRPCUser(namespace string) *nativeActorUserGRPC.User {
	return &nativeActorUserGRPC.User{
		Namespace: namespace,
		Uuid:      u.ID.Hex(),
		Login:     u.Login,
		Identity:  u.Identity,
		FullName:  u.FullName,
		Avatar:    u.Avatar,
		Email:     u.Email,
		Created:   timestamppb.New(u.Created),
		Updated:   timestamppb.New(u.Updated),
		Version:   u.Version,
	}
}
