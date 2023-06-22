package token

import (
	"time"

	goJWT "github.com/golang-jwt/jwt"
)

type JWTScope struct {
	Namespace            string   `json:"namespace"`
	Resources            []string `json:"resources"`
	Actions              []string `json:"actions"`
	NamespaceIndependent bool     `json:"namespaceIndependent"`
}

type JWTData struct {
	UUID      string `json:"uuid"`
	Namespace string `json:"namespace"`
	Identity  string `json:"identity"`

	Scopes []JWTScope `json:"scopes"`

	Refresh bool `json:"refresh"`

	goJWT.StandardClaims
}

func NewJWTData(uuid string, namespace string, identity string, scopes []JWTScope, refresh bool, expiresAt time.Time) *JWTData {
	return &JWTData{
		UUID:      uuid,
		Namespace: namespace,
		Identity:  identity,
		Scopes:    scopes,
		Refresh:   refresh,
		StandardClaims: goJWT.StandardClaims{
			ExpiresAt: expiresAt.UTC().Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
			Issuer:    "OBP native_iam_token",
		},
	}
}
