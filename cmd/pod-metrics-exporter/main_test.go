package main

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeClient "k8s.io/client-go/kubernetes/fake"
)

func TestGetPods(t *testing.T) {

	api := &k8sApi{
		Client: fakeClient.NewSimpleClientset(),
	}

	// Define the fake Pod object
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Labels: map[string]string{
				"app": "demo",
			},
		},
	}

	// Add the fake Pod object to the fake mocked API.
	api.Client.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})

	selector := createListOptions("app", "demo", "Running")
	pods, _ := api.getPods(selector)

	if pods[0].Name != "test-pod" {
		t.Errorf("test-pod is not equal to %s", pods[0].Name)
	}

}
