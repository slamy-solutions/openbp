# Namespace service
The `native_namespace` is a service for managing namespaces in the entire platform ecosystem. Namespace provides logical separation of the data for modules and services. They can be used, for example, to distinguish customers in the application.
After the creation of the namespace, the `native_namespace` service creates a new `openpb_namespace_<name>` database so that every new namespace data can be placed in a separate database. On `namespace` deletion, the service will delete its database, and all the data will be lost.

## Interaction examples
Those are several examples of how you can interact with the system from other kernel systems or modules.

??? example "Basic communication"
    In this example we will create new namespace. Then we will check if created namespace exists

    === "GO"
        ```golang
        package main

        import (
            "context"

            native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
            "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
        )

        const NATIVE_NAMESPACE_URL = "native_namespace:80"

        func main() {
            // Connect to the namespace service
            nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService())
            nativeStub.Connect()
            defer nativeStub.Close()

            // Create new namespace
            nativeStub.Services.Namespace.Create(context.TODO(), &namespace.EnsureNamespaceRequest{
                Name: "mynamespace",
                FullName: "My new namespace :)",
                Descriptions: "I created this namespace just for fun XD"
            })

            // Check if namespace was created
            response, _ := nativeStub.Services.Namespace.Exists(context.TODO(), &namespace.IsNamespaceExistRequest{
                Name:     "mynamespace",
                UseCache: true,
            })

            print(response.Exist)
        }
        ```

??? example "Listen for events"
    In this example, we will listen for events on the namespace service

    === "GO"
        ```golang
        package main

        import (
            "context"
            "fmt"

            "github.com/nats-io/nats.go"
            native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
            system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
            "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
        )


        func main() {
            // Connect to the namespace service
            nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService())
            nativeStub.Connect()
            defer nativeStub.Close()

            // Connect to NATS
            systemStub := system.NewSystemStub(system.NewSystemStubConfig().WithNats())
            systemStub.Connect()
            defer systemStub.Close()

            // Start listen for namespace creation events
            jet, _ := systemStub.Nats.JetStream()
            consumer, _ := jet.AddConsumer("native.namespace.event", &nats.ConsumerConfig{FilterSubject: "native.namespace.event.created", Name: "mymodule.myservice.namespace_consumer"})
            sub, _ := jet.SubscribeSync("", nats.Bind(consumer.Stream, consumer.Name))

            // Create new namespace
            nativeStub.Services.Namespace.Ensure(context.TODO(), &namespace.EnsureNamespaceRequest{
                Name: "mynamespace",
            })

            // Receive event with created namespace
            msg, _ := sub.NextMsg(5 * time.Second)
            var namespaceData namespace.Namespace
            proto.Unmarshal(msg.Data, &namespaceData)
            fmt.Printf("Created [%s] namespace\n", namespaceData.Name)
            msg.Ack()
        }
        ```

!!! info
    You can connect to the system using other languages than those provided in this example. You will have to manually compile [protobuf](https://developers.google.com/protocol-buffers) definitions of this service in the language you want. For more information, please have a look at [Use the language you want](../../concepts/languageYouWant.en.md) concept.

## Schema
The schema of the namespace is defined using `protobuf`. Definitions are provided by the `native` module.

`Namespace` schema:

| Property    | Type      | Description                                                         |
| ----------- | --------- | ------------------------------------------------------------------- |
| name        | string    | Unique short name.                                                  |
| fullName    | string    | Full, public name.                                                  |
| description | string    | Arbitrary description.                                              |
| created     | timestamp | Timestamp when namespace was created                                |
| updated     | timestamp | Timestamp when namespace was updated last time                      |
| version     | uint64    | Counter that increases on every update of the namespace information |

## API
Communication with the system is possible using the gRPC interface. Definitions of the interface (proto file) are provided by the `native` module.

??? example "rpc Create(CreateNamespaceRequest) returns (CreateNamespaceResponse) {};"
    Creates new namespace. Returns error if the namespace with same name already exists.

    === "Request"
        | Parameter name | Type   | Description                                          | Required | Format                                        |
        | -------------- | ------ | ---------------------------------------------------- | -------- | --------------------------------------------- |
        | name           | string | Unique name of the namespace. Can't be changed later | True     | Regex `^[A-Za-z0-9]+$`. Maximum 32 characters |
        | fullName       | string | Full, public name. Can be changed later.             | False    | Maximum 128 characters                        |
        | description    | string | Arbitrary description. Can be changed later.         | False    | Maximum 512 characters                        |
    === "OK"
        Namespace and the new `openbp_namespace_<name>` database were successfully created.

        | Response value | Type                 | Description             |
        | -------------- | -------------------- | ----------------------- |
        | namespace      | [Namespace](#Schema) | Newly created namespace |

        !!! info
            This response will raise the event on the [`system_nats`](../system/nats.en.md). Check the [Events](#events) section and the specific event for the namespace creation.
    === "ALREADY_EXISTS"
        Namespace with the same name already exists.
    === "IVALID_ARGUMENT"
        At least one of the parameters has a bad format.

??? example "rpc Ensure(EnsureNamespaceRequest) returns (EnsureNamespaceResponse) {};"
    Creates new namespace if it doesn't exist. Does nothing if the namespace already exists.

    === "Request"
        | Parameter name | Type   | Description                                          | Required | Format                                        |
        | -------------- | ------ | ---------------------------------------------------- | -------- | --------------------------------------------- |
        | name           | string | Unique name of the namespace. Can't be changed later | True     | Regex `^[A-Za-z0-9]+$`. Maximum 32 characters |
        | fullName       | string | Full, public name. Can be changed later.             | False    | Maximum 128 characters                        |
        | description    | string | Arbitrary description. Can be changed later.         | False    | Maximum 512 characters                        |
    === "OK"
        Namespace was successfully created (or it existed before request). Service created new `openbp_namespace_<name>` database.

        | Response value | Type                 | Description                                                                  |
        | -------------- | -------------------- | ---------------------------------------------------------------------------- |
        | namespace      | [Namespace](#Schema) | Newly created namespace                                                      |
        | created        | bool                 | Indicates if namespace was created or it already existed before this request |

        !!! info
            This response will raise the event on the [`system_nats`](../system/nats.en.md). Check the [Events](#events) section and the specific event for the namespace creation.
    === "IVALID_ARGUMENT"
        At least one of the parameters has a bad format.

??? example "rpc Update(UpdateNamespaceRequest) returns (UpdateNamespaceResponse) {};"
    Updates information of the existing namespace.

    === "Request"
        | Parameter name | Type   | Description                     | Required | Format                 |
        | -------------- | ------ | ------------------------------- | -------- | ---------------------- |
        | name           | string | Name of the namespace to update | True     | -                      |
        | fullName       | string | New full, public name           | False    | Maximum 128 characters |
        | description    | string | New description.                | False    | Maximum 512 characters |
    === "OK"
        Namespace was successfully updated.

        | Response value | Type                 | Description                            |
        | -------------- | -------------------- | -------------------------------------- |
        | namespace      | [Namespace](#Schema) | Namespace information after the update |

        !!! info
            This response will raise the event on the [`system_nats`](../system/nats.en.md). Check the [Events](#events) section and the specific event for the namespace update.
    === "NOT_FOUND"
        Namespace with this name wasn't founded.
    === "IVALID_ARGUMENT"
        At least one of the parameters has a bad format.

??? example "rpc Delete(DeleteNamespaceRequest) returns (DeleteNamespaceResponse) {};"
    Deletes namespace and corresponding database.
    !!! danger
        This is very danger operation. Service will delete the entire database and all data of all other services.
    === "Request"
        | Parameter name | Type   | Description                                |
        | -------------- | ------ | ------------------------------------------ |
        | name           | string | Unique name of the namespace to be deleted |
    === "OK"
        Namespace was successfully deleted or wasn't exist. Service has removed the entry from the `openbp_global` database and deleted `openbp_namespace_<name>` database.

        | Response value | Type | Description                                                 |
        | -------------- | ---- | ----------------------------------------------------------- |
        | existed        | bool | Indicates if namespace existed before this operation or not |

        !!! info
            This response will raise an event on the [`system_nats`](../system/nats.en.md). Check the [Events](#events) section and specific events raised on namespace deletion.

??? example "rpc Get(GetNamespaceRequest) returns (GetNamespaceResponse) {};"
    Returns namespace information by its name
    === "Request"
        | Parameter name | Type   | Description                                                            |
        | -------------- | ------ | ---------------------------------------------------------------------- |
        | name           | string | Unique name of the namespace to get                                    |
        | useCache       | bool   | Tell services if it can use cache in order to optimize response speed. |
        !!! tip
            Using cache greatly reduces response times. The cache may not be valid under rare conditions. Cache may not be valid only when there are a lot of read/write operations on namespace information simultaneously. In the real system, `namespace` probably will never change, so it's excellent practice to use cache wherever you can. The invalid cache deletes after some time (60 seconds by default). The system tries to keep the cache up to date, but it cants ensure its validity all the time.
    === "OK"
        Namespace was founded, and its data was successfully returned.

        | Response value | Type                 | Description           |
        | -------------- | -------------------- | --------------------- |
        | namespace      | [Namespace](#Schema) | Namespace information |
    === "NOT_FOUND"
        Service can't find this namespaces. Probably, it doesn't exist.

??? example "rpc GetAll(GetAllNamespacesRequest) returns (stream GetAllNamespacesResponse) {};"
    Gets stream of all namespaces
    === "Request"
        Request doesn't has parameters. It just returns all the namespaces.
    === "OK"
        Successfully got a stream of the namespaces.

        | Response value | Type                 | Description           |
        | -------------- | -------------------- | --------------------- |
        | namespace      | [Namespace](#Schema) | Namespace information |

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

??? example "rpc Stat(GetNamespaceStatisticsRequest) returns (GetNamespaceStatisticsResponse) {};"
    Get namespace statistics.
    === "Request"
        | Parameter name | Type   | Description                                                            |
        | -------------- | ------ | ---------------------------------------------------------------------- |
        | name           | string | Unique name of the namespace to get                                    |
        | useCache       | bool   | Tell services if it can use cache in order to optimize response speed. |
    === "OK"
        | DB statistic | Type   | Description                                                                          |
        | ------------ | ------ | ------------------------------------------------------------------------------------ |
        | objects      | uint64 | Number ob objects stored in the database                                             |
        | dataSize     | uint64 | Total size of the raw data stored (without pre-alocated allocated space and indexes) |
        | totalSize    | uint64 | Total memory usage of the namespace                                                  |

        | Response value | Type          | Description                                            |
        | -------------- | ------------- | ------------------------------------------------------ |
        | db             | DB statistics | Information about the database usage of this namespace |
    === "NOT_FOUND"
        Namespace wasnt founded. Probably, it doesn't exist.

## Events

This service raises events on the [`system_nats`](../system/nats.en.md) service.

| stream                 | subject                        | scheme               | conditions            |
| ---------------------- | ------------------------------ | -------------------- | --------------------- |
| native_namespace_event | native.namespace.event.created | Namespace (protobuf) | Namespace was created |
| native_namespace_event | native.namespace.event.updated | Namespace (protobuf) | Namespace was updated |
| native_namespace_event | native.namespace.event.deleted | Namespace (protobuf) | Namespace was deleted |

## Configuration
This service is controlled by environment variables.

| Env                                | default                                | description                                                                                                        |
| ---------------------------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| SYSTEM_DB_URL                      | mongodb://root:example@system_db/admin | [Mongo DB URL](https://www.mongodb.com/docs/manual/reference/connection-string/#standard-connection-string-format) |
| SYSTEM_CACHE_URL                   | redis://system_cache                   | System_cache redis connection URL                                                                                  |
| SYSTEM_TELEMETRY_EXPORTER_ENDPOINT | system_telemetry:55680                 | [OTEL connector](https://opentelemetry.io/docs/collector/) endpoint                                                |
| SYSTEM_NATS_URL                    | nats://system_nats                     | [Nats](../system/nats.en.md) endpoint                                                                              |