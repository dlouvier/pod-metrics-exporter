# Pod Metrics Exporter
It is an application written in Golang that once it is started, searchs
in Kubernetes for the number of pods matching the selected
labels and starts an HTTP server which exposes the number of pods as a Prometheus metric (`pod_count`).

The labels can be passed dinamically using flags.

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
./pod-metrics-exporter --label-name app \
                       --label-value demo-db \
                       --metricListenAddr "[OPTIONAL]" \
                       --kubeconfig "[OPTIONAL]"

```

`--label-name app` specified the label name to match. For this example the label *app*. This flag is mandatory.

`--label-value demo-db` specified the value for the specified label above to match. For this example the value *demo-db*. This flag is mandatory.

`--metricListenAddr host:port` specified the host:port where the HTTP Server will start. The port needs to be free. By default it is: *127.0.0.1:8080*. This flag is optional.

`--kubeconfig path` specifies the path to our kubeconfig file. by default is *~/.kube/config*. This flag is optional.
# Troubleshooting

- The Kubeconfig configuration is not correct
```
2021/07/04 20:06:45 It was not possible to connect to the Kubernetes cluster.
Check if it is possible to connect to the Kubernetes API.
invalid configuration: no configuration has been provided, try setting KUBERNETES_MASTER environment variable
```
Try with `kubectl get pods` the command responds correctly, if not, the kubectl config is not correct.

By default, the application tries to use `~/.kube/config` use a different one run the application using the following flag: `--kubeconfig "PATH"`


