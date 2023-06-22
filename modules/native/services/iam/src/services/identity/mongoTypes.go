package identity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	nativeIAmIdentityGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
)

const (
	identity_managed_no       = "NO"
	identity_managed_identity = "IDENTITY"
	identity_managed_service  = "SERVICE"
)

type identityPolicyInMongo struct {
	Namespace string `bson:"namespace"`
	UUID      string `bson:"uuid"`
}

type identityRoleInMongo struct {
	Namespace string `bson:"namespace"`
	UUID      string `bson:"uuid"`
}

type managementTypeInMongo struct {
	ManagementType string `bson:"_managementType"`

	// Builyin

	BuiltInType string `bson:"builtInType,omitempty"`

	// Identity

	IdentityNamespace string `bson:"identityNamespace,omitempty"`
	IdentityUUID      string `bson:"identityUUID,omitempty"`

	// Service

	ServiceName         string `bson:"serviceName,omitempty"`
	ServiceReason       string `bson:"serviceReason,omitempty"`
	ServiceManagementID string `bson:"serviceManagementId,omitempty"`
}

type identityInMongo struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `bson:"name"`
	Active bool               `bson:"active"`

	Managed *managementTypeInMongo `bson:"managed"`

	Policies []identityPolicyInMongo `bson:"policies"`
	Roles    []identityRoleInMongo   `bson:"roles"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version uint64    `bson:"version"`
}

func (i *identityInMongo) ToGRPCIdentity(namespace string) *nativeIAmIdentityGRPC.Identity {
	policies := make([]*nativeIAmIdentityGRPC.Identity_PolicyReference, len(i.Policies))
	for index, policy := range i.Policies {
		policies[index] = &nativeIAmIdentityGRPC.Identity_PolicyReference{
			Namespace: policy.Namespace,
			Uuid:      policy.UUID,
		}
	}

	roles := make([]*nativeIAmIdentityGRPC.Identity_RoleReference, len(i.Roles))
	for index, role := range i.Roles {
		roles[index] = &nativeIAmIdentityGRPC.Identity_RoleReference{
			Namespace: role.Namespace,
			Uuid:      role.UUID,
		}
	}

	identity := &nativeIAmIdentityGRPC.Identity{
		Namespace: namespace,
		Uuid:      i.ID.Hex(),
		Name:      i.Name,
		Active:    i.Active,
		Managed:   &nativeIAmIdentityGRPC.Identity_No{No: &nativeIAmIdentityGRPC.NotManagedData{}},
		Policies:  policies,
		Roles:     roles,
		Created:   timestamppb.New(i.Created),
		Updated:   timestamppb.New(i.Updated),
		Version:   i.Version,
	}

	if i.Managed != nil {
		switch i.Managed.ManagementType {
		case identity_managed_identity:
			identity.Managed = &nativeIAmIdentityGRPC.Identity_Identity{
				Identity: &nativeIAmIdentityGRPC.IdentityManagedData{
					IdentityNamespace: i.Managed.IdentityNamespace,
					IdentityUUID:      i.Managed.IdentityUUID,
				},
			}
		case identity_managed_service:
			identity.Managed = &nativeIAmIdentityGRPC.Identity_Service{
				Service: &nativeIAmIdentityGRPC.ServiceManagedData{
					Service:      i.Managed.ServiceName,
					Reason:       i.Managed.ServiceReason,
					ManagementId: i.Managed.ServiceManagementID,
				},
			}
		}
	}

	return identity
}
