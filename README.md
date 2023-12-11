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
export SERVER=http://172.11.0.6:9090
curl -XPOST '$SERVER/v1/k3d/clusters' -d '{"servers":1,"agents":1,"port":30000,"name":"hello"}'
curl '$SERVER/v1/k3d/clusters/hello/kubeconfig' -o test.yaml
helm --kubeconfig=test.yaml --kube-insecure-skip-tls-verify=true install atest oci://docker.io/linuxsuren/api-testing \
    --version v0.0.2-helm \
    --set service.nodePort=30000 \
    --set service.type=NodePort \
    --set image.tag=v0.0.15 \
    --set persistence.enabled=false

kubectl --kubeconfig=test.yaml --insecure-skip-tls-verify=true get pod -A

curl 'http://172.11.0.6:9090/v1/k3d/clusters/hello/portbinding?port=30000'
```
