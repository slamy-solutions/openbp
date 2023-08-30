package oauth

import (
	grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/oauth2"

	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/authentication/oauth/provider"
)

var ProviderTypes = []string{provider.NameGithub, provider.NameGitlab, provider.NameMicrosoft, provider.NameGoogle, provider.NameDiscord, provider.NameFacebook, provider.NameApple, provider.NameTwitter, provider.NameOIDC, provider.NameOIDC2, provider.NameOIDC3, provider.NameInstagram}
var ProviderStringTypeToGRPC = map[string]grpc.ProviderType{
	provider.NameGithub:    grpc.ProviderType_GITHUB,
	provider.NameGitlab:    grpc.ProviderType_GITLAB,
	provider.NameMicrosoft: grpc.ProviderType_MICROSOFT,
	provider.NameGoogle:    grpc.ProviderType_GOOGLE,
	provider.NameDiscord:   grpc.ProviderType_DISCORD,
	provider.NameFacebook:  grpc.ProviderType_FACEBOOK,
	provider.NameApple:     grpc.ProviderType_APPLE,
	provider.NameTwitter:   grpc.ProviderType_TWITTER,
	provider.NameOIDC:      grpc.ProviderType_OIDC,
	provider.NameOIDC2:     grpc.ProviderType_OIDC2,
	provider.NameOIDC3:     grpc.ProviderType_OIDC3,
	provider.NameInstagram: grpc.ProviderType_INSTAGRAM,
}
var ProviderGRPCTypeToString = map[grpc.ProviderType]string{
	grpc.ProviderType_GITHUB:    provider.NameGithub,
	grpc.ProviderType_GITLAB:    provider.NameGitlab,
	grpc.ProviderType_MICROSOFT: provider.NameMicrosoft,
	grpc.ProviderType_GOOGLE:    provider.NameGoogle,
	grpc.ProviderType_DISCORD:   provider.NameDiscord,
	grpc.ProviderType_FACEBOOK:  provider.NameFacebook,
	grpc.ProviderType_APPLE:     provider.NameApple,
	grpc.ProviderType_TWITTER:   provider.NameTwitter,
	grpc.ProviderType_OIDC:      provider.NameOIDC,
	grpc.ProviderType_OIDC2:     provider.NameOIDC2,
	grpc.ProviderType_OIDC3:     provider.NameOIDC3,
	grpc.ProviderType_INSTAGRAM: provider.NameInstagram,
}

type authProviderConfigInMongo struct {
	Name string `bson:"name"`

	Enabled      bool   `bson:"enabled"`
	ClientId     string `bson:"clientId"`
	ClientSecret []byte `bson:"clientSecret"`
	AuthUrl      string `bson:"authUrl"`
	TokenUrl     string `bson:"tokenUrl"`
	UserApiUrl   string `bson:"userApiUrl"`
}

func (c *authProviderConfigInMongo) ToGRPCProviderConfig(namespace string) *grpc.ProviderConfig {
	p, _ := provider.NewProviderByName(c.Name)
	authURL := c.AuthUrl
	if authURL == "" {
		authURL = p.AuthUrl()
	}
	tokenURL := c.TokenUrl
	if tokenURL == "" {
		tokenURL = p.TokenUrl()
	}
	userApiURL := c.UserApiUrl
	if userApiURL == "" {
		userApiURL = p.UserApiUrl()
	}

	return &grpc.ProviderConfig{
		Namespace:    namespace,
		Enabled:      c.Enabled,
		ClientId:     c.ClientId,
		ClientSecret: "",
		AuthUrl:      authURL,
		TokenUrl:     tokenURL,
		UserApiUrl:   userApiURL,
		Type:         ProviderStringTypeToGRPC[c.Name],
	}
}

type userDetailsInMongo struct {
	ID        string  `bson:"id"`
	Name      *string `bson:"name,omitempty"`
	Username  *string `bson:"username,omitempty"`
	Email     *string `bson:"email,omitempty"`
	AvatarUrl *string `bson:"avatarUrl,omitempty"`
}

func (d *userDetailsInMongo) ToGRPCUserDetails() *grpc.ProviderUserDetails {
	name := ""
	if d.Name != nil {
		name = *d.Name
	}
	username := ""
	if d.Username != nil {
		username = *d.Username
	}
	email := ""
	if d.Email != nil {
		email = *d.Email
	}
	avatarUrl := ""
	if d.AvatarUrl != nil {
		avatarUrl = *d.AvatarUrl
	}

	return &grpc.ProviderUserDetails{
		Id:        d.ID,
		Name:      name,
		Username:  username,
		Email:     email,
		AvatarUrl: avatarUrl,
	}
}

type registrationInMongo struct {
	Indetity string `bson:"identity"`

	GitHub    *userDetailsInMongo `bson:"github,omitempty"`
	GitLab    *userDetailsInMongo `bson:"gitlab,omitempty"`
	Microsoft *userDetailsInMongo `bson:"microsoft,omitempty"`
	Google    *userDetailsInMongo `bson:"google,omitempty"`
	Discord   *userDetailsInMongo `bson:"discord,omitempty"`
	Facebook  *userDetailsInMongo `bson:"facebook,omitempty"`
	Apple     *userDetailsInMongo `bson:"apple,omitempty"`
	Twitter   *userDetailsInMongo `bson:"twitter,omitempty"`
	OIDC      *userDetailsInMongo `bson:"oidc,omitempty"`
	OIDC2     *userDetailsInMongo `bson:"oidc2,omitempty"`
	OIDC3     *userDetailsInMongo `bson:"oidc3,omitempty"`
	Instagram *userDetailsInMongo `bson:"instagram,omitempty"`
}
