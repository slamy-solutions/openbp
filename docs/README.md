# OpenBP - Open Business platform

## Motivation
Have you ever created or integrated something to SAP / Comarch ERP / 1C / Odoo / …. ? Have you noticed how many years these platforms are behind? Some of them were created even before virtual machines appeared, I'm not even talking about Docker or Kubernetes :) The use of old patterns, dependencies, and very strong corporatization leads to the fact that they will not change.

Many new large enterprises are looking for which platforms to switch to and simply cannot find alternatives. The business cannot find developers to support the systems, simply because the person who works with these platforms is simply not called a developer. Most of them even have proprietary programming languages ​​for which there are no libraries anywhere on the internet. Seriously, why use ABAP when you can use Typescript?

These platforms do not use what has long been considered the norm for developers, for example, containerization, GIT for version control, popular programming languages, IaaC, CI/CD, and SPA as frontend on React/Vue/Angular, …. . They implement modern patterns only through dirty hacks. Such simple tasks as creating an API, opening a WebSocket, adding a module, setting up a high-availability cluster, and performing an update on the go - for these systems often turn into a nightmare.

We need a platform that will cut off old redundant business patterns and provide all new state-of-the-art capabilities out of the box.

## Capabilities
Here's what OpenBP offers out of the box:
- **SaaS-enabled** by default
- **Multitenant and multi namespace** system capable to manage multiple clients with non-standard/modified configurations.
- Unlimited Horizontal and Vertical **scalability** by default.
- Hight-availability option to ensure **99.9% uptime**.
- Dynamic workflow managing, resource sharing, and workload redistribution (provided by k3s)
- Industrial grade secured environment (provided by Vault)
- **Developer-friendly version-controlled** environment (with [GIT](https://git-scm.com/)) to create custom configurations.
- **CI/CD** pipelines for easy deployments and support
- Advanced **SaaS** tools for managing customers and their deployments.
- **Traceability** dashboards with current system status, logs, events, and errors.
- Intense set of highly detailed **documentation**.
- gRPC interservice communication. Use whatever language you want :)
- Use of state-of-the-art technologies and practices: [MongoDB](https://www.mongodb.com/), [RedisDB](https://redis.io/) cahing, Hashicorp [Vault](https://www.vaultproject.io), [RabbitMQ](https://www.rabbitmq.com/) messaging and task manager, [Docker](https://www.docker.com/), [k3s](https://k3s.io/), [OpenAPI3.0](https://swagger.io/specification/), [gRPC](https://grpc.io/) ....

Thats not all ;). Please, check [Architectual vision](./architecture/architecture_vision.md) if you want to know more.

## Sponsors
<div align="center">
  <a href="https://modulsoft.pl/">
    <img src="./sponsors/modulsoft.svg" width="20%" />
  </a>
</div>


## License
OBP is free and the source is available.
Software released under [Server Side Public License (SSPL) v1](LICENSE.md) license. This is GNU GPL v3 with small modifications provided by [MongoDB](https://www.mongodb.com/licensing/server-side-public-license/faq).