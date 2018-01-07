package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//statsMiddleware captures the start and end time of each request and updates the `Statistics` object
func statsMiddleware(webServer *WebServer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		responseTime := time.Now().Sub(startTime)
		webServer.Statistics.TotalRequests++
		webServer.Statistics.totalResponseTime += responseTime
		totalResponseMS := float32(webServer.Statistics.totalResponseTime / time.Millisecond)
		webServer.Statistics.AverageResponseTimeMS = totalResponseMS / float32(webServer.Statistics.TotalRequests)
	})
}

//statsHandler will return a `json` object of that statistics collected by the stats middleware`
func statsHandler(webServer *WebServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, e := json.Marshal(webServer.Statistics)
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", e)
			return
		}
		fmt.Fprintf(w, "%s", data)
	}
}
