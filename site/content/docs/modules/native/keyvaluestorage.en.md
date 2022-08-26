# Key-value-storage service
The main task of this server is to persistently store key-value entries. The main reason of creating this server to store values in such way that on the moment of the write, they are 100% have been written to the disk. So when you seting the value you are guaranteed that this value is persisted on the disk.

This service is very usefull for storing konfigurations or user preferences.

## API
Communication with the service is possible using gRPC interface. Definitions of the interface (proto file) are provided by native keyvaluestorage module.
### Set
```
    rpc Set(SetRequest) returns (SetResponse);
```
Sets the value for specified key. Value is guaranteed to be stored persistently.
If key already holds some value, this operation will overwrite value for this key.

<!-- tabs:start -->

#### **Request**

| Parameter name | Type   | Description                                                                                             |
| -------------- | ------ | ------------------------------------------------------------------------------------------------------- |
| namespace      | string | Namespace where to store the key-value pair. Can be empty for global key.                               |
| key            | string | Unique key inside namespace for wich to set value                                                       |
| value          | bytes  | Value to store. Can not be larger than 15mb. Values that have more than 1mb in size will not be cached. |

#### **OK**
Value was successfully seted or updated.

<!-- tabs:end -->

### Get
```
    rpc Get(GetRequest) returns (GetResponse);
```
Gets value for key.

<!-- tabs:start -->

#### **Request**

| Parameter name | Type   | Description                                                                                                                                                       |
| -------------- | ------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| namespace      | string | Namespace where to search for the key-value pair. Can be empty for global key.                                                                                    |
| key            | string | Unique key inside namespace for wich to get value                                                                                                                 |
| useCache       | bool   | Use cache to get value or not. Cache may not be valid under very rare conditions (simulatnious read and writes). Cache is automatically cleared every 60 seconds. |

#### **OK**
Value was successfully founded and returned.
| Property name | Type  | Description  |
| ------------- | ----- | ------------ |
| value         | bytes | Stored value |

#### **NOT_FOUND**
Namespace doesnt exist or there is no such key inside namespace 

<!-- tabs:end -->

## Configuration
This service is controled by environment variables.

| Env                                | default                                | description                                                                                                        |
| ---------------------------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| GLOBAL_PREFIX                      | openerp_                               | Prefix that will be applied to all created databases and amqp queues                                               |
| SYSTEM_DB_URL                      | mongodb://root:example@system_db/admin | [Mongo DB URL](https://www.mongodb.com/docs/manual/reference/connection-string/#standard-connection-string-format) |
| SYSTEM_CACHE_URL                   | redis://system_cache                   | System_cache redis connection URL                                                                                  |
| SYSTEM_TELEMETRY_EXPORTER_ENDPOINT | system_telemetry:55680                 | [OTEL connector](https://opentelemetry.io/docs/collector/) endpoint                                                |