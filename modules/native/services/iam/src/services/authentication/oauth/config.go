package oauth

import (
	"context"

	grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/oauth2"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConfigService struct {
	grpc.UnimplementedIAMAuthenticationOAuth2ConfigServiceServer

	systemStub *system.SystemStub
}

func NewConfigService(systemStub *system.SystemStub) *ConfigService {
	return &ConfigService{
		systemStub: systemStub,
	}
}

func encryptProviderClientSecret(ctx context.Context, systemStub *system.SystemStub, secret string) ([]byte, error) {
	response, err := systemStub.Vault.Encrypt(ctx, &vault.EncryptRequest{PlainData: []byte(secret)})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.FailedPrecondition {
			return nil, status.Error(codes.FailedPrecondition, "failed to encrypt client secret: vault is sealed: "+err.Error())
		}

		return nil, status.Error(codes.Internal, "failed to encrypt client secret: "+err.Error())
	}

	return response.EncryptedData, nil
}
func decryptProviderClientSecret(ctx context.Context, systemStub *system.SystemStub, secret []byte) (string, error) {
	response, err := systemStub.Vault.Decrypt(ctx, &vault.DecryptRequest{EncryptedData: secret})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.FailedPrecondition {
			return "", status.Error(codes.FailedPrecondition, "failed to decrypt client secret: vault is sealed: "+err.Error())
		}

		return "", status.Error(codes.Internal, "failed to decrypt client secret: "+err.Error())
	}

	return string(response.PlainData), nil
}

func (s *ConfigService) UpdateProviderConfig(ctx context.Context, in *grpc.UpdateProviderConfigRequest) (*grpc.UpdateProviderConfigResponse, error) {
	setQuery := bson.M{"enabled": in.Config.Enabled}
	if in.Config.ClientId != "" {
		setQuery["clientId"] = in.Config.ClientId
	}
	if in.Config.ClientSecret != "" {
		encryptedSecret, err := encryptProviderClientSecret(ctx, s.systemStub, in.Config.ClientSecret)
		if err != nil {
			return nil, err
		}

		setQuery["clientSecret"] = encryptedSecret
	}
	if in.Config.AuthUrl != "" {
		setQuery["authUrl"] = in.Config.AuthUrl
	}
	if in.Config.TokenUrl != "" {
		setQuery["tokenUrl"] = in.Config.TokenUrl
	}
	if in.Config.UserApiUrl != "" {
		setQuery["userApiUrl"] = in.Config.UserApiUrl
	}

	collection := configCollectionByNamespace(s.systemStub, in.Config.Namespace)
	_, err := collection.UpdateOne(ctx, bson.M{"name": ProviderGRPCTypeToString[in.Config.Type]}, bson.M{"$set": setQuery, "$setOnInsert": bson.M{"name": ProviderGRPCTypeToString[in.Config.Type]}}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update provider config in database: %s", err.Error())
	}

	return &grpc.UpdateProviderConfigResponse{}, status.Error(codes.OK, "")
}
func (s *ConfigService) ListProviderConfigs(ctx context.Context, in *grpc.ListProviderConfigsRequest) (*grpc.ListProviderConfigsResponse, error) {
	collection := configCollectionByNamespace(s.systemStub, in.Namespace)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &grpc.ListProviderConfigsResponse{Configs: []*grpc.ProviderConfig{}}, status.Error(codes.OK, "Namespace not found")
			}
		}

		return nil, status.Errorf(codes.Internal, "failed to list provider configs from database: failed to init cursor: %s", err.Error())
	}

	var configs []authProviderConfigInMongo
	if err = cursor.All(ctx, &configs); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list provider configs from database: failed to read from cursor: %s", err.Error())
	}

	grpcConfigs := make([]*grpc.ProviderConfig, 0, len(configs))
	for _, config := range configs {
		grpcConfigs = append(grpcConfigs, config.ToGRPCProviderConfig(in.Namespace))
	}

	return &grpc.ListProviderConfigsResponse{Configs: grpcConfigs}, status.Error(codes.OK, "")
}
func (s *ConfigService) GetAvailableProviders(ctx context.Context, in *grpc.GetAvailableProvidersRequest) (*grpc.GetAvailableProvidersResponse, error) {
	collection := configCollectionByNamespace(s.systemStub, in.Namespace)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &grpc.GetAvailableProvidersResponse{Providers: []grpc.ProviderType{}}, status.Error(codes.OK, "Namespace not found")
			}
		}

		return nil, status.Errorf(codes.Internal, "failed to list provider configs from database: failed to init cursor: %s", err.Error())
	}

	var configs []authProviderConfigInMongo
	if err = cursor.All(ctx, &configs); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list provider configs from database: failed to read from cursor: %s", err.Error())
	}

	providers := make([]grpc.ProviderType, 0, len(configs))
	for _, config := range configs {
		if config.Enabled {
			providers = append(providers, ProviderStringTypeToGRPC[config.Name])
		}
	}

	return &grpc.GetAvailableProvidersResponse{Providers: providers}, status.Error(codes.OK, "")
}
