package services

import (
	"time"

	nativeIAmPolicyGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	policy_managed_no       = "NO"
	policy_managed_builtin  = "BUILTIN"
	policy_managed_identity = "IDENTITY"
	policy_managed_service  = "SERVICE"
)

const (
	policy_builtin_global_root    = "GLOBAL_ROOT"
	policy_builtin_namespace_root = "NAMESPACE_ROOT"
	policy_builtin_empty          = "EMPTY"
)

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

type policyInMongo struct {
	UUID primitive.ObjectID `bson:"_id,omitempty"`

	Name        string `bson:"name"`
	Description string `bson:"description"`

	Managed *managementTypeInMongo `bson:"managed"`

	NamespaceIndependent bool     `bson:"namespaceIndependent"`
	Resources            []string `bson:"resources"`
	Actions              []string `bson:"actions"`

	Tags []string `bson:"tags"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version uint64    `bson:"version"`
}

func ManagedTypeToGRPC(roleType string) nativeIAmPolicyGRPC.BuiltInPolicyType {
	switch roleType {
	case policy_builtin_global_root:
		return nativeIAmPolicyGRPC.BuiltInPolicyType_GLOBAL_ROOT
	case policy_builtin_namespace_root:
		return nativeIAmPolicyGRPC.BuiltInPolicyType_NAMESPACE_ROOT
	case policy_builtin_empty:
		return nativeIAmPolicyGRPC.BuiltInPolicyType_EMPTY
	}

	return nativeIAmPolicyGRPC.BuiltInPolicyType_EMPTY
}

func (p *policyInMongo) ToGRPCPolicy(namespace string) *nativeIAmPolicyGRPC.Policy {
	grpcPolicy := &nativeIAmPolicyGRPC.Policy{
		Namespace:   namespace,
		Uuid:        p.UUID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		Managed:     &nativeIAmPolicyGRPC.Policy_No{No: &nativeIAmPolicyGRPC.NotManagedData{}},
		Resources:   p.Resources,
		Actions:     p.Actions,
		Tags:        p.Tags,
		Created:     timestamppb.New(p.Created),
		Updated:     timestamppb.New(p.Updated),
		Version:     p.Version,
	}

	if p.Managed != nil {
		switch p.Managed.ManagementType {
		case policy_managed_builtin:
			grpcPolicy.Managed = &nativeIAmPolicyGRPC.Policy_BuiltIn{
				BuiltIn: &nativeIAmPolicyGRPC.BuiltInManagedData{
					Type: ManagedTypeToGRPC(p.Managed.BuiltInType),
				},
			}
		case policy_managed_identity:
			grpcPolicy.Managed = &nativeIAmPolicyGRPC.Policy_Identity{
				Identity: &nativeIAmPolicyGRPC.IdentityManagedData{
					IdentityNamespace: p.Managed.IdentityNamespace,
					IdentityUUID:      p.Managed.IdentityUUID,
				},
			}
		case policy_managed_service:
			grpcPolicy.Managed = &nativeIAmPolicyGRPC.Policy_Service{
				Service: &nativeIAmPolicyGRPC.ServiceManagedData{
					Service:      p.Managed.ServiceName,
					Reason:       p.Managed.ServiceReason,
					ManagementId: p.Managed.ServiceManagementID,
				},
			}
		}
	}

	return grpcPolicy
}
