# Pod Metrics Exporter
It is a binary written in Golang that once started, it searchs
in Kubernetes for the number of pods matching the selected
labels and starts an HTTP server which exposes the number of pods as a Prometheus metric (`pod_count`).

# Requirements
Golang 1.16


## How to run the unit tests
``` 
make test
```

## How to build the binary locally
```
make build
```

## How to run the application
```
cd cmd/pod-metrics-exporter
./pod-metrics-exporter --label-name app --label-value demo-db
```


# Troubleshooting

- The Kubeconfig configuration is not correct
```
2021/07/04 20:06:45 It was not possible to connect to the Kubernetes cluster.
Check if it is possible to connect to the Kubernetes API.
invalid configuration: no configuration has been provided, try setting KUBERNETES_MASTER environment variable
```
Try with `kubectl get pods` the command responds correctly, if not, the kubectl config is not correct.

By default, the application tries to use `~/.kube/config` use a different one run the application using the following flag: `--kubeconfig "PATH"`


