package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestShutdownHandler(t *testing.T) {
	server := CreateServer(8080, sleepSeconds)
	request, _ := http.NewRequest("GET", "/shutdown", nil)
	response := httptest.NewRecorder()
	handler := *server.Handler

	handler.ServeHTTP(response, request)
	time.Sleep(500 * time.Millisecond)

	if server.Closed {
		t.Errorf("Shutdown should not occur before the time has elapsed")
	}

	time.Sleep(600 * time.Millisecond)
	if !server.Closed {
		t.Errorf("Shutdown should have occurred by now")
	}
}
