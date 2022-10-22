# Namespace service
The `native_namespace` is a service for managing namespaces in the entire platform ecosystem. Namespace provides logical separation of the data for modules and services. They can be used, for example, to distinguish customers in the application.
After the creation of the namespace, the `native_namespace` service creates a new `namespace_<name>` database so that every new namespace data can be placed in a separate database. On `namespace` deletion, the service will delete its database, and all the data will be lost.

## API
Communication with the system is possible using the gRPC interface. Definitions of the interface (proto file) are provided by the `native` module.

??? example "rpc Ensure(EnsureNamespaceRequest) returns (EnsureNamespaceResponse) {};"
    This endpoint allows you to create new namespace if it doesn't exist. If the namespace already exists, this endpoint will do nothing.

    === "Request"
        | Parameter name | Type   | Description                  | Required | Format                                         |
        | -------------- | ------ | ---------------------------- | -------- | ---------------------------------------------- |
        | name           | string | Unique name of the namespace | True     | Will be validated using regex `^[A-Za-z0-9]+$` |
    === "OK"
        Namespace was successfully created. Service created new `<GLOBAL_PREFIX>namespace_<name>` database.
        !!! info
            This response will raise the event on the amqp. Check the [Events](#events) section and the specific event for the namespace creation.
    === "IVALID_ARGUMENT"
        The name has a bad format. Check it against `^[A-Za-z0-9]+$` regex and try again

??? example "rpc Delete(DeleteNamespaceRequest) returns (DeleteNamespaceResponse) {};"
    This endpoint will try to delete namespace and corresponding database.
    !!! danger
        This is very danger operation. Service will delete the entire database and all data of all other services.
    === "Request"
        | Parameter name | Type   | Description                                | Required | Format |
        | -------------- | ------ | ------------------------------------------ | -------- | ------ |
        | name           | string | Unique name of the namespace to be deleted | True     | -      |
    === "OK"
        Namespace was successfully deleted or wasn't exist. Service has removed the entry from the `<GLOBAL_PREFIX>global` database and deleted `<GLOBAL_PREFIX>namespace_<name>` database.
        !!! info
            This response will raise an event on the amqp. Check the [Events](#events) section and specific events raised on namespace deletion.

??? example "rpc Get(GetNamespaceRequest) returns (GetNamespaceResponse) {};"
    Returns namespace information by its name
    === "Request"
        | Parameter name | Type   | Description                                                            | Required | Format |
        | -------------- | ------ | ---------------------------------------------------------------------- | -------- | ------ |
        | name           | string | Unique name of the namespace to get                                    | True     | -      |
        | useCache       | bool   | Tell services if it can use cache in order to optimize response speed. | True     | -      |
        !!! tip
            Using cache greatly reduces response times. The cache may not be valid under rare conditions. Cache may not be valid only when there are a lot of read/write operations on namespace information simultaneously. In the real system, `namespace` probably will never change, so it's excellent practice to use cache wherever you can. The invalid cache deletes after some time (60 seconds by default). The system tries to keep the cache up to date, but it cants ensure its validity all the time.
    === "OK"
        Namespace was founded, and its data was successfully returned.
        !!! info
            More information may be returned in future versions of the service
        | Response value | Type   | Description           | Format                 |
        | -------------- | ------ | --------------------- | ---------------------- |
        | name           | string | Name of the namespace | Regex `^[A-Za-z0-9]+$` |
    === "NOT_FOUND"
        Service can't find data about this namespaces. Probably, it doesn't exist.

??? example "rpc GetAll(GetAllNamespacesRequest) returns (stream GetAllNamespacesResponse) {};"
    Gets stream of all namespaces
    === "Request"
        Request doesn't has parameters. It just returns all the namespaces.
    === "OK"
        Successfully geted stream of the namespaces.
        !!! info
            More information may be returned in the future versions on the service
        | Response value | Type   | Description           | Format                 |
        | -------------- | ------ | --------------------- | ---------------------- |
        | name           | string | Name of the namespace | Regex `^[A-Za-z0-9]+$` |

??? example "rpc Exists(IsNamespaceExistRequest) returns (IsNamespaceExistResponse) {};"
    Checks if namespace with provided name exists
    === "Request"
        | Parameter name | Type   | Description                                                            | Required |
        | -------------- | ------ | ---------------------------------------------------------------------- | -------- |
        | name           | string | Unique name of the namespace to check                                  | True     |
        | useCache       | bool   | Tell services if it can use cache in order to optimize response speed. | True     |
        !!! tip
            Using cache greatly reduces response times. Assuming you will delete the namespace only once, it is an outstanding practice to use cache whenever possible. The invalid cache will be deleted after a short time (default is 60 seconds).
    === "OK"
        Successfully checked if namespace exists.
        | Response value | Type | Description            |
        | -------------- | ---- | ---------------------- |
        | exist          | bool | Namespace exist or not |

## Events
!!! bug
    This feature is ***NOT IMPLEMENTED***

| system_rabbitmq exchange               | routing key | scheme               | conditions            |
| -------------------------------------- | ----------- | -------------------- | --------------------- |
| <GLOBAL_PREFIX>native_namespace_events | created     | Namespace (protobuf) | Namespace was created |
| <GLOBAL_PREFIX>native_namespace_events | deteled     | Namespace (protobuf) | Namespace was deleted |

## Configuration
This service is controlled by environment variables.

| Env                                | default                                | description                                                                                                        |
| ---------------------------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| SYSTEM_DB_URL                      | mongodb://root:example@system_db/admin | [Mongo DB URL](https://www.mongodb.com/docs/manual/reference/connection-string/#standard-connection-string-format) |
| SYSTEM_CACHE_URL                   | redis://system_cache                   | System_cache redis connection URL                                                                                  |
| SYSTEM_TELEMETRY_EXPORTER_ENDPOINT | system_telemetry:55680                 | [OTEL connector](https://opentelemetry.io/docs/collector/) endpoint                                                |