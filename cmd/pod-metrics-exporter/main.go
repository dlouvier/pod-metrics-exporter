package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// Prometheus dependencies
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// Kubernetes API dependencies
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	labelName        = flag.String("label-name", "", "Set the label name to filter")
	labelValue       = flag.String("label-value", "", "Set the value to search from the set label")
	metricListenAddr = flag.String("metrics-listen-addr", "127.0.0.1:8080", "Address to listen blabla")
	kubeconfig       = flag.String("kubeconfig", "", "Path to Kubeconfig")
	api              = k8sApi{Client: nil}
	kubeConfigPath   string
)

var podCountMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "pod_count",
		Help: "Number of pods.",
	},
	[]string{"label_name", "label_value", "phase"},
)

// k8s implements kubernetes interface, so it is easier
// to create mock the API during testing.
type k8sApi struct {
	Client kubernetes.Interface
}

// k8sclient tries to figure out wherever the application is running
// inside of a Kubernetes cluster or outside and depending of this it
// uses a different configuration to start.
//
// If using Kubernetes, it tries to use the assigned service account
// (it has to be created previously) if not, it will read the kubeconfig
// file from user home directory or provided as a argument.
func k8sClient() *kubernetes.Clientset {
	if *kubeconfig == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Println("Error trying to discover user home directory.", err)
		}
		kubeConfigPath = homeDir + "/.kube/config"
	} else {
		kubeConfigPath = *kubeconfig
	}

	log.Println("Loading kubeconfig from", kubeConfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatal("It was not possible to connect to the Kubernetes cluster.\nCheck if it is possible to connect to the Kubernetes API.\n",
			err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("It was not possible to create the Kubernetes configuration.\n", err.Error())
	}
	return clientset
}

// getPods returns an array with the list of Pods that match the criteria.
func (k k8sApi) getPods(listOptions metav1.ListOptions, namespace string) ([]v1.Pod, error) {
	podList, err := k.Client.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	return podList.Items, err
}

// createListOptions create the list option necessary to filter for the pods
// that match the query.
func createListOptions(labelName string, labelValue string, phase string) metav1.ListOptions {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{labelName: labelValue}}
	fieldSelector := fmt.Sprintf("status.phase=%s", phase)

	return metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
		FieldSelector: fieldSelector,
	}
}

// startWs starts the http server necessary to export the prometheus metrics.
func startWs() {
	log.Println("Starting webserver in", *metricListenAddr)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*metricListenAddr, nil))
}

// checkFlags verifies the minimum the app has started with the minimum required
// flags
func checkFlags() {
	if *labelName == "" {
		log.Fatal("Flag --label-name is missing.")
	} else if *labelValue == "" {
		log.Fatal("Flag --label-value is missing.")
	}
}

func init() {
	prometheus.MustRegister(podCountMetric)
}

func main() {
	flag.Parse()
	checkFlags()

	api := &k8sApi{
		Client: k8sClient(),
	}

	go startWs()

	for {
		// https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase
		podPhases := [5]string{"Pending", "Running", "Succeeded", "Failed", "Unknown"}

		for _, phase := range podPhases {
			selector := createListOptions(*labelName, *labelValue, phase)
			pods, err := api.getPods(selector, "demo-production")
			if err != nil {
				log.Println("Error getting the pod list.\n", err)
			}

			podCountMetric.WithLabelValues(*labelName, *labelValue, phase).Set(float64(len(pods)))
		}

		time.Sleep(10 * time.Second)
	}
}
