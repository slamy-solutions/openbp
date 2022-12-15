# Deploy for production using docker-compose
This is recommended way to go with small installations that don't require high availability. Some of the key features of this type of deployment are:

- Easy to deploy/maintain/update
- Backup is just copying a folder with data.
- You will need to stop the deployment for 30-60 seconds on every update.
- Not scalable. Everything is on one server.
- If a machine (VPS or physical) is down - the entire deployment is down.
- It doesn't have data redundancy - if you do not backup regularly, you can lose all your data.

## Deploy
Before starting, please make sure your system has minimal requirements:
   
| Requirement | Minimal                                    | Recommended                 |
| ----------- | ------------------------------------------ | --------------------------- |
| OS          | Linux/Windows                              | Linux                       |
| Arch        | linux/amd64, linux/arm64/v8, windows/amd64 | linux/amd64, linux/arm64/v8 |
| CPU         | 1 vCPU                                     | 4 vCPU                      |
| RAM         | 1GB (more with Windows)                    | 4GB                         |

!!! info
    We suggest using Linux operating system. Windows installation may use much more resources and not be stable (as far as Docker for Windows is not stable).

1. Install [Docker](https://docs.docker.com/get-docker/) and [docker-compose](https://docs.docker.com/compose/install/)
2. Get the `docker-compose.all-in-one.yml` file from the [repository](https://github.com/slamy-solutions/openbp/blob/master/scripts/compose/docker-compose.all-in-one.yml)
3. Start OpenBP using `docker-compose -f docker-compose.all-in-one.yml up -d` command.

??? example "All-in-one script for Linux"
    ```sh
    echo "Installing docker"
    curl -s https://get.docker.io | bash -s
    sudo groupadd docker
    sudo usermod -aG docker $USER
    newgrp docker

    echo "Installing docker-compose using python3"
    python3 -m pip install docker-compose
     
    echo "Getting docker-compose.yml file"
    mkdir openbp
    cd openbp
    curl -s https://github.com/slamy-solutions/openbp/blob/master/scripts/compose/docker-compose.all-in-one.yml >> docker-compose.all-in-one.yml

    echo "Starting OpenBP"
    docker-compose -f docker-compose.all-in-one.yml up -d
    ```

After startup, OpenBP will create folder `data` near the `docker-compose.all-in-one.yml` file. This is folder, where OpenBP will persist all its data.

## Configuring
### Selecting the version
If you don't have a version selected, the OpenBP will run with the latest possible. If you want a specific version, first of all, download it from the branch with `tag`. Use the environment variable or `docker-compose` [.env](https://docs.docker.com/compose/environment-variables/) file to set a value for `OPENBP_VERSION` with the version you want.


## Update
The system will not be available for ~20-30 seconds during the update. If you want high availability, consider using [Kubernetes Deployment](./kubernetes.en.md). 

First of all, you have to update `docker-compose.*.yml` files. Go to the [repository](https://github.com/slamy-solutions/openbp/tree/master/scripts/compose), select `tag` with the version you want and download files with the newer version.

Set the environment variable with the selected version:
```sh
export OPENBP_VERSION=v1.0.13
```
!!! hint
    You can also put the version into the `docker-compose` [.env](https://docs.docker.com/compose/environment-variables/) file.

Update deployment by pulling newer versions of the `docker` containers.
```sh
docker-compose -f docker-compose.all-in-one.yml pull
docker-compose -f docker-compose.all-in-one.yml down
docker-compose -f docker-compose.all-in-one.yml up -d
```

## Partial deployment
You can find more `docker-compose` files in the compose scripts [repository](https://github.com/slamy-solutions/openbp/tree/master/scripts/compose).

If you just want to install all the modules and use all features, use `docker-compose.all-in-one.yml`, but if you don't need all the functionalities, you can combine modules using standard `docker-compose` syntax. 

!!! example
    `docker-compose up -f docker-compose.yml -f docker-compose.system.yml -f docker-compose.native.yml up` will start up only the `system` and `native` modules.

Table with specifications for every file you can find below:

| File                          | Description                             | Dependencies                                                             |
| ----------------------------- | --------------------------------------- | ------------------------------------------------------------------------ |
| docker-compose.all-in-one.yml | Run all the modules and features        | -                                                                        |
| docker-compose.yml            | Create internal network and data volume | -                                                                        |
| docker-compose.system.yml     | Run only `system` module                | docker-compose.yml                                                       |
| docker-compose.native.yml     | Run only `native` module                | docker-compose.yml, docker-compose.system.yml                            |
| docker-compose.tools.yml      | Run only `tools` module                 | docker-compose.yml, docker-compose.system.yml, docker-compose.native.yml |

## Backup
All the data of the application is located in the `data` folder (near the `docker-compose.*.yml` files). You can back up data by doing a copy of the folder.