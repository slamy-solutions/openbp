# IAM Identity service

The `native_iam_identity` is service for managing authentication identities. Its mainly used by other actors (users/applications/servers) as a way to outsource identification capabilities. For example, in OpenBP, users don't have a login/password/security key. They have one-to-one binding with `native_iam_identity` identity instead.

Identity implements authentification capabilities using `native_iam_authentication` methods (like `native_iam_authentication_password` to authenticate using a secret password).

Authorization implemented using `native_iam_policy` policies. Policies can be added and removed from the identity. One identity can have many policies assigned.

## Schema
The schema of the identity is defined using `protobuf`. Definitions are provided by the `native` module.

`Identity` schema:

| Property  | Type              | Description                                                                      |
| --------- | ----------------- | -------------------------------------------------------------------------------- |
| namespace | string            | Namespaces of the identity. Can be empty for global identities.                  |
| uuid      | string            | Unique identity identifier inside namespace. Assigned automatically by service   |
| name      | string            | User-defined public name                                                         |
| active    | bool              | If identity is not active, it will not be able to login and perform any actions. |
| policies  | PolicyReference[] | `native_iam_policy` policies assigned to the identity                            |

`PolicyReference` schema:

| Property  | Type   | Description                                  |
| --------- | ------ | -------------------------------------------- |
| namespace | string | Namespace where policy located               |
| uuid      | string | Unique identifier of policy inside namespace |

??? example
    This is an example of identity that belongs to the admin user. 
    
    Assume you have defined `native_iam_policy` policy with "542c2b97bac0595474108125" uuid that grants access to all resources.
    
    ```json
        {
            "namespace": "",
            "uuid": "142c2b97bac0595474173823",
            "name": "User admin",
            "active": true,
            "policies": [
                {
                    "namespace": "",
                    "uuid": "542c2b97bac0595474108125"
                }
            ]
        }
    ```

## API
Communication with the service is possible using the gRPC interface. Definitions of the interface (proto file) are provided by the `native` module.

??? example "rpc Create(CreateIdentityRequest) returns (CreateIdentityResponse);"
    Creates new identity

    === "Request"
        | Parameter name  | Type   | Description                                                                                        |
        | --------------- | ------ | -------------------------------------------------------------------------------------------------- |
        | namespace       | string | Namespace where identity will be located                                                           |
        | name            | string | Public name for newly created identity. It may not be unique - this is just a human-readable name. |
        | initiallyActive | bool   | Should the identity be active at the start or not                                                  |
    === "OK"
        The identity was successfully created. Returns newly created identity with all the data and assigned UUID.
        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for identity creation.
    === "FAILED_PRECONDITION"
        Namespace doesn't exist

??? example "rpc Get(GetIdentityRequest) returns (GetIdentityResponse);"
    Gets identity by its namespace and UUID

    === "Request"
        | Parameter name | Type   | Description                                                                                                                                                                        |
        | -------------- | ------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
        | namespace      | string | Namespace of the identity                                                                                                                                                          |
        | uuid           | string | Unique identifier of the identity inside namespace                                                                                                                                 |
        | useCache       | bool   | Use cache for this request or not. The cache may not be valid in rare circumstances. The invalid cache automatically deletes after a short period of time (30 seconds by default). |
    === "OK"
        Returns identity. See [Schema](#Schema).
    === "NOT_FOUND"
        The namespace doesn't exist, or there is no identity with a specified UUID inside the namespace.
    === "INVALID_ARGUMENT"
        UUID has a bad format

??? example "rpc Delete(DeleteIdentityRequest) returns (DeleteIdentityResponse);"
    Deletes identity by its namespace and UUID

    === "Request"
        | Parameter name | Type   | Description                                        |
        | -------------- | ------ | -------------------------------------------------- |
        | namespace      | string | Namespace of the identity                          |
        | uuid           | string | Unique identifier of the identity inside namespace |
    === "OK"
        Identity doesn't exist in the system. It was deleted during this operation or earlier. The service also cleared all the cache related to this identity.
        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for identity deletion.


??? example "rpc AddPolicy(AddPolicyRequest) returns (AddPolicyResponse);"
    Adds policy to the identity. If identity already has this policy attached - it does nothing. So you can use this method to ensure identity has the policies you need.

    === "Request"
        | Parameter name    | Type   | Description                                         |
        | ----------------- | ------ | --------------------------------------------------- |
        | identityNamespace | string | Namespace of the identity                           |
        | identityUUID      | string | Unique identifier of the identity inside namespace  |
        | policyNamespace   | string | Namespace of the `native_iam_policy` policy         |
        | policyUUID        | string | Unique identifier of the `native_iam_policy` policy |
    === "OK"
        The policy was successfully assigned to the identity.
        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for identity policies changes.
    === "INVALID_PRECONDITION"
        Policy doesnt exist
    === "NOT_FOUND"
        The identity namespace doesn't exist, or there is no identity with a specified UUID inside the namespace.
    === "INVALID_ARGUMENT"
        Identity UUID has bad format

??? example "rpc RemovePolicy(RemovePolicyRequest) returns (RemovePolicyResponse);"
    Removes policy from the identity. If identity doesn't have an assigned policy - it does nothing.

    === "Request"
        | Parameter name    | Type   | Description                                         |
        | ----------------- | ------ | --------------------------------------------------- |
        | identityNamespace | string | Namespace of the identity                           |
        | identityUUID      | string | Unique identifier of the identity inside namespace  |
        | policyNamespace   | string | Namespace of the `native_iam_policy` policy         |
        | policyUUID        | string | Unique identifier of the `native_iam_policy` policy |
    === "OK"
        The policy was successfully removed from to the identity.
        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for identity policies changes.
    === "NOT_FOUND"
        The identity namespace doesn't exist, or there is no identity with a specified UUID inside the namespace.
    === "INVALID_ARGUMENT"
        Identity UUID has bad format

??? example "rpc SetActive(SetIdentityActiveRequest) returns (SetIdentityActiveResponse);"
    Set identity active or not.

    === "Request"
        | Parameter name | Type   | Description                                        |
        | -------------- | ------ | -------------------------------------------------- |
        | namespace      | string | Namespace of the identity                          |
        | uuid           | string | Unique identifier of the identity inside namespace |
        | active         | bool   | Setactive or not                                   |
    === "OK"
        Active state of the identity was successfully changed.
        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for identity active state changes.
    === "NOT_FOUND"
        The identity namespace doesn't exist, or there is no identity with a specified UUID inside the namespace.
    === "INVALID_ARGUMENT"
        Identity UUID has bad format

## Events
!!! bug
    This feature is ***NOT IMPLEMENTED***

| system_amqp exchange       | routing key | scheme                      | conditions           |
| -------------------------- | ----------- | --------------------------- | -------------------- |
| native_iam_identity_events | created     | Created identity (protobuf) | Identity was created |
| native_iam_identity_events | updated     | ?                           | Identity was updated |
| native_iam_identity_events | deteled     | ?                           | Identity was deleted |

## Configuration
This service is controlled by environment variables.

| Env                                | default                                | description                                                                                                        |
| ---------------------------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| SYSTEM_DB_URL                      | mongodb://root:example@system_db/admin | [Mongo DB URL](https://www.mongodb.com/docs/manual/reference/connection-string/#standard-connection-string-format) |
| SYSTEM_CACHE_URL                   | redis://system_cache                   | System_cache redis connection URL                                                                                  |
| SYSTEM_TELEMETRY_EXPORTER_ENDPOINT | system_telemetry:55680                 | [OTEL connector](https://opentelemetry.io/docs/collector/) endpoint                                                |
| NATIVE_NAMESPACE_URL               | native_namespace:80                    | Native_namespace server gRPC URL                                                                                   |
| NATIVE_IAM_POLICY_URL              | native_iam_policy:80                   | Native_iam_policy server gRPC URL                                                                                  |
