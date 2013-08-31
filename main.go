package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/marconi/jsonstash/bucket"
	"github.com/marconi/jsonstash/utils"
)

var stash *bucket.Stash

type BucketPayload struct {
	Key string
}

func WriteResponse(w http.ResponseWriter, r *http.Request, data string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body: ", err)
			err = errors.New("Invalid posted payload")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		payload := new(BucketPayload)
		if err := json.Unmarshal(b, payload); err != nil {
			log.Println("Error parsing body: ", err)
			err = errors.New("Invalid posted payload")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if payload.Key == "" {
			errStr := "Empty bucket key."
			log.Println(errStr)
			err = errors.New(errStr)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// we have a valid payload
		stash.Add(payload.Key)
		msg := fmt.Sprintf("Bucket %s has been added.", payload.Key)
		WriteResponse(w, r, msg)
		return
	}

	// GET request
	names := utils.ToJSON(stash.GetBucketNames())
	WriteResponse(w, r, names)
}

func main() {
	stash = bucket.NewStash()

	router := mux.NewRouter()
	router.HandleFunc("/buckets", ListBuckets).Methods("GET", "POST")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
