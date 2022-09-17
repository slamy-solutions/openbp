# Namespace service
Native namespace is service for managing namespaceses in the entire platform ecosystem. Namespace provides logical separation of the data for modules and services. They can be used for example to distinguishe customers in the application.

After creation of the `namespace`, native_namespace service creates new `<GLOBAL_PREFIX>namespace_<name>` database, so every new namespace data can be placed in separate database. On `namespace` deletion, service will delete its database and all the data will be lost.

## API
Communication with system is possible using gRPC interface. Definitions of the interface (proto file) are provided by native module.
### Ensure
```
    rpc Ensure(EnsureNamespaceRequest) returns (EnsureNamespaceResponse) {};
```
This endpoint allows you to create new namespace if it doesnt exist. If namespace already exists, this endpoint will do nothing.
<!-- tabs:start -->

#### **Request**

| Parameter name | Type   | Description                  | Required | Format                                         |
| -------------- | ------ | ---------------------------- | -------- | ---------------------------------------------- |
| name           | string | Unique name of the namespace | True     | Will be validated using regex `^[A-Za-z0-9]+$` |

#### **OK**
Namespace was successfully created. Service created new `<GLOBAL_PREFIX>namespace_<name>` database.

?> This response will raise event on the amqp. Check the [Events](#events) section and specific event raised on namespace creation.

#### **IVALID_ARGUMENT**
Name has bad format. Check it against `^[A-Za-z0-9]+$` regex and try again

<!-- tabs:end -->

### Delete
```
    rpc Delete(DeleteNamespaceRequest) returns (DeleteNamespaceResponse) {};
```
This endpoint will try to delete namespace and correcponding database.

!> This is very danger operation. Entire database will all data of all other services will be deleted.

<!-- tabs:start -->

#### **Request**

| Parameter name | Type   | Description                                | Required | Format |
| -------------- | ------ | ------------------------------------------ | -------- | ------ |
| name           | string | Unique name of the namespace to be deleted | True     | -      |

#### **OK**
Namespace was successfully deleted or wasnt exist. Service removed entry from the `<GLOBAL_PREFIX>global` datadase and deleted `<GLOBAL_PREFIX>namespace_<name>` database.

?> This response will raise event on the amqp. Check the [Events](#events) section and specific event raised on namespace deletion.

<!-- tabs:end -->

### Get
```
    rpc Get(GetNamespaceRequest) returns (GetNamespaceResponse) {};
```
Returns namespace information by its name

<!-- tabs:start -->

#### **Request**

| Parameter name | Type   | Description                                                            | Required | Format |
| -------------- | ------ | ---------------------------------------------------------------------- | -------- | ------ |
| name           | string | Unique name of the namespace to get                                    | True     | -      |
| useCache       | bool   | Tell services if it can use cache in order to optimize response speed. | True     | -      |

?> Using cache greatly reduses response times. Cache may not be valid under very rare conditions. Cache may not be valid only when there are a lot of read/write operations on namespace information at same time. In the real system, namespace probably will never change so its very good practice to use cache whereever you can. Invalid cache automatically deletes after very short period of time (60 seconds by default). System tries to keep cache up to date, but it cant ensure its validity all the time.

#### **OK**
Namespace was founded and its data was successfully returned.

?> More information may be returned in the future versions of the service

| Response value | Type   | Description           | Format                 |
| -------------- | ------ | --------------------- | ---------------------- |
| name           | string | Name of the namespace | Regex `^[A-Za-z0-9]+$` |

#### **NOT_FOUND**
Service cant find data about this namespaces. Probably, it doesnt exist.

<!-- tabs:end -->

### GetAll
```
    rpc GetAll(GetAllNamespacesRequest) returns (stream GetAllNamespacesResponse) {};
```
Gets stream of all namespaces

<!-- tabs:start -->

#### **Request**

Request doesnt has parameters. It just return all the namespaces.

#### **OK**

Successfully geted stream of the namespaces.

?> More information may be returned in the future versions on the service

| Response value | Type   | Description           | Format                 |
| -------------- | ------ | --------------------- | ---------------------- |
| name           | string | Name of the namespace | Regex `^[A-Za-z0-9]+$` |

<!-- tabs:end -->

### Exists
```
    rpc Exists(IsNamespaceExistRequest) returns (IsNamespaceExistResponse) {};
```
Checks if namespace with provided name exists

<!-- tabs:start -->

#### **Request**

| Parameter name | Type   | Description                                                            | Required |
| -------------- | ------ | ---------------------------------------------------------------------- | -------- |
| name           | string | Unique name of the namespace to check                                  | True     |
| useCache       | bool   | Tell services if it can use cache in order to optimize response speed. | True     |

?> Using cache greatly reduses response times. Assuming you are are going to delete namespace only once, it is very good practice to use cache whenever its possible. Invalid cache will be deleted after short period of time (default is 60 seconds).

#### **OK**

Successfully checked if namespace exists.

| Response value | Type | Description            |
| -------------- | ---- | ---------------------- |
| exist          | bool | Namespace exist or not |

<!-- tabs:end -->
    

## Events
!> This feature is ***NOT IMPLEMENTED***

Events this service can rise

| system_rabbitmq exchange               | routing key | scheme               | conditions            |
| -------------------------------------- | ----------- | -------------------- | --------------------- |
| <GLOBAL_PREFIX>native_namespace_events | created     | Namespace (protobuf) | Namespace was created |
| <GLOBAL_PREFIX>native_namespace_events | deteled     | Namespace (protobuf) | Namespace was deleted |

## Configuration
This service is controled by environment variables.

| Env                                | default                                | description                                                                                                        |
| ---------------------------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| SYSTEM_DB_URL                      | mongodb://root:example@system_db/admin | [Mongo DB URL](https://www.mongodb.com/docs/manual/reference/connection-string/#standard-connection-string-format) |
| SYSTEM_CACHE_URL                   | redis://system_cache                   | System_cache redis connection URL                                                                                  |
| SYSTEM_TELEMETRY_EXPORTER_ENDPOINT | system_telemetry:55680                 | [OTEL connector](https://opentelemetry.io/docs/collector/) endpoint                                                |
| SYSTEM_AMQP_PREFIX                 | system_rabbitmq                        | System rabbitmq connection link                                                                                    |