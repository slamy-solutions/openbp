# System Redis Cache
System Redis cache is a [Redis](https://redis.io/) server for dealing with cache. It is configured to have max available memory and [LRU](https://redis.io/docs/manual/eviction/#eviction-policies) keys eviction policy. Check Redis [caching manual](https://redis.io/docs/manual/config/#configuring-redis-as-a-cache).

## Usage
Cache service is used directly throung the Redis client libraries. This ensures very hight speed and uses build-in Redis load balancing features.

### Configuration
This service is controled by build-in constant config with several environment variables.

| env       | range     | default | description                                                                                                              |
|-----------|-----------|---------|--------------------------------------------------------------------------------------------------------------------------|
| MAXMEMORY | 0mb - Xmb | 512mb   | Maximum ammount of memory that can be used by cache instance. After reaching this limit, service will start to drop keys |

### Keys
Cache access is shared between all the services so its very important to have unique cache keys withing entire system.

Please use `external_<developer>_<whatever + namespace>` format, where:
- developer is unique name of the deveporer. Developer name should not have "_" character. In that way modules installed by different developers will not have collisions
- whatever is your custom path. Dont forget, that cache service is also shared between namespaces. Add namespace name somewhere in the key.
Examples:
- external_scribesystems_rcp_integrator_devices_list_customnamespace
- external_scribesystems_rcp_integrator_devices_data_customnamespace_84529
- external_litesolutions_workload_participant_customnamespace_381204

Native modules names differ from the external. They start with prefix "native". Examples:
 - native_namespace_list_customnamespace
 - native_namespace_data_customnamespace_testing
 - native_cognito_users_list_customnamespace

### Values
DONT store big values in this cache. This service is optimized for small values. When you are storing big files here, it can evict hundreds of thousands of entries. If you have binary file/image/uncompressable data/something big (like more than 100 kilobytes compressed) - use system_redis_bigCache service for this. Remember, the smaller is expected value - the better.

Values in this service should be compressed with some algorithm to reduce memory footprint if possible. Compression must be implemented on the client side. The recomended way to go is [LZO](https://en.wikipedia.org/wiki/Lempel%E2%80%93Ziv%E2%80%93Oberhumer), because it provides fast decompression speeds. Also, if you are storing mostly JSON or Text data, check [Brotli](https://en.wikipedia.org/wiki/Brotli).

### Commands
You can use all the commands from Redis, that directly modify data. Rememder, that value under the key can be removed at any time.
Dont use Redis commands that are not related to the data, like configuration / options / etc ...

## Security
Cache can be accessed by every service and it stores data from all the service in same place. It is very important that services will not allow direct access to the keys by calers. Every service has responsibility to make sure, that under any circumstances it will not read and retrieve key that belongs to other service.

## Versions and compatibility
System Redis cache will ensure reverse compatibility of the key formats, so you dont have to verify this documentation for every new version :) New formats and key prefixes will not brake compatibility with old one and will not make new key collisions.

## Optimizations
As part of the "big values" optimization, server system should allocate more CPU resources, but less RAM (as it will not need it so much).

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.