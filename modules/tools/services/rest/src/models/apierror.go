package models

type APIErrorCode = string

const (
	ErrorAuthPasswordIdentityNotActive    = "AUTH_PASSWORD_IDENTITY_NOT_ACTIVE"
	ErrorAuthPasswordCredentialsInvalid   = "AUTH_PASSWORD_CREDENTIALS_INVALID"
	ErrorAuthPasswordNotEnoughtPrivileges = "AUTH_PASSWORD_NOT_ENOUGHT_PRIVILEGES"
	ErrorAuthPasswordUnauthorizedUnknown  = "AUTH_PASSWORD_UNAUTHORIZED_UNKNOWN"

	ErrorAuthTokenRefreshTokenInvalid            = "AUTH_TOKEN_REFRESH_TOKEN_INVALID"
	ErrorAuthTokenRefreshTokenNotFound           = "AUTH_TOKEN_REFRESH_TOKEN_NOT_FOUND"
	ErrorAuthTokenRefreshTokenDisabled           = "AUTH_TOKEN_REFRESH_TOKEN_DISABLED"
	ErrorAuthTokenRefreshTokenExpired            = "AUTH_TOKEN_REFRESH_TOKEN_EXPIRED"
	ErrorAuthTokenRefreshTokenIsNotRefreshToken  = "AUTH_TOKEN_REFRESH_TOKEN_IS_NOT_REFRESH_TOKEN"
	ErrorAuthTokenRefreshIdentityNotFound        = "AUTH_TOKEN_REFRESH_IDENTITY_NOT_FOUND"
	ErrorAuthTokenRefreshIdentityNotActive       = "AUTH_TOKEN_REFRESH_IDENTITY_NOT_ACTIVE"
	ErrorAuthTokenRefreshIdentityUnauthenticated = "AUTH_TOKEN_REFRESH_IDENTITY_UNAUTHENTICATED"
	ErrorAuthTokenRefreshUnknown                 = "AUTH_TOKEN_REFRESH_UNKNOWN"
)

var errorMessage = map[APIErrorCode]string{
	ErrorAuthPasswordIdentityNotActive:    "Identity not active. Maybe it was manually disabled.",
	ErrorAuthPasswordCredentialsInvalid:   "Login or password is invalid.",
	ErrorAuthPasswordNotEnoughtPrivileges: "Not enought privileges to create token with specified scopes.",
	ErrorAuthPasswordUnauthorizedUnknown:  "Not authorized for unknown reasons.",

	ErrorAuthTokenRefreshTokenInvalid:            "Refresh token is invalid. It may have bad format or signature.",
	ErrorAuthTokenRefreshTokenNotFound:           "Token not found. Most probably it was deleted from the system.",
	ErrorAuthTokenRefreshTokenDisabled:           "Auth token was manually disabled.",
	ErrorAuthTokenRefreshTokenExpired:            "Refresh token expired and cant be used to create new access tokens.",
	ErrorAuthTokenRefreshTokenIsNotRefreshToken:  "This token is not a refresh token. Probably access token was sended onstead of refresh one.",
	ErrorAuthTokenRefreshIdentityNotFound:        "Token identity was deleted.",
	ErrorAuthTokenRefreshIdentityNotActive:       "Token identity is not active. Probably it was manually disabled.",
	ErrorAuthTokenRefreshIdentityUnauthenticated: "Identity privilages changed and now you cant create new token with same scopes using this refresh token.",
	ErrorAuthTokenRefreshUnknown:                 "Cant refresh token. Unknown or hidden user error.",
}

type APIError struct {
	Code    APIErrorCode `json:"code"`
	Message string       `json:"message"`
}

func NewAPIError(code APIErrorCode) *APIError {
	return &APIError{
		Code:    code,
		Message: errorMessage[code],
	}
}

func ErrorMessageFromCode(code APIErrorCode) string {
	return errorMessage[code]
}
