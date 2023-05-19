# Modules

"Module" is a top-level abstraction that distinguishes resources within the OpenBP system. Here are the key aspects of a module:

- A module can be enabled or disabled; there is no partial state.
- The same module cannot be enabled multiple times.
- Modules can be either internal or external. Internal modules refer to those developed within the OpenBP repository.
- Modules form the core of the system. There are no security checks enforced between modules, allowing all modules to access each other freely.

## Internal modules
Here is the list of available internal modules:

| Service                                   | Description                                                                                                                                                 |
| ----------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ***[`System`](./system/overview.md)***    | Provides low-level functionalities such as database management, security, caching, message queues ...                                                       |
| ***[`Native`](./native/overview.en.md)*** | Provides extensive collection of services and tools that will most likely be used in every module or solution. For example: authorization, namespacing, ... |
