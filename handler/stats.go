package handler

import "net/http"
import "fmt"
import "encoding/json"

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
