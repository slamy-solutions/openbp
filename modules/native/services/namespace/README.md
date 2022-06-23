# Native namespace
Native namespace is service for managing namespaceses in the entire application. Namespace provides logical separation of the data. They can bu used for example to distinguishe customers in the application. Communication with system is possible by gRPC interface. Definitions of the interface (proto file) is provided by native module.

## Usage
After creation of the `namespace`, native_namespace service creates entry in the `<GLOBAL_PREFIX>global` datadase and new `<GLOBAL_PREFIX>namespace_<name>` database, so every new namespace data can be placed in separate database. On `namespace` deletion, service will delete undelying `namespace` database

## Events
Events that are rising by this service are

| system_rabbitmq exchange                 | routing key | scheme               | conditions            |
|------------------------------------------|-------------|----------------------|-----------------------|
| <GLOBAL_PREFIX>native_namespace_events   | created     | Namespace (protobuf) | Namespace was created |
| <GLOBAL_PREFIX>native_namespace_events   | deteled     | Namespace (protobuf) | Namespace was deleted |

## Configuration
This service is controled by environment variables.

| env                                | default                                | description                                                                                                        |
|------------------------------------|----------------------------------------|--------------------------------------------------------------------------------------------------------------------|
| GLOBAL_PREFIX                      | openerp_                               | Prefix that will be applied to all created databases and amqp queues                                               |
| SYSTEM_DB_URL                      | mongodb://root:example@system_db/admin | [Mongo DB URL](https://www.mongodb.com/docs/manual/reference/connection-string/#standard-connection-string-format) |
| SYSTEM_CACHE_URL                   | redis://system_cache                   | System_cache redis connection URL                                                                                  |
| SYSTEM_TELEMETRY_EXPORTER_ENDPOINT | system_telemetry:55680                 | [OTEL connector](https://opentelemetry.io/docs/collector/) endpoint                                                |
| SYSTEM_AMQP_PREFIX                 | system_rabbitmq                        | System rabbitmq connection link                                                                                    |