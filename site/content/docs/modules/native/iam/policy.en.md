# IAM Policy service

The `native_iam_policy` is a service for managing user-defined authorization policies. Those policies can be assigned to the `native_identity` identity. The policy specifies what resources identity can access and what actions it can perform on those resources.

## Schema
The schema is defined using `protobuf`. Definitions are provided by the `native` module.

| Property  | Type     | Description                                                                                                                                                                                                                   |
| --------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| uuid      | string   | Unique identifier. Service will automatically assign it to the policy.                                                                                                                                                        |
| name      | string   | User-defined name. It can be whatever you want. It is just for you to distinguish them                                                                                                                                        |
| namespace | string   | `native_namespace` namespace where policy is defined. Policies are valid only inside theirs namespaces. If the namespace is empty, the policy is "global." Global policies are working and granting access to all namespaces. |
| resources | string[] | List of the resources that this policy can access.                                                                                                                                                                            |
| actions   | string[] | List of activities that can be performed on specified resources by this policy.                                                                                                                                               |

!!! tip
    Resources and actions can be defined with wildcards to grant access to the group of targets. Use the `*` symbol at the end of the resource or action to make it a wildcard.

## Examples
Those are several examples of defined policies in JSON format:
??? example "Root access to entire system"
    ```json
        {
            "name": "Root access",
            "uuid": "542c2b97bac0595474108123",
            "namespace": "",
            "resources": ["*"],
            "actions": ["*"]
        }
    ```

??? example "Strict access to specific resource"
    ```json
        {
            "name": "Read access to my service",
            "uuid": "542c2b97bac0595474108124",
            "namespace": "somenamespace",
            "resources": ["mycompany.myproject.someservice.resource1"],
            "actions": ["mycompany.myproject.someservice.get", "mycompany.myproject.someservice.list"]
        }
    ```
??? example "Wildcard access to the group of the resources"
    ```json
        {
            "name": "Full access for my service",
            "uuid": "542c2b97bac0595474108125",
            "namespace": "somenamespace",
            "resources": ["mycompany.myproject.someservice.*"],
            "actions": ["mycompany.myproject.someservice.*"]
        }
    ```
!!! info
    There is no strict structure for the resource and action names, but it's recommended to use `<companyname>.<project/module>.<service>.<......>` prefix. This way, all the resources and actions will be unique and will not have unwanted collisions.

## API
Communication with the service is possible using the gRPC interface. Definitions of the interface (proto file) are provided by the `native` module.

??? example "rpc Create(CreatePolicyRequest) returns (CreatePolicyResponse);"
    This endpoint allows you to create new policy.

    === "Request"
        | Parameter name | Type            | Description                                    |
        | -------------- | --------------- | ---------------------------------------------- |
        | namespace      | string          | Namespace where policy will be located         |
        | name           | string          | User-defined name                              |
        | resources      | array of string | Resources that can be accessed by this policy  |
        | actions        | array of string | Actions that can be performed on the resources |
    === "OK"
        The policy was successfully created. Returns newly created policy with all the data and assigned UUID.
        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for policy creation.
    === "FAILED_PRECONDITION"
        Namespace doesn't exist

??? example "rpc Get(GetPolicyRequest) returns (GetPolicyResponse);"
    Gets policy using its UUID and namespace.

    === "Request"
        | Parameter name | Type   | Description                                                                      |
        | -------------- | ------ | -------------------------------------------------------------------------------- |
        | namespace      | string | Namespace where policy located                                                   |
        | uuid           | string | Unique identifier of policy inside namespace                                     |
        | useCache       | bool   | Use cache for this request or not. Cache may be invalid under rare circumstances |
    === "OK"
        Returns policy data. See [Schema](#Schema).
    === "NOT_FOUND"
        Namespace doesn't exist, or there is no policy with a specified UUID in this namespace
    === "INVALID_ARGUMENT"
        UUID has a bad format

??? example "rpc Exist(ExistPolicyRequest) returns (ExistPolicyResponse);"
    Checks if policy with specified UUID exist in namespace

    === "Request"
        | Parameter name | Type   | Description                                                                      |
        | -------------- | ------ | -------------------------------------------------------------------------------- |
        | namespace      | string | Namespace where policy located                                                   |
        | uuid           | string | Unique identifier of policy inside namespace                                     |
        | useCache       | bool   | Use cache for this request or not. Cache may be invalid under rare circumstances |
    === "OK"
        | Parameter name | Type | Description         |
        | -------------- | ---- | ------------------- |
        | exist          | bool | Policy exist or not |

??? example "rpc Update(UpdatePolicyRequest) returns (UpdatePolicyResponse);"
    Updates policy information

    === "Request"
        | Parameter name | Type            | Description                                    |
        | -------------- | --------------- | ---------------------------------------------- |
        | namespace      | string          | Namespace where policy located                 |
        | uuid           | string          | Unique identifier of the policy                |
        | name           | string          | User-defined name                              |
        | resources      | array of string | Resources that can be accessed by this policy  |
        | actions        | array of string | Actions that can be performed on the resources |
    === "OK"
        Information was successfully updated. Returns updated information as a response to this request. See [Schema](#Schema). The service also cleared all the cache related to this policy.
        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for policy updates.
    === "NOT_FOUND"
        Namespace doesn't exist, or there is no policy with a specified UUID in this namespace
    === "INVALID_ARGUMENT"
        UUID has a bad format

??? example "rpc Delete(DeletePolicyRequest) returns (DeletePolicyResponse);"
    Deletes policy

    === "Request"
        | Parameter name | Type   | Description                     |
        | -------------- | ------ | ------------------------------- |
        | namespace      | string | Namespace where policy located  |
        | uuid           | string | Unique identifier of the policy |
    === "OK"
        The policy was successfully deleted. The service also cleared all the cache related to this policy.
        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for policy deletion.

??? example "rpc List(ListPoliciesRequest) returns (stream ListPoliciesResponse);"
    Streams list of policies in namespace. See [Schema](#Schema).

    === "Request"
        | Parameter name | Type         | Description                                                 |
        | -------------- | ------------ | ----------------------------------------------------------- |
        | namespace      | string       | Namespace where to search policies                          |
        | skip           | unsigned int | How many entries to skip before returning actual policies   |
        | limit          | unsigned int | Maximum number of policies to return. 0 to ignore the limit |

## Events
!!! bug
    This feature is ***NOT IMPLEMENTED***

| system_amqp exchange     | routing key | scheme                           | conditions         |
| ------------------------ | ----------- | -------------------------------- | ------------------ |
| native_iam_policy_events | created     | Created policy (protobuf)        | Policy was created |
| native_iam_policy_events | updated     | New version of policy (protobuf) | Policy was updated |
| native_iam_policy_events | deteled     | ?                                | Policy was deleted |

## Configuration
This service is controlled by environment variables.

| Env                                | default                                | description                                                                                                        |
| ---------------------------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| SYSTEM_DB_URL                      | mongodb://root:example@system_db/admin | [Mongo DB URL](https://www.mongodb.com/docs/manual/reference/connection-string/#standard-connection-string-format) |
| SYSTEM_CACHE_URL                   | redis://system_cache                   | System_cache redis connection URL                                                                                  |
| SYSTEM_TELEMETRY_EXPORTER_ENDPOINT | system_telemetry:55680                 | [OTEL connector](https://opentelemetry.io/docs/collector/) endpoint                                                |
| NATIVE_NAMESPACE_URL               | native_namespace:80                    | Native_namespace server gRPC URL                                                                                   |