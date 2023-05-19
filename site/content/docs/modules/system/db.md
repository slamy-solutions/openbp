# System db service

The `system_db` is a ***[MongoDB](https://www.mongodb.com/)*** service designed to serve as a general-purpose database for various functionalities within the OpenBP project. It acts as the main persistence provider, ensuring data durability and availability. The container of the MongoDB is used as is to ensure, that the license is not violated. We suggest checking MongoDB Enterprise because it has a lot of features that you would probably want to incorporate into your project.

## Resource naming convention
To ensure that developers' resources do not overlap, it is important to establish a naming convention.

The service utilizes two types of databases:

- Namespaced: As OpenBP implements multitenancy, the service divides information between each tenant. This allows for easy understanding of the database usage by each tenant. The naming convention for these databases is `openbp_namespace_<namespace name>`.
- Global: There is only one global database per installation, which contains information relevant to all tenants. The name of this database is `openbp_global`.

While you have the freedom to create additional databases, it is advised to exercise caution, as this may negatively impact other services within OpenBP, particularly those related to technical aspects such as statistics and backup.

To prevent overlaps with other developers, it is important to use the ***`external_<org name>_<project name>_`*** prefix for collections within the databases. This ensures that collections are unique to each developer. However, it's worth noting that OpenBP internal collections will have a shorter name starting with the module name, but they will never begin with the `external_ prefix`.