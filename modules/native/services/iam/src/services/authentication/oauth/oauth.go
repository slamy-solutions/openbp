package oauth

import (
	"context"
	"errors"
	"log"
	"time"

	grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/oauth2"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/authentication/oauth/provider"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errOAuthProviderDisabled = errors.New("provider disabled")
var errOAuthErrorWhileRetrievingToken = errors.New("error while retrieving token")
var errOAuthErrorWhileFetchingUserDetails = errors.New("error while fetching user details")

type OAuthServer struct {
	grpc.UnimplementedIAMAuthenticationOAuth2ServiceServer

	systemStub *system.SystemStub
}

func NewOAuthServer(ctx context.Context, systemStub *system.SystemStub) (*OAuthServer, error) {
	err := EnsureIndexesForNamespace(ctx, "", systemStub)
	if err != nil {
		return nil, errors.Join(errors.New("failed to ensure indexes for global namespace"), err)
	}

	return &OAuthServer{
		systemStub: systemStub,
	}, nil
}

func (s *OAuthServer) getProviderAuthUser(ctx context.Context, namespace string, providerName string, code string, codeVerifier string, redirectURL string) (*provider.AuthUser, error) {
	oauthProvider, err := provider.NewProviderByName(providerName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid provider: cant find provider implementation: %v", providerName)
	}

	configCollection := configCollectionByNamespace(s.systemStub, namespace)
	var providerConfig authProviderConfigInMongo
	if err := configCollection.FindOne(ctx, bson.M{"name": providerName}).Decode(&providerConfig); err != nil {
		if err == mongo.ErrNoDocuments {

		} else if err, ok := err.(mongo.WriteException); ok && err.HasErrorLabel("InvalidNamespace") {

		} else {
			return nil, status.Errorf(codes.Internal, "cant load provider configuration: %v", err)
		}
	}

	if !providerConfig.Enabled {
		return nil, errOAuthProviderDisabled
	}
	if providerConfig.AuthUrl != "" {
		oauthProvider.SetAuthUrl(providerConfig.AuthUrl)
	}
	if providerConfig.TokenUrl != "" {
		oauthProvider.SetTokenUrl(providerConfig.TokenUrl)
	}
	if providerConfig.UserApiUrl != "" {
		oauthProvider.SetUserApiUrl(providerConfig.UserApiUrl)
	}
	if providerConfig.ClientId != "" {
		oauthProvider.SetClientId(providerConfig.ClientId)
	}
	if len(providerConfig.ClientSecret) != 0 {
		secret, err := decryptProviderClientSecret(ctx, s.systemStub, providerConfig.ClientSecret)
		if err != nil {
			return nil, err
		}
		oauthProvider.SetClientSecret(secret)
	}

	oauthProvider.SetRedirectUrl(redirectURL)

	providerContext, cancelProviderContext := context.WithTimeout(ctx, time.Second*10)
	defer cancelProviderContext()

	oauthProvider.SetContext(providerContext)

	token, err := oauthProvider.FetchToken(code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		log.Printf("error while retrieving token: %v", err)
		return nil, errOAuthErrorWhileRetrievingToken
	}

	authUser, err := oauthProvider.FetchAuthUser(token)
	if err != nil {
		return nil, errOAuthErrorWhileFetchingUserDetails
	}

	return authUser, nil
}

func (s *OAuthServer) Authenticate(ctx context.Context, in *grpc.AuthenticateRequest) (*grpc.AuthenticateResponse, error) {
	providerName, ok := ProviderGRPCTypeToString[in.Provider]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid provider: provider not supported: %v", in.Provider)
	}

	authUser, err := s.getProviderAuthUser(ctx, in.Namespace, providerName, in.Code, in.CodeVerifier, in.RedirectUrl)
	if err != nil {
		if err == errOAuthProviderDisabled {
			return &grpc.AuthenticateResponse{Status: grpc.AuthenticateResponse_PROVIDER_DISABLED}, status.Error(codes.OK, "")
		}
		if err == errOAuthErrorWhileRetrievingToken {
			return &grpc.AuthenticateResponse{Status: grpc.AuthenticateResponse_ERROR_WHILE_RETRIEVING_AUTH_TOKEN}, status.Error(codes.OK, "")
		}
		if err == errOAuthErrorWhileFetchingUserDetails {
			return &grpc.AuthenticateResponse{Status: grpc.AuthenticateResponse_ERROR_WHILE_FETCHING_USER_DETAILS}, status.Error(codes.OK, "")
		}

		return nil, err
	}

	registrationCollection := registrationCollectionByNamespace(s.systemStub, in.Namespace)
	var registration registrationInMongo
	if err := registrationCollection.FindOne(ctx, bson.M{providerName + ".id": authUser.Id}).Decode(&registration); err != nil {
		if err == mongo.ErrNoDocuments {
			return &grpc.AuthenticateResponse{Status: grpc.AuthenticateResponse_UNAUTHENTICATED}, status.Error(codes.OK, "")
		}
		if err, ok := err.(mongo.WriteException); ok && err.HasErrorLabel("InvalidNamespace") {
			return &grpc.AuthenticateResponse{Status: grpc.AuthenticateResponse_UNAUTHENTICATED}, status.Error(codes.OK, "Invalid namespace")
		}
		return nil, status.Errorf(codes.Internal, "cant load registration: %v", err)
	}

	return &grpc.AuthenticateResponse{
		Status:   grpc.AuthenticateResponse_OK,
		Identity: registration.Indetity,
		UserDetails: &grpc.ProviderUserDetails{
			Id:        authUser.Id,
			Name:      authUser.Name,
			Email:     authUser.Email,
			AvatarUrl: authUser.AvatarUrl,
			Username:  authUser.Username,
		},
	}, status.Error(codes.OK, "")
}
func (s *OAuthServer) RegisterProviderForIdentity(ctx context.Context, in *grpc.RegisterProviderForIdentityRequest) (*grpc.RegisterProviderForIdentityResponse, error) {
	providerName, ok := ProviderGRPCTypeToString[in.Provider]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid provider: provider not supported: %v", in.Provider)
	}

	authUser, err := s.getProviderAuthUser(ctx, in.Namespace, providerName, in.Code, in.CodeVerifier, in.RedirectUrl)
	if err != nil {
		if err == errOAuthProviderDisabled {
			return &grpc.RegisterProviderForIdentityResponse{Status: grpc.RegisterProviderForIdentityResponse_PROVIDER_DISABLED}, status.Error(codes.OK, "")
		}
		if err == errOAuthErrorWhileRetrievingToken {
			return &grpc.RegisterProviderForIdentityResponse{Status: grpc.RegisterProviderForIdentityResponse_ERROR_WHILE_RETRIEVING_AUTH_TOKEN}, status.Error(codes.OK, "")
		}
		if err == errOAuthErrorWhileFetchingUserDetails {
			return &grpc.RegisterProviderForIdentityResponse{Status: grpc.RegisterProviderForIdentityResponse_ERROR_WHILE_FETCHING_USER_DETAILS}, status.Error(codes.OK, "")
		}

		return nil, err
	}

	collection := registrationCollectionByNamespace(s.systemStub, in.Namespace)
	findFilter := bson.M{"identity": in.Identity}
	userDetails := userDetailsInMongo{
		ID:        authUser.Id,
		Name:      &authUser.Name,
		Email:     &authUser.Email,
		AvatarUrl: &authUser.AvatarUrl,
		Username:  &authUser.Username,
	}
	update := bson.M{"$set": bson.M{providerName: userDetails}, "$setOnInsert": bson.M{"identity": in.Identity}}
	options := options.Update().SetUpsert(true)

	if _, err := collection.UpdateOne(ctx, findFilter, update, options); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &grpc.RegisterProviderForIdentityResponse{Status: grpc.RegisterProviderForIdentityResponse_ALREADY_REGISTERED, UserDetails: userDetails.ToGRPCUserDetails()}, status.Error(codes.OK, "")
		}

		return nil, status.Errorf(codes.Internal, "cant update registration in database: cant udpate user details for provider: %v", err)
	}

	return &grpc.RegisterProviderForIdentityResponse{Status: grpc.RegisterProviderForIdentityResponse_OK, UserDetails: userDetails.ToGRPCUserDetails()}, status.Error(codes.OK, "")
}
func (s *OAuthServer) ForgetIdentityProvider(ctx context.Context, in *grpc.ForgetIdentityProviderRequest) (*grpc.ForgetIdentityProviderResponse, error) {
	providerName, ok := ProviderGRPCTypeToString[in.Provider]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid provider: provider not supported: %v", in.Provider)
	}

	registrationCollection := registrationCollectionByNamespace(s.systemStub, in.Namespace)
	if _, err := registrationCollection.UpdateOne(ctx, bson.M{"identity": in.Identity}, bson.M{"$unset": bson.M{providerName + ".id": ""}}); err != nil {
		if err, ok := err.(mongo.WriteException); ok && err.HasErrorLabel("InvalidNamespace") {
			return &grpc.ForgetIdentityProviderResponse{}, status.Error(codes.OK, "Invalid namespace")
		}
		return nil, status.Errorf(codes.Internal, "cant forget provider: cant update registration: %v", err)
	}

	return &grpc.ForgetIdentityProviderResponse{}, status.Error(codes.OK, "")
}
func (s *OAuthServer) GetRegisteredIdentityProviders(ctx context.Context, in *grpc.GetRegisteredIdentityProvidersRequest) (*grpc.GetRegisteredIdentityProvidersResponse, error) {
	collection := registrationCollectionByNamespace(s.systemStub, in.Namespace)
	var registration registrationInMongo
	if err := collection.FindOne(ctx, bson.M{"identity": in.Identity}).Decode(&registration); err != nil {
		if err == mongo.ErrNoDocuments {
			return &grpc.GetRegisteredIdentityProvidersResponse{Providers: []*grpc.GetRegisteredIdentityProvidersResponse_RegisteredProvider{}}, status.Error(codes.OK, "")
		}
		if err, ok := err.(mongo.WriteException); ok && err.HasErrorLabel("InvalidNamespace") {
			return &grpc.GetRegisteredIdentityProvidersResponse{Providers: []*grpc.GetRegisteredIdentityProvidersResponse_RegisteredProvider{}}, status.Error(codes.OK, "Invalid namespace")
		}
		return nil, status.Errorf(codes.Internal, "cant load registration from database: %v", err)
	}

	providers := []*grpc.GetRegisteredIdentityProvidersResponse_RegisteredProvider{}
	addProvider := func(providerType grpc.ProviderType, details *userDetailsInMongo) {
		if details == nil {
			return
		}
		providers = append(providers, &grpc.GetRegisteredIdentityProvidersResponse_RegisteredProvider{
			Provider:    providerType,
			UserDetails: details.ToGRPCUserDetails(),
		})
	}

	addProvider(grpc.ProviderType_GITHUB, registration.GitHub)
	addProvider(grpc.ProviderType_GOOGLE, registration.Google)
	addProvider(grpc.ProviderType_GITLAB, registration.GitLab)
	addProvider(grpc.ProviderType_APPLE, registration.Apple)
	addProvider(grpc.ProviderType_MICROSOFT, registration.Microsoft)
	addProvider(grpc.ProviderType_DISCORD, registration.Discord)
	addProvider(grpc.ProviderType_FACEBOOK, registration.Facebook)
	addProvider(grpc.ProviderType_INSTAGRAM, registration.Instagram)
	addProvider(grpc.ProviderType_OIDC, registration.OIDC)
	addProvider(grpc.ProviderType_OIDC2, registration.OIDC2)
	addProvider(grpc.ProviderType_OIDC3, registration.OIDC3)
	addProvider(grpc.ProviderType_TWITTER, registration.Twitter)

	return &grpc.GetRegisteredIdentityProvidersResponse{Providers: providers}, status.Error(codes.OK, "")
}
