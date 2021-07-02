package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	labelName        = flag.String("label-name", "", "Set the label name to filter")
	labelValue       = flag.String("label-value", "", "Set the value to search from the set label")
	metricListenAddr = flag.String("metrics-listen-addr", "127.0.0.1:8080", "Address to listen blabla")
	kubeconfig       = flag.String("kubeconfig", "/home/dani/.kube/config", "Path to Kubeconfig")
)

var podCountMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "pod_count",
		Help: "Number of pods.",
	},
	[]string{"label_name", "label_value", "phase"},
)

// k8sclient tries to figure out wherever the application is running
// inside of a Kubernetes cluster or outside and depending of this it
// uses a different configuration to start.
//
// If using Kubernetes, it tries to use the assigned service account
// (it has to be created previously) if not, it will read the kubeconfig
// file from user home directory or provided as a argument.
func k8sClient() *kubernetes.Clientset {
	_, runningInK8s := os.LookupEnv("KUBERNETES_SERVICE_HOST")

	if runningInK8s {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		return clientset
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		return clientset
	}
}

// getPods returns an array with the list of Pods that match the criteria.
func getPods(listOptions metav1.ListOptions, namespace string, clientset *kubernetes.Clientset) ([]v1.Pod, error) {
	podList, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	return podList.Items, err
}

// createListOptions parse the argument that the user specifies for
// filtering for the pods to list only the pods that match the query.
func createListOptions(labelName string, labelValue string, phase string) metav1.ListOptions {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{labelName: labelValue}}
	fieldSelector := fmt.Sprintf("status.phase=%s", phase)

	return metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
		FieldSelector: fieldSelector,
	}
}

func startWs() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*metricListenAddr, nil))
}

func init() {
	flag.Parse()
	prometheus.MustRegister(podCountMetric)
}

func main() {
	client := k8sClient()
	go startWs()

	for {
		// https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase
		podPhases := [5]string{"Pending", "Running", "Succeeded", "Failed", "Unknown"}

		for _, phase := range podPhases {
			selector := createListOptions(*labelName, *labelValue, phase)
			pods, err := getPods(selector, "demo-production", client)
			if err != nil {
				fmt.Printf("Error getting the pod list. {}", err)
			}
			podCountMetric.WithLabelValues(*labelName, *labelValue, phase).Set(float64(len(pods)))
		}

		time.Sleep(10 * time.Second)
	}
}
