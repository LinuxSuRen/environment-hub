# environment-hub

This project aims to provide a Kubernetes compatible environment quickly. Especially for the testing purpose.

## Get started
```shell
docker run -v /var/run/docker.sock:/var/run/docker.sock \
    -p 9090:8080 \
    ghcr.io/linuxsuren/environment-hub:master
```
