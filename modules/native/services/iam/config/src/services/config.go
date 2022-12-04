package services

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/cache"
	grpccodes "google.golang.org/grpc/codes"

	iamGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/config"
	keyValueStorageGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
)

type IAMConfigServer struct {
	iamGRPC.UnimplementedIAMConfigServiceServer

	cacheClient           cache.Cache
	keyvaluestorageClient keyValueStorageGRPC.KeyValueStorageServiceClient
}

const (
	STORAGE_KEY = "native_iam_config_global"
)

var DEFAULT_CONFIG = iamGRPC.Configuration{
	AccessTokenTTL:  15 * 60 * 1000,
	RefreshTokenTTL: 48 * 60 * 60 * 1000,
	PasswordAuth: &iamGRPC.Configuration_PasswordAuth{
		Enabled:           true,
		AllowRegistration: false,
	},
	GoogleOAuth2: &iamGRPC.Configuration_OAuth2{
		Enabled:           false,
		ClientId:          "",
		ClientSecret:      "",
		AllowRegistration: false,
	},
	FacebookOAuth2: &iamGRPC.Configuration_OAuth2{
		Enabled:           false,
		ClientId:          "",
		ClientSecret:      "",
		AllowRegistration: false,
	},
	GithubOAuth2: &iamGRPC.Configuration_OAuth2{
		Enabled:           false,
		ClientId:          "",
		ClientSecret:      "",
		AllowRegistration: false,
	},
	GitlabOAuth2: &iamGRPC.Configuration_OAuth2{
		Enabled:           false,
		ClientId:          "",
		ClientSecret:      "",
		AllowRegistration: false,
	},
}

func NewIAMConfigServer(cacheClient cache.Cache, keyvaluestorageClient keyValueStorageGRPC.KeyValueStorageServiceClient) *IAMConfigServer {
	return &IAMConfigServer{
		cacheClient:           cacheClient,
		keyvaluestorageClient: keyvaluestorageClient,
	}
}

func (s *IAMConfigServer) Set(ctx context.Context, in *iamGRPC.SetConfigRequest) (*iamGRPC.SetConfigResponse, error) {
	data, err := proto.Marshal(in.Configuration)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	_, err = s.keyvaluestorageClient.Set(ctx, &keyValueStorageGRPC.SetRequest{
		Namespace: "",
		Key:       STORAGE_KEY,
		Value:     data,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &iamGRPC.SetConfigResponse{}, status.Error(grpccodes.OK, "")
}

func (s *IAMConfigServer) Get(ctx context.Context, in *iamGRPC.GetConfigRequest) (*iamGRPC.GetConfigresponse, error) {
	r, err := s.keyvaluestorageClient.Get(ctx, &keyValueStorageGRPC.GetRequest{
		Namespace: "",
		Key:       STORAGE_KEY,
		UseCache:  in.UseCache,
	})
	if err != nil {
		code := grpc.Code(err)
		if code == grpccodes.NotFound {
			return &iamGRPC.GetConfigresponse{Configuration: &DEFAULT_CONFIG}, status.Error(grpccodes.OK, "")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	var config iamGRPC.Configuration
	err = proto.Unmarshal(r.Value, &config)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &iamGRPC.GetConfigresponse{Configuration: &config}, status.Error(grpccodes.OK, "")
}
