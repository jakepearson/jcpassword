package handler

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/jakepearson/jcpassword/encoder"
)

func hashGetHandler(webServer *WebServer, w http.ResponseWriter, r *http.Request) {
	hashID, e := strconv.Atoi(path.Base(r.URL.Path))
	if e != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	hash := webServer.Hashes[hashID]
	if hash == nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if !hash.Complete {
		w.WriteHeader(http.StatusProcessing)
	}
	fmt.Fprint(w, hash.Value)
}

func hashPostHandler(webServer *WebServer, w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Password missing (use password query parameter)")
		return
	}

	hash := Hash{
		ID:       len(webServer.Hashes),
		Complete: false,
	}
	webServer.Hashes[hash.ID] = &hash

	go func() {
		time.Sleep(time.Duration(webServer.SleepSeconds) * time.Second) //Slow down request to meet requirement
		hash.Value = encoder.HashAndEncode(password)
		hash.Complete = true
	}()

	fmt.Fprintf(w, "%d", hash.ID)
}

func hashHandler(webServer *WebServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			hashGetHandler(webServer, w, r)
		case "POST":
			hashPostHandler(webServer, w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
