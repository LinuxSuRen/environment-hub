# environment-hub

## Get started
```shell
docker run -v /var/run/docker.sock:/var/run/docker.sock \
    -e ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.21 \
    -p 9090:8080 \
    surenpi/environment-hub:master
```
