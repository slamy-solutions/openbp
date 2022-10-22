# IAM Token service

The `native_iam_token` is a service for managing authorization tokens.

## Schema
The schema is defined using `protobuf`. Definitions are provided by the `native` module.

`Token` schema:

| Property         | Type     | Description                                                                                                                                                                                                      |
| ---------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| namespace        | string   | [`native_namespace`](../namespace.en.md) namespace where token and identity are located. Epmty for global token (without namespace)                                                                              |
| uuid             | string   | Token unique identifier inside namespace. It is assigned automatically by service                                                                                                                                |
| identity         | string   | [`native_iam_identity`](./identity.en.md) identity unique identifier inside namespace                                                                                                                            |
| disabled         | string   | Identifies if token was manually disabled. Disabled token always fails on authorization and can not be re-enabled                                                                                                |
| expiresAt        | DateTime | Date and time after with token will not be valid and will fail on Refresh and Authorize attempts                                                                                                                 |
| scopes           | Scope[]  | List of token scopes. Describes what actions can token perform on what resources                                                                                                                                 |
| createdAt        | DateTime | Date and time when token was created                                                                                                                                                                             |
| creationMetadata | string   | Arbitrary metadata added on token creation. For example MAC/IP/information of the actor/application/browser/machine that created this token. The exact format of metadata is not defined, but JSON is suggested. |

`Scope` schema:

| Property  | Type     | Description                                                                                                                                         |
| --------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| namespace | string   | [`native_namespace`](../namespace.en.md) namespace to which scope is bounded. Empty string if the scope is not bounded to any namespace (is global) |
| resources | string[] | Resources that can be accessed using this token                                                                                                     |
| actions   | string[] | Actions that can be performed on accessible resources                                                                                               |

??? example
    This example shows a token created by admin user. Using this token, you can access all the resources and perform all possible actions on them.

    ```json
        {
            "namespace": "",
            "uuid": "542c2b97bac0595474108125",
            "identity": "734c2b97bac0595474108526",
            "disabled": false,
            "expiresAt": "2022-10-20T20:20:39.945Z",
            "scopes": [
                {
                    "namespace": "",
                    "resources": ["*"],
                    "actions": ["*"]
                }
            ],
            "createAt": "2022-10-19T20:20:39.945Z",
            "creationMetadata": {"ip": "32.43.12.123", "mac": "2C:54:91:88:C2:E4", "user-agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4942.12 Safari/536.45"}
        }
    ```

## API
Communication with the service is possible using the gRPC interface. Definitions of the interface (proto file) are provided by the `native` module.


??? example "rpc Create(CreateRequest) returns (CreateResponse);"
    Create new token

    === "Request"
        | Parameter | Type    | Description                                                                                                                                                      |
        | --------- | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
        | namespace | string  | [`native_namespace`](../namespace.en.md) namespace where the token will be located and where is related [`native_iam_identity`](identity.en.md) identity located |
        | identity  | string  | Unique identifier of the identity                                                                                                                                |
        | scopes    | Scope[] | Scopes that will be applied to the token                                                                                                                         |
        | metadata  | string  | Actions that can be performed on the resources                                                                                                                   |
    === "OK"
        The token was successfully created.

        | Property     | Type   | Description                                  |
        | ------------ | ------ | -------------------------------------------- |
        | token        | string | Actual access token formated to the string.  |
        | refreshToken | string | Refreshtoken is used to refresh access token |
        | tokenData    | Token  | Token data                                   |

        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for token creation.
    === "FAILED_PRECONDITION"
        [`native_namespace`](../namespace.en.md) namespace doesn't exist

??? example "rpc Get(GetRequest) returns (GetResponse);"
    Get token data using its unique identifier

    === "Request"
        | Parameter | Type   | Description                                                                                                                                         |
        | --------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------- |
        | namespace | string | [`native_namespace`](../namespace.en.md) namespace of the token                                                                                     |
        | uuid      | string | Unique identifier of the token                                                                                                                      |
        | useCache  | bool   | Use cache for this request or not. Cache may be invalid. The invalid cache automatically deletes after short period of time (30 seconds by default) |
    === "OK"
        Returns token data. See [Schema](#Schema).
    === "NOT_FOUND"
        The token [`native_namespace`](../namespace.en.md) namespace doesn't exist, or there is no token with a specified UUID inside the namespace.
    === "INVALID_ARGUMENT"
        UUID has a bad format

??? example "rpc RawGet(RawGetRequest) returns (RawGetResponse);"
    Get token data using token in string format

    === "Request"
        | Parameter | Type   | Description                                                                                                                                         |
        | --------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------- |
        | token     | string | Refresh or access token                                                                                                                             |
        | useCache  | bool   | Use cache for this request or not. Cache may be invalid. The invalid cache automatically deletes after short period of time (30 seconds by default) |
    === "OK"
        Returns token data. See [Schema](#Schema).
    === "NOT_FOUND"
        The token is valid, but it wasn't founded. Maybe [`native_namespace`](../namespace.en.md) namespace or token was deleted.
    === "INVALID_ARGUMENT"
        The token has a bad format

??? example "rpc Delete(DeleteRequest) returns (DeleteResponse);"
    Delete token. If the token doesn't exist - it does nothing.

    === "Request"
        | Parameter | Type   | Description                                                     |
        | --------- | ------ | --------------------------------------------------------------- |
        | namespace | string | [`native_namespace`](../namespace.en.md) namespace of the token |
        | uuid      | string | Unique identifier of the token                                  |
    === "OK"
        The token was deleted. All the cache for this token was cleared.
        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for token deletion.
    === "INVALID_ARGUMENT"
        The token has a bad format


??? example "rpc Disable(DisableRequest) returns (DisableResponse);"
    Disable token using its unique identifier. The disabled token can not be used for authorization and authentication. The token cannot be re-enabled.

    === "Request"
        | Parameter | Type   | Description                                                     |
        | --------- | ------ | --------------------------------------------------------------- |
        | namespace | string | [`native_namespace`](../namespace.en.md) namespace of the token |
        | uuid      | string | Unique identifier of the token                                  |
    === "OK"
        The token was disabled
        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for token disabling.
    === "NOT_FOUND"
        The token wasn't founded. Maybe [`native_namespace`](../namespace.en.md) namespace or token was deleted.
    === "INVALID_ARGUMENT"
        The token has a bad format

??? example "rpc Validate(ValidateRequest) returns (ValidateResponse);"
    Check if token:
    
    1. has valid format
    2. was signed by this service
    3. not expired
    4. exists (there is information about this token in service)
    5. active (wasn't manually disabled)

    === "Request"
        | Parameter | Type   | Description                                                                                                                                                                      |
        | --------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
        | token     | string | Token to validate                                                                                                                                                                |
        | useCache  | bool   | Use cache for faster validation or not. Cache has a very low chance to not be valid. If cache is not valid it will be deleted after short period of time (30 seconds by default) |
    === "OK"
        As the response you will receive:

        | Property  | Type   | Description                                                                                |
        | --------- | ------ | ------------------------------------------------------------------------------------------ |
        | status    | Status | Validation status                                                                          |
        | tokenData | Token  | Token data. See [Schema](#Schema) for token. It will only be returned if the status is OK. |

        Where `Status` is:

        | Status    | Description                                              |
        | --------- | -------------------------------------------------------- |
        | OK        | Everything is ok. The token is valid.                    |
        | INVALID   | The token has a bad format or invalid signature          |
        | NOT_FOUND | The token wasn't founded. Most probably, it was deleted. |
        | DISABLED  | The token was manually disabled.                         |
        | EXPIRED   | The token expired.                                       |
    
??? example "rpc Refresh(RefreshRequest) returns (RefreshResponse);"
    Validate refresh token and create new token based on it.

    === "Request"
        | Parameter    | Type   | Description                                               |
        | ------------ | ------ | --------------------------------------------------------- |
        | refreshToken | string | Refresh token, based on which new token will be generated |
    
    === "OK"
        As the response you will receive:

        | Property  | Type   | Description                                                                                |
        | --------- | ------ | ------------------------------------------------------------------------------------------ |
        | status    | Status | Validation and refresh status                                                              |
        | token     | string | New token. It will only be returned if the status is OK.                                   |
        | tokenData | Token  | Token data. See [Schema](#Schema) for token. It will only be returned if the status is OK. |

        Where `Status` is:

        | Status            | Description                                                               |
        | ----------------- | ------------------------------------------------------------------------- |
        | OK                | Everything is ok.                                                         |
        | INVALID           | The token has a bad format or invalid signature                           |
        | NOT_FOUND         | The token wasn't founded. Most probably, it was deleted.                  |
        | DISABLED          | The token was manually disabled.                                          |
        | EXPIRED           | The token expired.                                                        |
        | NOT_REFRESH_TOKEN | This token was founded, and it is valid, but this is not a refresh token. |

??? example "rpc GetTokensForIdentity(GetTokensForIdentityRequest) returns (stream GetTokensForIdentityResponse);"
    Get tokens for specified [`native_iam_identity`](identity.en.md) identity

    === "Request"
        | Parameter    | Type         | Description                                                                                                             |
        | ------------ | ------------ | ----------------------------------------------------------------------------------------------------------------------- |
        | namespace    | string       | [`native_namespace`](../namespace.en.md) namespace of the tokens and [`native_iam_identity`](identity.en.md) identities |
        | identity     | string       | Unique identifier of the [`native_iam_identity`](identity.en.md) identity                                               |
        | activeFilter | ActiveFilter | How to filter on "active" property of the token. See schema below.                                                      |
        | skip         | uint32       | Skip of results before returning actual tokens. Set to 0 in order not to skip                                           |
        | limit        | uint32       | Limit of the returned results. Set to 0 in order to remove the limit and return all possible results up to the end.     |

        Where `ActiveFilter` is:

        | ActiveFilter    | Description                                           |
        | --------------- | ----------------------------------------------------- |
        | ALL             | Don't filter. Get all tokens.                         |
        | ONLY_ACTIVE     | Only get tokens that weren't disabled and not expired |
        | ONLY_NOT_ACTIVE | Only get tokens that are disabled or expired          |
    
    === "OK"
        The service will stream list of tokens. See [Schema](#Schema) for token data to be returned.
        
        Tokens are ordered using their creation time. The service will return newly created tokens first.

## Events
!!! bug
    This feature is ***NOT IMPLEMENTED***

| system_amqp exchange     | routing key | scheme                           | conditions         |
| ------------------------ | ----------- | -------------------------------- | ------------------ |
| native_iam_token_events | created     | Created token (protobuf)        | Token was created |
| native_iam_token_events | updated     | ? | Token was updated |
| native_iam_token_events | deteled     | ?                                | Token was deleted |

## Configuration
This service is controlled by environment variables.

| Env                                | default                                | description                                                                                                        |
| ---------------------------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| SYSTEM_DB_URL                      | mongodb://root:example@system_db/admin | [Mongo DB URL](https://www.mongodb.com/docs/manual/reference/connection-string/#standard-connection-string-format) |
| SYSTEM_CACHE_URL                   | redis://system_cache                   | System_cache redis connection URL                                                                                  |
| SYSTEM_TELEMETRY_EXPORTER_ENDPOINT | system_telemetry:55680                 | [OTEL connector](https://opentelemetry.io/docs/collector/) endpoint                                                |
| NATIVE_NAMESPACE_URL               | native_namespace:80                    | Native_namespace server gRPC URL                                                                                   |