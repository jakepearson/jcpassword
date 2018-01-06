package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func executeIsolatedRequest(webServer *WebServer, uri string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest("POST", uri, nil)
	response := httptest.NewRecorder()
	handler := *webServer.Handler
	handler.ServeHTTP(response, request)
	return response
}

func TestStatsEndpoint(t *testing.T) {
	var server = CreateServer(8080, false, sleepSeconds)
	for i := 0; i < 10; i++ {
		executeIsolatedRequest(server, "/hash?password=test")
	}

	response := executeIsolatedRequest(server, "/stats")

	if response.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d", response.Code)
	}

	responseData := make(map[string]interface{})

	if err := json.Unmarshal(response.Body.Bytes(), &responseData); err != nil {
		t.Errorf("Invalid JSON: %v", err)
	}
	totalRequests := responseData["TotalRequests"]
	fmt.Printf("Stats: %v\n", responseData)
	if totalRequests != 10.0 {
		t.Errorf("Total Response Count Invalid: %d", totalRequests)
	}
}
