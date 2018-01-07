package handler

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/jakepearson/jcpassword/encoder"
)

//hasGetHandler accepts a request to convert an id into a completed hash
//if the id is not in the "database", it will return `StatusNotFound`
//if the id is in the "database", but the hash isn't ready the handler will return `StatusProcessing`
//if the id is in the "database", and the hash is ready, the hash will be returned
func hashGetHandler(webServer *WebServer, w http.ResponseWriter, r *http.Request) {
	hashID, e := strconv.Atoi(path.Base(r.URL.Path))
	if e != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	hash := webServer.Hashes[hashID]
	if hash == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	if !hash.Complete {
		w.WriteHeader(http.StatusProcessing)
		fmt.Fprintf(w, "Processing in progress")
		return
	}
	fmt.Fprint(w, hash.Value)
}

//hashPostHandler accepts a request to process a password
//It kicks off the background job to generate the hash and returns an id that can be used to lookup
//the hash at a future time
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
