# OpenBP runtime module library

This library is created for inter-service communication in the OpenBP.
You can only use this library for developing new OpenBP modules because it establishes a connection with services using internal communication protocols.
This will not work if you want to communicate with OpenBP from external application. For this use-case check `@openbp/sdk`.

## Basic usage
```ts
import { ERPStub } from '@openbp/runtime'

// Creating services communication "proxy". This will only connect to the selected services
const stub = new ERPStub()
await stub.connect()

// Now you can communicate with runtime services
const runtimes = []
await stub.services.manager.core.GetAllRuntimes().forEach((r) => runtimes.push(r))
console.log(runtimes)

// Dont forget to close connection
stub.close()
```