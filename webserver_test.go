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

func executeHashRequest(password *string) *httptest.ResponseRecorder {
	handler := handler.Create()
	params := make(url.Values)
	if password != nil {
		params.Add("password", *password)
	}
	uri := fmt.Sprintf("%s?%s", "/hash", params.Encode())
	request, _ := http.NewRequest("GET", uri, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func TestHashRoute(t *testing.T) {
	startTime := time.Now()
	password := "angryMonkey"
	response := executeHashRequest(&password)

	finishTime := time.Now()
	if finishTime.Sub(startTime).Seconds() < 5 {
		t.Errorf("Request was faster than 5 seconds")
	}

	if response.Code != 200 {
		t.Fatalf("Wrong code returned: %d", response.Code)
	}

	body := response.Body.String()
	expected := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	if body != expected {
		t.Errorf("Wrong hash returned")
	}
}

func TestHashRouteMissingParameter(t *testing.T) {
	response := executeHashRequest(nil)
	if response.Code != 400 {
		t.Fatalf("Wrong code returned: %d", response.Code)
	}

	body := response.Body.String()
	expected := "Password missing (use password query parameter"
	if body != expected {
		t.Errorf("Wrong hash returned")
	}
}

func TestHashRouteBlankParameter(t *testing.T) {
	password := ""
	response := executeHashRequest(&password)
	if response.Code != 400 {
		t.Fatalf("Wrong code returned: %d", response.Code)
	}

	body := response.Body.String()
	expected := "Password missing (use password query parameter"
	if body != expected {
		t.Errorf("Wrong hash returned")
	}
}
