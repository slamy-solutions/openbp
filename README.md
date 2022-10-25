# OpenBP - Open Business platform

## Motivation
Have you ever created or integrated something to SAP / Comarch ERP / 1C / Odoo / …. ? Have you noticed how many years these platforms are behind? Some of them were created even before virtual machines appeared, and I'm not even talking about Docker or Kubernetes :) The use of old patterns, dependencies, and very strong corporatization leads to the fact that they will not change.

Many new large enterprises are looking for which platforms to switch to and cannot find alternatives. The business cannot find developers to support the systems simply because the person who works with these platforms is not called a developer. Most even have proprietary programming languages ​​for which there are no libraries on the internet. Seriously, why use ABAP when you can use Typescript?

These platforms do not use what has long been considered the norm for developers, for example, containerization, GIT for version control, popular programming languages, IaC, CI/CD, and SPA as frontend on React/Vue/Angular, …. . They implement modern patterns only through dirty hacks. Such simple tasks as creating an API, opening a WebSocket, adding a module, setting up a high-availability cluster, and performing an update on the go - for these systems often turn into a nightmare.

We need a platform that will cut off old redundant business patterns and provide all new state-of-the-art capabilities out of the box.

## Capabilities
Here's what OpenBP offers out of the box for developers, implementation companies and business:

### Developers
- Created by developers for developers. We want this to be as easy and beautifulll as possible. 
- **Developer-friendly version-controlled** environment (with [GIT](https://git-scm.com/)) to create custom configurations.
- **CI/CD** pipelines for easy deployments and support
- **OpenTelemetry** tracing to debug
- Intense set of highly detailed **documentation**
- gRPC interservice communication. Use whatever language you want :)
- There are no limit in what you can create. Just create one more microservice and OpenBP will handle everything.
- Use of state-of-the-art technologies and practices: [MongoDB](https://www.mongodb.com/), [RedisDB](https://redis.io/) cahing, Hashicorp [Vault](https://www.vaultproject.io), [RabbitMQ](https://www.rabbitmq.com/) messaging and task manager, [Docker](https://www.docker.com/), [k3s](https://k3s.io/), [OpenAPI3.0](https://swagger.io/specification/), [gRPC](https://grpc.io/) ....

### Implementation companies
- **SaaS-enabled** by default. One place for all customers
- **Multitenant and multi namespace** system capable to manage multiple clients with non-standard/modified configurations.
- Dynamic workflow managing, resource sharing, and workload redistribution. All customers are using same resources and are working on the same cluster. Very efficient and easy to maintain.
- Advanced **SaaS** tools for managing customers and their deployments.
- **Traceability** dashboards with current system status, logs, events, and errors.
- Automatic billing based on resource usage

### Business
- Open source and free
- Unlimited **scalability** by default. Dont think about the scale of your bisness and ammount of data. This system will handle everything.
- Hight-availability option to ensure **99.9% uptime**.
- **On-the-fly updates** that will not stop your critical workflows.
- Industrial grade security environment complied with **ISO/IEC 27001** requirements.
- **GDPR** complied storage system and data flow.
- Highly flexible and extendable.
- Integrated with edge computing and device management.
- Includes set of the solutions: ERP, CRM, Edge computing, etc.

That's not all ;). Please, check [Architectural vision](./docs/architecture/architecture_vision.md) if you want to know more.

## Sponsors
<div align="center">
  <a href="https://modulsoft.pl/">
    <img src="./site/content/assets/images/sponsors/modulsoft.svg" width="20%" />
  </a>
</div>


## License
OpenBP is free and open source. It is distributed under AGPL3 license. If you need an Enterprise license - contact us.