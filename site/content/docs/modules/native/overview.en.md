# Native module

The `native` module is an extensive collection of services and tools that will most likely be used in every module or solution. The services found here do not have any specifics of use and are inert to the final product (that is, they do not have dependencies and the context of what exactly will be done).

!!! note
    For example, a whole subgroup of "IAm" services is responsible for authorization and authentication. Such services are likely to be used in any final solution or module.

Therefore, we can say that this module is a collection of `meta-services` that are most often used as tools for other modules or solutions.

!!! warning
    There is a big difference between the `system` and `native` modules. Services of `native` module arent external solutions (entire code is inside this repository).

To summarize, those attributes service must have to be in the `native` module:

- **Abstract**. Doesn't depend on the final solution. Applicable to everything.
- **Internal**. The entire code is located in this repository.
- **Dependency-less**. The only dependency is the `system` module.