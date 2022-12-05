# Deploy for production using Raspberry Pi

OpenBP has natively compiled containers for the `linux/arm64/v8` architecture which is used in Raspberry Pi boards. You can use [docker-compose](./compose.en.md) or [Kubernetes](./kubernetes.en.md) to run OpenBP on Raspberry Pi. 

!!! hint
    If you want to do this with Kubernetes, use **[K3S](https://k3s.io/)** instead of full Kubernetes installation. K3S uses only 500mb of RAM and runs perfectly on ARM. You don't need full Kubernetes installation.

## Preparing the OS
The only thing you need to do for Raspberry Pi is to make file system **ReadOnly** and use external storage for the data.

Normal SD cards connected to the RPI will stop working after several months of usage. Flash will last longer, but it is also not reliable enough for storing database data with intense I/O operations.