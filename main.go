package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var ClientSet *kubernetes.Clientset
var k8sAPI KubernetesAPI

func main() {
	fmt.Println("Starting hello-world server...")
	// Create kubernetes client config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	k8sAPI = KubernetesAPI{ClientSet: ClientSet}
	// Number of pods
	http.HandleFunc("/npods", npods)
	// Pods with sort function
	http.HandleFunc("/pods", pods)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// Return number of pods in bbedward namespace
func npods(w http.ResponseWriter, r *http.Request) {
	npods, err := k8sAPI.GetNPods("bbedward")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	fmt.Fprintf(w, "%d", npods)
}

// Return pods with sort functionality
func pods(w http.ResponseWriter, r *http.Request) {
	podResp, err := k8sAPI.GetPods("bbedward")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
		return
	}
	// Get sort type if present
	sort := strings.ToLower(r.URL.Query().Get("sort"))
	if sort != "" && sort != string(SortName) && sort != string(SortAge) && sort != string(SortRestarts) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "invalid sort option. Valid options are "%s", "%s", or "%s""}`, SortName, SortAge, SortRestarts)))
		return
	} else if sort != "" {
		sortDirection := strings.ToLower(r.URL.Query().Get("order"))
		if sortDirection != "" && sortDirection != string(SortDescending) && sortDirection != string(SortAscending) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(`{"message": "invalid sort direction. Valid options are "%s" or "%s""}`, SortAscending, SortDescending)))
			return
		} else if sortDirection == "" {
			// Ascending by default
			sortDirection = string(SortAscending)
		}
		// Perform sort
		k8sAPI.SortPods(podResp, PodSortMethod(sort), PodSortDirection(sortDirection))
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
