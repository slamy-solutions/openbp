# The goal of the project
The goal of the project is to create a competitive high-tech platform capable of solving modern business problems. OBP allows the creation of an ecosystem with ERP/CRM/Edge computing/etc solutions for businesses.

The bird-eye view of the project is presented in the diagram below.

The project focuses on delivering a structure capable of handling business solutions. Also, the project provides basic modules with the most usable services.

Solutions that will be delivered at the end of the project include:

Solutions that will be delivered at the end of the project include:
- **Multitenant and multi namespace** system capable to manage multiple clients with non-standard/modified configurations using the same SaaS environment.
- The system with unlimited Horizontal and Vertical **scalability** by default.
- Hight-availability option to ensure **99.9% uptime**.
- A system capable of dynamic workflow managing, resource sharing, and workload redistribution
- Industrial grade security environment complied with **ISO/IEC 27001** requirements
- **GDPR** complied storage system and data flow
- **Developer-friendly version-controlled** environment (with [GIT](https://git-scm.com/)) to create custom configurations.
- **CI/CD** pipelines for easy deployments and support
- Advanced **SaaS** tools for managing customers and their deployments.
- **Traceability** dashboards with current system status, logs, events, and errors.
- Ready to use “system” services for the developers:
    - **DB** - multi-cluster, scalable, highly available database
    - **Cache** - highly available in-memory key-value storage cluster designed for “cache” with LRU cache eviction strategy.
    - **BigCache** -  highly available in-memory key-value storage cluster designed for “cache” with LRU cache eviction strategy. Can be used to efficiently store big cache entries (like binary files).
    - **Vault** - highly available secrets storage engine with HSM capabilities and implementing  ISO/IEC 27001 requirements. 
    - **AMQP** - distributed queue service with AMQP communication protocol.
    - **Entrypoint** - high available reverse proxy with load balancing, dynamic paths configuration, and automatic SSL certificates.
    - **OTEL** - OpenTelemetry scraper service
- Ready to use “native” services for the developers:
    - **Namespace** - managing the multi-tenant environment
    - **IAm** - a group of services to manage authentication and authorization; implementing authentication methods using OAuth2, SSO, X509 certificates, etc; managing identity policies.
    - **KeyValueStorage** - persistent and reliable key-value storage for internal uses.
    - **File** - s3-like storage for storing big chunks of data.
    - **Lambda** - a group of services for serverless computing.
    - **Configuration** - a group of services to manage/update/apply configurations for namespaces.
    - **API** - default, build-in native API.
    - **BFF** - dynamic API used by configurations to implement BFF (backend for frontend) pattern
    - **Cron** - task scheduler
    - **Backup** - managing backups and data snapshots
- Intense set of highly detailed **documentation** and courses for developers.

# Selected quality attributes
Selected quality attributes, which will be implemented by the OpenERP, originate from best practices for creating IT systems. They were selected based on the reference list **ISO 25000** of quality attributes:
- **Interoperability (Compatibility)** - business automation solutions exist for a long period of time. They are a lot of already working components and mechanisms that can not be replaced or their replacement is too costly. That means a new solution has to be able to easily adapt and use third-party business components. However, compatibility has not to be forced at the price of cutting functionality or complicating workflow. The solution uses a “new-first” approach to adapt to the new market needs and in some cases drops compatibility with redundant business patterns.
- **Learnability (Usability)** - this is one of the key factors that will distinguish the product from its competitors. The product will try to make the learning curve as smooth as possible. The system will ensure the best possible level of effectiveness, efficiency, freedom from risk, and satisfaction by default. The automatic level of ensuring the quality of the processes will only be bounded and limited by other quality attributes. Learnability has less priority than other quality attributes in all conflicting situations.
- **Operability (Usability)** - the system has to only use the libraries and tools that are well known to the developers and have a big community. Creating solutions based on this platform has to have the same workflow/or as close as possible to the workflow with ordinary solutions (like using GIT,CI/CD,....). With the support of a multi-tenant environment, the number of managed namespaces can grow dramatically. Same for the complexity of the overall system. That means the solution has to provide a full set of tools to easily work with SaaS infrastructure. The platform also provides traceability tools for developers to make it easier to handle edge cases and complex scenarios.
- **User error protection (Usability)** - the system provides error protection by creating clearer structure and simpler usage patterns. Error protection is guaranteed inside module services by validating inputs that can hurt internal structure. On other hand, the system doesn't guarantee error protection while communicating between modules services. That was done in order to maintain a clearer structure of the system.
- **Availability (Reliability)** - product has to have the option to work with 99.9% uptime. Those requirements are dictated by the scale of clients and their 24/7 workflows, which cannot be terminated, or termination is more costly than the cost of maintaining a high-available setup.
- **Recoverability (Reliability)** - platform potentially handles critical business information and workflows. The system has to be able to automatically recover critical failures as fast as possible and notify operators about incidents. Business data recoverability has to be guaranteed. The end user can decide what level of recoverability he wants and at what cost.
- **Confidentiality (Security)** - product provides an extensive set of authorization and authentication tools to ensure a high level of security. The solution has to use best practices from the ISO/IEC 27001 specification
- **Modularity (Maintainability)** - system work in a highly dynamic environment that ofter differs from client to client. That means, the system has to be highly modular so those modules can be adjusted to client needs with minimal impact on other components.

# Selected functional attributes

The proposed functional attributes focus on material aspects and benefits that give the product the ability to perform its intended task better and more efficiently than others, thus providing measurable value for consumers. Attributes are separated into categories in order to make a clearer overall view.

## Native services
### Namespace
- Create new namespaces
- Get a namespace by name
- Delete namespace
- List/Search namespaces
### IAm
Manages authentication and authorization information for identities. Provides API for other services to protect them from unauthorized usage.
#### Configuration
- Securely store OAuth2 provider secrets
- Store custom configuration parameters for IAm services
#### Policy
- Create global policies and policies in a specified namespace. The policy must contain information about what resources can be accessed by policy and what actions can be performed on those resources.
- Get an existing policy by its unique identifier and namespace
- Change policy resources and actions at runtime
- Delete policies at runtime
- List/Search for policies
#### Identity
- Ability to create new global identities and identities for the specified namespace. Identity must contain information about what policies are assigned to this identity
- Get existing identity by its unique identifier and namespace
- Change policies assigned to the identity at runtime
- Block/Unblock identities at runtime
- List/Search for identities
#### Auth
- Create an access token using authentication and authorization mechanisms
- Use login/password as one of the authentication mechanism
- Use OAuth2 as one of the authentication mechanisms. Use github.com, google.com, facebook.com, and gitlab.com as possible authentication providers.
- Use SSO as one of the authentication mechanisms.
- Add the ability to use two-factor authentication mechanisms.
- Add TOTP two-factor authentication mechanism.
- Provide authorization mechanisms for other services
### KeyValueStorage
- Reliably create key-value pairs and store them in persistent storage
- Get existing values by the corresponding keys
- Remove keys
### File
- Create and upload arbitrary binary data (as file) to the persistent storage
- Get file by unique identifier supporting reading from the middle
- Get file statistics (size, hash sum, name, metadata)
- Delete files
### Lambda
#### Manager
- CRUD operations for global lambdas and lambdas in namespaces
#### Entrypoint
- Call lambda without waiting for the response
- Execute lambda and get a response
- Stream data to lambda and wait for the response
- Execute lambda and read stream response
- Bidirectional stream with lambda execution
### Configuration
- CRUD operations for configurations
- Apply configuration to namespace
- Get the list of configurations applied to the namespace
### API
Provide REST API for the native services
### BFF
- CRUD user-defined APIs
- Show OpenAPI specification for defined APIs if possible
### Cron
- Schedule one-time task
- Schedule repeatable task
- Check task results/statistics/errors
### Backup
- Create namespace data snapshots
- Delete old snapshots
- Download snapshots
- Apply data snapshot to the namespace

# Technical architecture
