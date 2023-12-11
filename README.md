# environment-hub

This project aims to provide a Kubernetes compatible environment quickly. Especially for the testing purpose.

## Get started
```shell
docker run -v /var/run/docker.sock:/var/run/docker.sock \
    -p 9090:8080 \
    ghcr.io/linuxsuren/environment-hub:master
```

## Use Case 1

```shell
curl 'http://172.11.0.6:9090/v1/k3d/clusters' \
  --data-raw '{"servers":1,"agents":1,"port":30000,"name":"1234"}'

curl 'http://172.11.0.6:9090/v1/k3d/clusters/1234/portbinding?port=30000'
curl 'http://172.11.0.6:9090/v1/k3d/clusters/1234/kubeconfig' -o test.yaml

kubectl --kubeconfig=test.yaml get ns
```
