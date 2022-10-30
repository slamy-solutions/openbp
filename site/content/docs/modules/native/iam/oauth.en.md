# OAuth service

The `native_iam_oauth` is a meta-service, that connects most of the `native_iam` services to deliver OAuth capabilities.

## API
Communication with the service is possible using the gRPC interface. Definitions of the interface (proto file) are provided by the `native` module.


??? example "rpc CreateTokenWithPassword(CreateTokenWithPasswordRequest) returns (CreateTokenWithPasswordResponse);"
    Creates new [`native_iam_token`](./token.en.md) token using [`native_iam_authentication_password`](./authentication/password.en.md) password authentication method.

    === "Request"
        | Parameter name | Type                                      | Description                                                                                                                                                                                                                                                         |
        | -------------- | ----------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
        | namespace      | string                                    | [`native_iam_namespace`](../namespace.en.md) namespace where [`native_iam_identity`](./identity.en.md) identity located. May be null for global identity.                                                                                                           |
        | identity       | string                                    | [`native_iam_identity`](./identity.en.md) identity unique identifier                                                                                                                                                                                                |
        | password       | string                                    | [`native_iam_authentication_password`](./authentication/password.en.md) password                                                                                                                                                                                    |
        | metadata       | string                                    | Arbitrary metadata. For example MAC/IP/information of the actor/application/browser/machine that makes this request. The exact format of metadata is not defined, but JSON is suggested. It will be added to the created [`native_iam_token`](./token.en.md) token. |
        | scopes         | [`native_iam_token`](./token.en.md) Scope | Scopes of the created [`native_iam_token`](./token.en.md) token. Empty for creating token with all possible scopes for identity.                                                                                                                                    |
    === "OK"
        | Property     | Type   | Description                                                                                                             |
        | ------------ | ------ | ----------------------------------------------------------------------------------------------------------------------- |
        | status       | Status | Status of the [`native_iam_identity`](./identity.en.md) identity authentication and authorization                       |
        | accessToken  | string | [`native_iam_token`](./token.en.md) token used for authentication and authorization. If status is not OK - empty string |
        | refreshToken | string | [`native_iam_token`](./token.en.md) token used for refreshing accessToken. If status is not OK - empty string           |

        Where status is:

        | Status              | Description                                                                                                                                                                                                                                  |
        | ------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
        | OK                  | Everything is ok. Access and refresh tokens were successfully created                                                                                                                                                                        |
        | CREDENTIALS_INVALID | [`native_iam_identity`](./identity.en.md) identity or [`native_iam_authentication_password`](./authentication/password.en.md) password is not valid. Maybe password authentication is not enabled for identity or the password doesn't match |
        | IDENTITY_NOT_ACTIVE | [`native_iam_identity`](./identity.en.md) identity was manually disabled                                                                                                                                                                     |
        | UNAUTHORIZED        | Not enough privileges to create [`native_iam_token`](./token.en.md) token with specified scopes                                                                                                                                              |

??? example "rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);"
    Creates new OAuth access token using refresh token.

    === "Request"
        | Parameter name | Type   | Description                                       |
        | -------------- | ------ | ------------------------------------------------- |
        | refreshToken   | string | [`native_iam_token`](./token.en.md) refresh token |

    === "OK"
        | Property    | Type   | Description                                                                                                                    |
        | ----------- | ------ | ------------------------------------------------------------------------------------------------------------------------------ |
        | status      | Status | Status of the [`native_iam_token`](./token.en.md) token and [`native_iam_identity`](./identity.en.md) identity reauthorization |
        | accessToken | string | [`native_iam_token`](./token.en.md) token used for authentication and authorization. If status is not OK - empty string        |

        Where status is:

        | Status                     | Description                                                                                                                                                                                                             |
        | -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
        | OK                         | Everything is ok. New [`native_iam_token`](./token.en.md) access token was successfully created                                                                                                                         |
        | TOKEN_INVALID              | Received [`native_iam_token`](./token.en.md) token has bad format or its signature doesnt match                                                                                                                         |
        | TOKEN_NOT_FOUND            | Most probably [`native_iam_token`](./token.en.md) token was deleted after its creation                                                                                                                                  |
        | TOKEN_DISABLED             | [`native_iam_token`](./token.en.md) token was manually disabled                                                                                                                                                         |
        | TOKEN_EXPIRED              | [`native_iam_token`](./token.en.md) token expired                                                                                                                                                                       |
        | TOKEN_IS_NOT_REFRESH_TOKEN | Provided [`native_iam_token`](./token.en.md) token was recognized but most probably it is ordinary access token (not refresh one)                                                                                       |
        | IDENTITY_NOT_FOUND         | [`native_iam_identity`](./identity.en.md) identity wasn't founded. Most probably it was deleted after token creation                                                                                                    |
        | IDENTITY_NOT_ACTIVE        | [`native_iam_identity`](./identity.en.md) identity was manually disabled.                                                                                                                                               |
        | IDENTITY_UNAUTHENTICATED   | Most probably [`native_iam_identity`](./identity.en.md) identity [`native_iam_policy`](./policy.en.md) policies changed and now it's not possible to create [`native_iam_token`](./token.en.md) tokens with same scopes |

??? example "rpc CheckAccess(CheckAccessRequest) returns (CheckAccessResponse);"
    Checks if a token is allowed to perform actions from the specified scopes

    === "Request"
        | Parameter name | Type   | Description                     |
        | -------------- | ------ | ------------------------------- |
        | accessToken    | string | Token to check                  |
        | scopes         | Scope  | Scopes for with to check access |
    
    === "OK"
        | Property | Type   | Description                                                                   |
        | -------- | ------ | ----------------------------------------------------------------------------- |
        | status   | Status | Status of the check                                                           |
        | message  | string | Details of the status, that can be safelly returned and displayed to the user |

        Where status is:

        | Status          | Description                                                                                     |
        | --------------- | ----------------------------------------------------------------------------------------------- |
        | OK              | Everything is ok. The provided token allows to access scopes.                                   |
        | TOKEN_INVALID   | Received [`native_iam_token`](./token.en.md) token has bad format or its signature doesnt match |
        | TOKEN_NOT_FOUND | Most probably [`native_iam_token`](./token.en.md) token was deleted after its creation          |
        | TOKEN_DISABLED  | [`native_iam_token`](./token.en.md) token was manually disabled                                 |
        | TOKEN_EXPIRED   | [`native_iam_token`](./token.en.md) token expired                                               |
        | UNAUTHORIZED    | [`native_iam_token`](./token.en.md) token has not enough privileges to access specified scopes  |


## Configuration
This service is controlled by environment variables.

| Env                                    | default                                | description                                                                                                        |
| -------------------------------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| SYSTEM_DB_URL                          | mongodb://root:example@system_db/admin | [Mongo DB URL](https://www.mongodb.com/docs/manual/reference/connection-string/#standard-connection-string-format) |
| SYSTEM_CACHE_URL                       | redis://system_cache                   | System_cache redis connection URL                                                                                  |
| SYSTEM_TELEMETRY_EXPORTER_ENDPOINT     | system_telemetry:55680                 | [OTEL connector](https://opentelemetry.io/docs/collector/) endpoint                                                |
| NATIVE_IAM_POLICY_URL                  | native_iam_policy:80                   | Native_iam_policy server gRPC URL                                                                                  |
| NATIVE_IAM_TOKEN_URL                   | native_iam_token:80                    | Native_iam_token server gRPC URL                                                                                   |
| NATIVE_IAM_IDENTITY_URL                | native_iam_identity:80                 | Native_iam_identity server gRPC URL                                                                                |
| NATIVE_IAM_AUTHENTICATION_PASSWORD_URL | native_iam_authentication_password:80  | Native_iam_authentication_password server gRPC URL                                                                 |
