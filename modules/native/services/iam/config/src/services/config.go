package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/status"

	"github.com/slamy-solutions/open-erp/modules/system/libs/go/cache"

	iamGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/config/src/grpc/native_iam"
)

type IAMConfigServer struct {
	iamGRPC.UnimplementedIAMConfigServiceServer

	mongoClient           *mongo.Client
	mongoDbPrefix         string
	mongoBundleCollection *mongo.Collection
	cacheClient           cache.Cache
	bigCacheClient        cache.Cache
}

func (s *IAMConfigServer) Set(ctx context.Context, in *iamGRPC.SetConfigRequest) (*iamGRPC.SetConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}

func (s *IAMConfigServer) Get(ctx context.Context, in *iamGRPC.GetConfigRequest) (*iamGRPC.GetConfigresponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
