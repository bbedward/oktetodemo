package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNPodsAPI(t *testing.T) {
	// Setup pods
	SetUpk8ApiTests()

	req, err := http.NewRequest("GET", "/npods", nil)
	if err != nil {
		t.Fatal(err)
	}

	controller := OKtetoAPIController{K8sApi: &K8sMockAPI}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Npods)

	handler.ServeHTTP(rr, req)

	// Assert 200 response code
	assert.Equal(t, http.StatusOK, rr.Code)
	bodyInt, err := strconv.Atoi(rr.Body.String())
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, bodyInt)
}

func TestGetPodsAPI(t *testing.T) {
	// Setup pods
	SetUpk8ApiTests()

	req, err := http.NewRequest("GET", "/pods?sort=name&order=asc", nil)
	if err != nil {
		t.Fatal(err)
	}

	controller := OKtetoAPIController{K8sApi: &K8sMockAPI}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Pods)

	handler.ServeHTTP(rr, req)

	// Assert 200 response code
	assert.Equal(t, http.StatusOK, rr.Code)
	// Unmarshall
	var resp []PodResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)
	// Assertions
	assert.Equal(t, 3, len(resp))
	assert.Equal(t, "pod1", resp[0].Name)
	assert.Equal(t, "pod2", resp[1].Name)
	assert.Equal(t, "pod3", resp[2].Name)
	assert.Equal(t, "2022-02-01T00:00:00Z", resp[0].CreatedTS.Format(time.RFC3339))
	assert.Equal(t, "2022-02-02T00:00:00Z", resp[1].CreatedTS.Format(time.RFC3339))
	assert.Equal(t, "2022-02-03T00:00:00Z", resp[2].CreatedTS.Format(time.RFC3339))
	assert.Equal(t, 0, resp[0].Restarts)
	assert.Equal(t, 1, resp[1].Restarts)
	assert.Equal(t, 2, resp[2].Restarts)
}

func TestGetPodsAPIInvalidSort(t *testing.T) {
	// Setup pods
	SetUpk8ApiTests()

	// Invalid sort param
	req, err := http.NewRequest("GET", "/pods?sort=id", nil)
	if err != nil {
		t.Fatal(err)
	}

	controller := OKtetoAPIController{K8sApi: &K8sMockAPI}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Pods)

	handler.ServeHTTP(rr, req)

	// Assert 400 response code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	// Unmarshall
	var resp map[string]string
	fmt.Print(rr.Body.String())
	json.Unmarshal(rr.Body.Bytes(), &resp)
	// Assertions
	assert.Equal(t, `invalid sort option. Valid options are "name", "age", or "restarts"`, resp["message"])
}

func TestGetPodsAPIInvalidSortORder(t *testing.T) {
	// Setup pods
	SetUpk8ApiTests()

	// Test invalid sort order
	req, err := http.NewRequest("GET", "/pods?sort=name&order=sideways", nil)
	if err != nil {
		t.Fatal(err)
	}

	controller := OKtetoAPIController{K8sApi: &K8sMockAPI}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Pods)
	handler.ServeHTTP(rr, req)

	// Assert 400 response code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	// Unmarshall
	var resp map[string]string
	fmt.Print(rr.Body.String())
	json.Unmarshal(rr.Body.Bytes(), &resp)
	// Assertions
	assert.Equal(t, `invalid sort direction. Valid options are "asc" or "desc"`, resp["message"])
}
