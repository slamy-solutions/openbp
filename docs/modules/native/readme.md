# Native module

The `native` module is a large collection of services and tools that will most likely be used in every module or solution. The services found here do not have any specifics of use and are inert to the final product (that is, they do not have dependencies and the context of what exactly will be done).

?> For example, there is a whole subgroup of "IAm" services that is responsible for authorization and authentication. Such services are very likely to be used in any final solution or module.

Therefore, we can say that this module is a collection of `meta-services` that are most often used as tools for other modules or solutions.

!> There is big difference between `system` and `native` modules. Services of `native` module arent external solutions (entire code is inside this repository).

To summarize, what attributes service must have to be in the `native` module:
- **Abstract**. Doesnt depend on the final solution. Applyable to everything.
- **Internal**. Entire code is located in this repository.
- **Dependency-less**. The noly dependency is `system` module.