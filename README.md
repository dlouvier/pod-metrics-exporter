# Pod Metrics Exporter
This application checks periodically (every 10s) how many
Pods matching the selected label name/value exist in the the Kubernetes cluster (all namespaces).

This number is later exposed as Prometheus metric (`pod_count`).

# Requirements
- Golang 1.16
- Kubernetes Cluster (>1.20)

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

Once the application starts, it will let know us in which URL we can check the metrics:

Output:
```
$ cd cmd/pod-metrics-exporter
$ ./pod-metrics-exporter --label-name app \
                       --label-value demo-db
2021/07/04 21:20:46 Loading kubeconfig from /home/dani/.kube/config
2021/07/04 21:20:46 Starting webserver in 127.0.0.1:8080
2021/07/04 21:20:46 Metrics are available in http://127.0.0.1:8080/metrics
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


