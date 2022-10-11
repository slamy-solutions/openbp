# Key-value-storage service
The main task of this server is to persistently store key-value entries. The main reason for creating this server to keep values in such a way that at the moment of the write, they are 100% have been written to the disk. So when you are setting the value, you are guaranteed that this value is persisted on the disk.

This service is very useful for storing configurations and user preferences.

## API
Communication with the service is possible using gRPC interface. The native `keyvaluestorage` module provides definitions of the interface (proto file).
??? example "rpc Set(SetRequest) returns (SetResponse);"
    Sets the value for the specified key. Value is guaranteed to be stored persistently. If the key already holds some value, this operation will overwrite the value for this key.
    === "Request"
        | Parameter name | Type   | Description                                                                                                  |
        | -------------- | ------ | ------------------------------------------------------------------------------------------------------------ |
        | namespace      | string | Namespace where to store the key-value pair. It can be empty for a global key.                               |
        | key            | string | Unique key inside namespace to set value                                                                     |
        | value          | bytes  | Value to store. It cannot be larger than 15MB. Values that have more than one Mb in size will not be cached. |
    === "OK"
        Value was successfully seted or updated.

??? example "rpc Get(GetRequest) returns (GetResponse);"
    Gets value for key.
    === "Request"
        | Parameter name | Type   | Description                                                                                                                                                       |
        | -------------- | ------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
        | namespace      | string | Namespace where to search for the key-value pair. It can be empty for a global key.                                                                                    |
        | key            | string | Unique key inside namespace to get value                                                                                                                 |
        | useCache       | bool   | Use cache to get value or not. The cache may not be valid under rare conditions (simultaneous reads and writes). The cache is automatically cleared every 60 seconds. |
    === "OK"
        Value was successfully founded and returned.
        | Property name | Type  | Description  |
        | ------------- | ----- | ------------ |
        | value         | bytes | Stored value |
    === "NOT_FOUND"
        Namespace doesn't exist or there is no such key inside namespace 

## Configuration
This service is controlled by environment variables.

| Env                                | default                                | description                                                                                                        |
| ---------------------------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| SYSTEM_DB_URL                      | mongodb://root:example@system_db/admin | [Mongo DB URL](https://www.mongodb.com/docs/manual/reference/connection-string/#standard-connection-string-format) |
| SYSTEM_CACHE_URL                   | redis://system_cache                   | System_cache redis connection URL                                                                                  |
| SYSTEM_TELEMETRY_EXPORTER_ENDPOINT | system_telemetry:55680                 | [OTEL connector](https://opentelemetry.io/docs/collector/) endpoint                                                |