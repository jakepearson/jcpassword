package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	hashReference, hasKey := webServer.Hashes.Load(hashID)
	if !hasKey {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	hash, isHash := hashReference.(*Hash)
	if !isHash {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong type stored in hashes")
	}
	if !hash.Complete {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Processing in progress")
		return
	}
	fmt.Fprint(w, hash.Value)
}

//getNextID will grab the idLock, increment the id, release the lock then return the new id
func getNextID(webServer *WebServer) int {
	webServer.idLock.Lock()
	webServer.nextID++
	result := webServer.nextID
	webServer.idLock.Unlock()
	return result
}

//getPasswordInput will read the body from the request and parse it as a query string
func getPasswordInput(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", errors.New("No body defined")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	parameters, err := url.ParseQuery(string(body))
	if err != nil {
		return "", err
	}
	password := parameters.Get("password")
	if password == "" {
		return "", errors.New("Password not set")
	}
	return password, nil
}

//hashPostHandler accepts a request to process a password
//It kicks off the background job to generate the hash and returns an id that can be used to lookup
//the hash at a future time
func hashPostHandler(webServer *WebServer, w http.ResponseWriter, r *http.Request) {
	password, error := getPasswordInput(r)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad post body: %v", error)
		return
	}

	hash := Hash{
		ID:       getNextID(webServer),
		Complete: false,
	}
	webServer.Hashes.Store(hash.ID, &hash)

	go func() {
		time.Sleep(time.Duration(webServer.SleepSeconds) * time.Second) //Slow down request to meet requirement
		hash.Value = encoder.HashAndEncode(password)
		hash.Complete = true
	}()

	fmt.Fprintf(w, "%d", hash.ID)
}

func hashHandler(webServer *WebServer) func(w http.ResponseWriter, r *http.Request) {
	return statsMiddleware(webServer, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			hashGetHandler(webServer, w, r)
		case "POST":
			hashPostHandler(webServer, w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
