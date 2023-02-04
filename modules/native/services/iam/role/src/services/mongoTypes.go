package services

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	nativeIAmRoleGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
)

const (
	role_managed_no       = "NO"
	role_managed_builtin  = "BUILTIN"
	role_managed_identity = "IDENTITY"
	role_managed_service  = "SERVICE"
)

const (
	role_builtin_global_root    = "GLOBAL_ROOT"
	role_builtin_namespace_root = "NAMESPACE_ROOT"
	role_builtin_empty          = "EMPTY"
)

type assignedPolicyInMongo struct {
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

type roleInMongo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`

	Managed *managementTypeInMongo `bson:"managed"`

	Policies []*assignedPolicyInMongo `bson:"policies"`

	Tags []string `bson:"tags"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version uint64    `bson:"version"`
}

func ManagedTypeToGRPC(roleType string) nativeIAmRoleGRPC.BuiltInRoleType {
	switch roleType {
	case role_builtin_global_root:
		return nativeIAmRoleGRPC.BuiltInRoleType_GLOBAL_ROOT
	case role_builtin_namespace_root:
		return nativeIAmRoleGRPC.BuiltInRoleType_NAMESPACE_ROOT
	case role_builtin_empty:
		return nativeIAmRoleGRPC.BuiltInRoleType_EMPTY
	}

	return nativeIAmRoleGRPC.BuiltInRoleType_EMPTY
}

func (r *roleInMongo) ToGRPCRole(namespace string) *nativeIAmRoleGRPC.Role {
	policies := make([]*nativeIAmRoleGRPC.AssignedPolicy, len(r.Policies))
	for policyIndex, policy := range r.Policies {
		policies[policyIndex] = &nativeIAmRoleGRPC.AssignedPolicy{
			Namespace: policy.Namespace,
			Uuid:      policy.UUID,
		}
	}

	grpcRole := &nativeIAmRoleGRPC.Role{
		Namespace:   namespace,
		Uuid:        r.ID.Hex(),
		Name:        r.Name,
		Description: r.Description,

		Policies: policies,
		Managed:  &nativeIAmRoleGRPC.Role_No{No: &nativeIAmRoleGRPC.NotManagedData{}},

		Tags: r.Tags,

		Created: timestamppb.New(r.Created),
		Updated: timestamppb.New(r.Updated),
		Version: r.Version,
	}

	if r.Managed != nil {
		switch r.Managed.ManagementType {
		case role_managed_builtin:
			grpcRole.Managed = &nativeIAmRoleGRPC.Role_BuiltIn{
				BuiltIn: &nativeIAmRoleGRPC.BuiltInManagedData{
					Type: ManagedTypeToGRPC(r.Managed.BuiltInType),
				},
			}
		case role_managed_identity:
			grpcRole.Managed = &nativeIAmRoleGRPC.Role_Identity{
				Identity: &nativeIAmRoleGRPC.IdentityManagedData{
					IdentityNamespace: r.Managed.IdentityNamespace,
					IdentityUUID:      r.Managed.IdentityUUID,
				},
			}
		case role_managed_service:
			grpcRole.Managed = &nativeIAmRoleGRPC.Role_Service{
				Service: &nativeIAmRoleGRPC.ServiceManagedData{
					Service:      r.Managed.ServiceName,
					Reason:       r.Managed.ServiceReason,
					ManagementId: r.Managed.ServiceManagementID,
				},
			}
		}
	}

	return grpcRole
}
