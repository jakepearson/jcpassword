package main

import "testing"
import "context"
import "net/http"
import "fmt"

func TestEnsureServerStarts(t *testing.T) {

	webServer := startServer()
	defer webServer.Server.Shutdown(context.Background())

	url := fmt.Sprintf("http://localhost:%d/stats", port())
	response, e := http.Get(url)

	if e != nil {
		t.Errorf("Request failed: %v", e)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code returned: %d", response.StatusCode)
	}
}
