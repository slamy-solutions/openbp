# System nats service

The `system_nats` service is a managed service that leverages ***[NATS](https://nats.io/)***, a versatile messaging system. The service is preconfigured with default settings and utilizes unmodified NATS containers. Moreover, it incorporates the ***JetStream*** extension to enable data persistency.

## Resource naming convention

To ensure proper resource management and prevent overlapping resource names within the shared `system_nats` service, it is crucial to establish a consistent naming convention. This convention becomes particularly important when developing your own modules and services.

For your organization's modules and services, it is recommended to use the ***`external_<your org name>_<project name>_`*** prefix for all the resource names. This prefix helps clearly identify and differentiate your resources from others, minimizing the risk of conflicts. By incorporating your organization and project names into the prefix, you create a unique namespace within the NATS system.

On the other hand, resource names generated by OpenBP modules will have a shorter format and will directly start with the module name. However, these resources will never begin with the ***`external_`*** prefix.