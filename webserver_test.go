package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/jakepearson/jcpassword/handler"
)

func createHashRequest(password string) *http.Request {
	params := make(url.Values)
	params.Add("password", "angryMonkey")
	uri := fmt.Sprintf("%s?%s", "/hash", params.Encode())
	request, _ := http.NewRequest("GET", uri, nil)
	return request
}

func TestHashAPI(t *testing.T) {
	handler := handler.Create()
	request := createHashRequest("angryMonkey")

	response := httptest.NewRecorder()
	startTime := time.Now()
	handler.ServeHTTP(response, request)

	finishTime := time.Now()
	if finishTime.Sub(startTime).Seconds() < 5 {
		t.Errorf("Request was faster than 5 seconds")
	}

	expected := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	if response.Code != 200 {
		t.Fatalf("Wrong code returned: %d", response.Code)
	}
	body := response.Body.String()
	if body != expected {
		t.Errorf("Wrong hash returned")
	}
}
