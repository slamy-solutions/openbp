package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTScope struct {
	Namespace string   `json:"namespace"`
	Resources []string `json:"resources"`
	Actions   []string `json:"actions"`
}

type JWTData struct {
	UUID      string `json:"uuid"`
	Namespace string `json:"namespace"`
	Identity  string `json:"identity"`

	Scopes []JWTScope `json:"scopes"`

	Refresh bool `json:"refresh"`

	jwt.StandardClaims
}

var jwtkey = []byte("my_secret_key")

var ErrInvalidToken = errors.New("Invalid token")
var ErrTokenExpired = errors.New("Token expired")

func NewJWTData(uuid string, namespace string, identity string, scopes []JWTScope, refresh bool, expiresAt time.Time) *JWTData {
	return &JWTData{
		UUID:      uuid,
		Namespace: namespace,
		Identity:  identity,
		Scopes:    scopes,
		Refresh:   refresh,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.UTC().Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
			Issuer:    "OBP native_iam_token",
		},
	}
}

func JWTDataFromString(input string) (*JWTData, error) {
	// TODO: use Vault

	data := &JWTData{}
	_, err := jwt.ParseWithClaims(input, data, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		if err, ok := err.(jwt.ValidationError); ok {
			if err.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
			return nil, ErrInvalidToken
		}
	}

	return data, nil
}

func (t *JWTData) ToSignedString() (string, error) {
	// TODO: use Vault

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	return token.SignedString(jwtkey)
}
