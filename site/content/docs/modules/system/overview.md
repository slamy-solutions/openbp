# System module
This module consists of services that provide low-level functionalities such as database management, security, caching, message queues ...

The system services include:

| Service             | Description                                                                                                                                                      |
| ------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [DB](./db.md)       | NoSQL database - main storage of all the information. Driven by [MongoDB](https://www.mongodb.com/).                                                             |
| [Cache](./cache.md) | In memory key-value DB for cache with automatic keys-eviction. Driven by [Redis](https://redis.io/).                                                             |
| [NATS](./nats.md)   | [NATS](https://nats.io/) service for distributed messages and queues.                                                                                            |
| [Redis](./redis.md) | General purpose, key-value, in memory [Redis](https://redis.io/) database.                                                                                       |
| [Vault](./vault.md) | Security service with HSM (hardware security module) integration for state of the art secrets protection, certificates management, encryption, signing and more. |

## Licensing

Most of these services are used as-is, without any modifications, to ensure that their licenses are not violated. The majority of them have Enterprise licenses with additional features. If you require additional functionalities, the enterprise services should integrate with OpenBP without any problems. These projects are open-source, so if you genuinely appreciate what OpenBP brings to the table, don't forget to support each of them separately.