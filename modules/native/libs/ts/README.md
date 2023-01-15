# OpenBP native module library

This library is created for inter-service communication in the OpenBP.
You can only use this library for developing new OpenBP modules because it establishes a connection with services using internal communication protocols.
This will not work if you want to communicate with OpenBP from external application. For this use-case check `@openbp/sdk`.

## Basic usage
```ts
import { NativeStub, Service } from '@openbp/native'

// Creating services communication "proxy". This will only connect to the selected services
const stub = new NativeStub([Service.NAMESPACE])
await stub.connect()

// Now you can communicate with native services
const namespaces = []
await stub.services.namespace.GetAll({ useCache: true }).forEach((n) => namespaces.push(n))
console.log(namespaces)

// Dont forget to close connection
stub.close()
```