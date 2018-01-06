package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

const sleepSeconds = 1

func executeRequest(request *http.Request) *httptest.ResponseRecorder {
	server := CreateServer(8080, false, sleepSeconds)
	response := httptest.NewRecorder()
	handler := *server.Handler
	handler.ServeHTTP(response, request)
	return response
}

func executeHashRequest(password *string) *httptest.ResponseRecorder {
	params := make(url.Values)
	if password != nil {
		params.Add("password", *password)
	}
	uri := fmt.Sprintf("%s?%s", "/hash", params.Encode())
	request, _ := http.NewRequest("POST", uri, nil)
	return executeRequest(request)
}

func TestHashRoute(t *testing.T) {
	startTime := time.Now()
	password := "angryMonkey"
	response := executeHashRequest(&password)

	finishTime := time.Now()
	if finishTime.Sub(startTime).Seconds() < sleepSeconds {
		t.Errorf("Request was faster than %d seconds", sleepSeconds)
	}

	if response.Code != 200 {
		t.Errorf("Wrong code returned: %d", response.Code)
	}

	body := response.Body.String()
	expected := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	if body != expected {
		t.Errorf("Wrong hash returned: %s", body)
	}
}

func TestHashRouteMissingParameter(t *testing.T) {
	response := executeHashRequest(nil)
	if response.Code != 400 {
		t.Errorf("Wrong code returned: %d", response.Code)
	}

	body := response.Body.String()
	expected := "Password missing (use password query parameter)"
	if body != expected {
		t.Errorf("Wrong error returned: %s", body)
	}
}

func TestHashRouteBlankParameter(t *testing.T) {
	password := ""
	response := executeHashRequest(&password)
	if response.Code != 400 {
		t.Errorf("Wrong code returned: %d", response.Code)
	}

	body := response.Body.String()
	expected := "Password missing (use password query parameter)"
	if body != expected {
		t.Errorf("Wrong error returned: %s", body)
	}
}

func TestWrongMethod(t *testing.T) {
	params := make(url.Values)
	params.Add("password", "test")
	uri := fmt.Sprintf("%s?%s", "/hash", params.Encode())
	request, _ := http.NewRequest("GET", uri, nil)
	response := executeRequest(request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Wrong error code returned: %d", response.Code)
	}
}
