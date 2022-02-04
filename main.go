package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	fmt.Println("Starting hello-world server...")
	// Create kubernetes client config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	ClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	// Create controller
	controller := OKtetoAPIController{K8sApi: &KubernetesAPI{ClientSet: ClientSet}}
	// Number of pods
	http.HandleFunc("/npods", controller.Npods)
	// Pods with sort function
	http.HandleFunc("/pods", controller.Pods)
	// Prometheus metrics
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "bbedward",
		Name:      "pod_count",
		Help:      "Number of pods in the bbedward namespace",
	}, func() float64 {
		npods, err := controller.K8sApi.GetNPods("bbedward")
		if err != nil {
			return -1
		}
		return float64(npods)
	})
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// Define a "controller" so k8sAPI can be available
type OKtetoAPIController struct {
	K8sApi *KubernetesAPI
}

// Return number of pods in bbedward namespace
func (controller *OKtetoAPIController) Npods(w http.ResponseWriter, r *http.Request) {
	npods, err := controller.K8sApi.GetNPods("bbedward")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	fmt.Fprintf(w, "%d", npods)
}

// Return pods with sort functionality
func (controller *OKtetoAPIController) Pods(w http.ResponseWriter, r *http.Request) {
	podResp, err := controller.K8sApi.GetPods("bbedward")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
		return
	}
	// Get sort type if present
	sort := strings.ToLower(r.URL.Query().Get("sort"))
	if sort != "" && sort != string(SortName) && sort != string(SortAge) && sort != string(SortRestarts) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "invalid sort option. Valid options are \"%s\", \"%s\", or \"%s\""}`, SortName, SortAge, SortRestarts)))
		return
	} else if sort != "" {
		sortDirection := strings.ToLower(r.URL.Query().Get("order"))
		if sortDirection != "" && sortDirection != string(SortDescending) && sortDirection != string(SortAscending) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(`{"message": "invalid sort direction. Valid options are \"%s\" or \"%s\""}`, SortAscending, SortDescending)))
			return
		} else if sortDirection == "" {
			// Ascending by default
			sortDirection = string(SortAscending)
		}
		// Perform sort
		controller.K8sApi.SortPods(podResp, PodSortMethod(sort), PodSortDirection(sortDirection))
	}
	// Serialize and return response
	marshalled, err := json.Marshal(podResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
		return
	}
	w.Write(marshalled)
}
