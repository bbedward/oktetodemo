package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

var K8sMockAPI KubernetesAPI

func SetUpk8ApiTests() {
	// Setup
	K8sMockAPI = KubernetesAPI{MockClientSet: testclient.NewSimpleClientset()}

	// Create a few pods
	// Feb 1
	str := "2022-02-01T00:00:00.000Z"
	t1, _ := time.Parse(time.RFC3339, str)
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "pod1",
			Namespace:         "bbedward",
			CreationTimestamp: metav1.Time{t1},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "okteto:go",
					Image:           "okteto:go",
					ImagePullPolicy: "Always",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:         "okteto:go",
					RestartCount: 0,
				},
			},
		},
	}

	// Feb 2
	str = "2022-02-02T00:00:00.000Z"
	t2, _ := time.Parse(time.RFC3339, str)
	pod2 := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "pod2",
			Namespace:         "bbedward",
			CreationTimestamp: metav1.Time{t2},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "okteto:go",
					Image:           "okteto:go",
					ImagePullPolicy: "Always",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:         "okteto:go",
					RestartCount: 1,
				},
			},
		},
	}

	// Feb 3
	str = "2022-02-03T00:00:00.000Z"
	t3, _ := time.Parse(time.RFC3339, str)
	pod3 := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "pod3",
			Namespace:         "bbedward",
			CreationTimestamp: metav1.Time{t3},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "okteto:go",
					Image:           "okteto:go",
					ImagePullPolicy: "Always",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:         "okteto:go",
					RestartCount: 2,
				},
			},
		},
	}

	_, err := K8sMockAPI.MockClientSet.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error occured while creating pod %s: %s", pod.Name, err.Error())
		panic("Error in setUp creating pod")
	}
	_, err = K8sMockAPI.MockClientSet.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod2, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error occured while creating pod %s: %s", pod.Name, err.Error())
		panic("Error in setUp creating pod")
	}
	_, err = K8sMockAPI.MockClientSet.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod3, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error occured while creating pod %s: %s", pod.Name, err.Error())
		panic("Error in setUp creating pod")
	}
}

// Actual tests
func TestGetNPods(t *testing.T) {
	SetUpk8ApiTests()

	npods, err := K8sMockAPI.GetNPods("bbedward")
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, npods)
}

func TestGetPods(t *testing.T) {
	SetUpk8ApiTests()

	pods, err := K8sMockAPI.GetPods("bbedward")
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(pods))

	// Sort and check data
	K8sMockAPI.SortPods(pods, SortName, SortAscending)
	assert.Equal(t, pods[0].Name, "pod1")
	assert.Equal(t, pods[0].CreatedTS.Format(time.RFC3339), "2022-02-01T00:00:00Z")
	assert.Equal(t, pods[0].Restarts, 0)
	assert.Equal(t, pods[1].Name, "pod2")
	assert.Equal(t, pods[1].CreatedTS.Format(time.RFC3339), "2022-02-02T00:00:00Z")
	assert.Equal(t, pods[1].Restarts, 1)
	assert.Equal(t, pods[2].Name, "pod3")
	assert.Equal(t, pods[2].CreatedTS.Format(time.RFC3339), "2022-02-03T00:00:00Z")
	assert.Equal(t, pods[2].Restarts, 2)
}

func TestSortPods(t *testing.T) {
	SetUpk8ApiTests()

	pods, _ := K8sMockAPI.GetPods("bbedward")

	// Sort by name ascending
	K8sMockAPI.SortPods(pods, SortName, SortAscending)
	assert.Equal(t, pods[0].Name, "pod1")
	assert.Equal(t, pods[1].Name, "pod2")
	assert.Equal(t, pods[2].Name, "pod3")
	// Sort by name descending
	K8sMockAPI.SortPods(pods, SortName, SortDescending)
	assert.Equal(t, pods[0].Name, "pod3")
	assert.Equal(t, pods[1].Name, "pod2")
	assert.Equal(t, pods[2].Name, "pod1")
	// Sort by age ascending
	K8sMockAPI.SortPods(pods, SortAge, SortAscending)
	assert.Equal(t, pods[0].Name, "pod3")
	assert.Equal(t, pods[1].Name, "pod2")
	assert.Equal(t, pods[2].Name, "pod1")
	// Sort by age descending
	K8sMockAPI.SortPods(pods, SortAge, SortDescending)
	assert.Equal(t, pods[0].Name, "pod1")
	assert.Equal(t, pods[1].Name, "pod2")
	assert.Equal(t, pods[2].Name, "pod3")
	// Sort by restarts ascending
	K8sMockAPI.SortPods(pods, SortRestarts, SortAscending)
	assert.Equal(t, pods[0].Name, "pod1")
	assert.Equal(t, pods[1].Name, "pod2")
	assert.Equal(t, pods[2].Name, "pod3")
	// Sort by restarts descending
	K8sMockAPI.SortPods(pods, SortRestarts, SortDescending)
	assert.Equal(t, pods[0].Name, "pod3")
	assert.Equal(t, pods[1].Name, "pod2")
	assert.Equal(t, pods[2].Name, "pod1")
}

func TestFormatAgeString(t *testing.T) {
	// Our base date:
	str := "2022-02-01T00:00:00.000Z"
	t1, _ := time.Parse(time.RFC3339, str)

	// Fake now seconds to replace real time
	fakeNow := func() time.Time {
		str := "2022-02-01T00:00:30.000Z"
		t, _ := time.Parse(time.RFC3339, str)
		return t
	}
	assert.Equal(t, "30.00 seconds", FormatAgeString(t1, fakeNow))
	// 5 minutes
	fakeNow = func() time.Time {
		str := "2022-02-01T00:05:00.000Z"
		t, _ := time.Parse(time.RFC3339, str)
		return t
	}
	assert.Equal(t, "5.00 minutes", FormatAgeString(t1, fakeNow))
	// 3 hours
	fakeNow = func() time.Time {
		str := "2022-02-01T03:00:00.000Z"
		t, _ := time.Parse(time.RFC3339, str)
		return t
	}
	assert.Equal(t, "3.00 hours", FormatAgeString(t1, fakeNow))
	// 10 days
	fakeNow = func() time.Time {
		str := "2022-02-11T00:00:00.000Z"
		t, _ := time.Parse(time.RFC3339, str)
		return t
	}
	assert.Equal(t, "10 days", FormatAgeString(t1, fakeNow))
}
