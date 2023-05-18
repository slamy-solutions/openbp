# System cache service
System cache is a [Redis](https://redis.io/) server for dealing with the cache. It is configured to have max available memory and [LRU](https://redis.io/docs/manual/eviction/#eviction-policies) keys eviction policy. Check Redis [caching manual](https://redis.io/docs/manual/config/#configuring-redis-as-a-cache).

## Usage
The cache service can be directly accessed through Redis client libraries, ensuring high speed and utilizing built-in Redis load balancing features. Furthermore, OpenBP offers libraries that encapsulate Redis communication logic and provide a clearer interface with OTel integration.

### Configuration
This service is controlled by a built-in configuration that consists of various environment variables:

| env       | range     | default | description                                                                                                              |
|-----------|-----------|---------|--------------------------------------------------------------------------------------------------------------------------|
| MAXMEMORY | 0mb - Xmb | 512mb   | Maximum ammount of memory that can be used by cache instance. After reaching this limit, service will start to drop keys |

### Keys
To ensure unique cache keys throughout the entire system, it is crucial to use the `external_<developer>_<whatever + namespace> `format. Dont forget to include namespace name to the key. Here are some examples:

- external_scribesystems_rcp_integrator_devices_list_customnamespace
- external_scribesystems_rcp_integrator_devices_data_customnamespace_84529
- external_litesolutions_workload_participant_customnamespace_381204
  
For native modules, the names differ and start with the prefix "native". Examples:

- native_namespace_list_customnamespace
- native_namespace_data_customnamespace_testing
- native_cognito_users_list_customnamespace
  
By following this naming convention, collisions between modules installed by different developers can be avoided, and the cache access remains shared appropriately.

### Values
It is important not to store large values in this cache as the service is optimized for small values. Storing big files can result in the eviction of hundreds of thousands of entries. Remember that smaller values are preferred for optimal performance.

When storing values in this service, it is recommended to compress them using a suitable algorithm to reduce memory footprint. Compression should be implemented on the client side. The recommended approach is to use [LZO](https://en.wikipedia.org/wiki/Lempel%E2%80%93Ziv%E2%80%93Oberhumer) as it provides fast decompression speeds. If you primarily store JSON or text data, you may also consider using [Brotli](https://en.wikipedia.org/wiki/Brotli) compression.

### Commands
You are free to use all Redis commands that directly modify data in the cache. However, please keep in mind that the value associated with a key can be removed at any time.

Avoid using Redis commands that are not related to data manipulation, such as configuration or options settings. Stick to the commands that interact with and modify the stored data.

## Security
The cache is accessible by every service, and it serves as a centralized storage for data from all services. It is crucial to ensure that services do not allow direct access to keys by callers. Each service carries the responsibility of ensuring that it never reads or retrieves keys that belong to other services, under any circumstances. Implement proper access controls and safeguards within each service to maintain data isolation and prevent unauthorized access to keys belonging to other services.

## Versions and compatibility
The system cache is designed to ensure reverse compatibility of key formats. This means that you don't have to refer to the documentation for every new version to verify the key format. The introduction of new formats and key prefixes will not break compatibility with the old ones, and it will not cause collisions with existing keys.

You can rely on the system cache to handle key format changes seamlessly, allowing you to update and evolve your application without worrying about compatibility issues or key collisions.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.