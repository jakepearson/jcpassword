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

var server = CreateServer(8080, false, sleepSeconds)

func executeRequest(request *http.Request) *httptest.ResponseRecorder {
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
	password := "angryMonkey"
	response := executeHashRequest(&password)

	if response.Code != 200 {
		t.Errorf("Wrong code returned: %d", response.Code)
	}

	hashID := response.Body.String()
	hashURI := fmt.Sprintf("/hash/%s", hashID)

	time.Sleep(500 * time.Millisecond) // Sleep for half the time of the processing
	readHashRequest, _ := http.NewRequest("GET", hashURI, nil)
	readHashResponse := executeRequest(readHashRequest)

	if readHashResponse.Code != http.StatusProcessing {
		t.Errorf("Request should still be processing: %d", readHashResponse.Code)
	}

	time.Sleep(600 * time.Millisecond) // Sleep until processing is complete

	readHashRequest, _ = http.NewRequest("GET", hashURI, nil)
	readHashResponse = executeRequest(readHashRequest)

	if response.Code != http.StatusOK {
		t.Errorf("Wrong code returned: %d", response.Code)
	}

	body := readHashResponse.Body.String()
	expected := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	if body != expected {
		t.Errorf("Wrong hash returned: %s", body)
	}
}

func TestGetBadHashId(t *testing.T) {
	hashID := 1234
	hashURI := fmt.Sprintf("/hash/%d", hashID)

	readHashRequest, _ := http.NewRequest("GET", hashURI, nil)
	readHashResponse := executeRequest(readHashRequest)

	if readHashResponse.Code != http.StatusNotFound {
		t.Errorf("Wrong code returned: %d", readHashResponse.Code)
	}
}

func TestInvalidHashId(t *testing.T) {
	hashID := "invalidID"
	hashURI := fmt.Sprintf("/hash/%s", hashID)

	readHashRequest, _ := http.NewRequest("GET", hashURI, nil)
	readHashResponse := executeRequest(readHashRequest)

	if readHashResponse.Code != http.StatusNotFound {
		t.Errorf("Wrong code returned: %d", readHashResponse.Code)
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
	request, _ := http.NewRequest("PATCH", uri, nil)
	response := executeRequest(request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Wrong error code returned: %d", response.Code)
	}
}
