package services

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/slamy-solutions/open-erp/modules/system/libs/go/cache"

	nativeIAmAuthGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/auth/src/grpc/native_iam_auth"
	nativeIAmIdentityGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/auth/src/grpc/native_iam_identity"
	nativeIAmPolicyGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/auth/src/grpc/native_iam_policy"
	nativeNamespaceGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/auth/src/grpc/native_namespace"
)

type IAmAuthServer struct {
	nativeIAmAuthGRPC.UnimplementedIAMAuthServiceServer

	mongoClient             *mongo.Client
	mongoDbPrefix           string
	cacheClient             cache.Cache
	nativeNamespaceClient   nativeNamespaceGRPC.NamespaceServiceClient
	nativeIAmPolicyClient   nativeIAmPolicyGRPC.IAMPolicyServiceClient
	nativeIAmIdentityClient nativeIAmIdentityGRPC.IAMIdentityServiceClient
}
