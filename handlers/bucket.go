package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/marconi/jsonstash/bucket"
)

type BucketPayload struct {
	Key string
}

type BucketView struct {
	RESTView *RESTView
	Stash    *bucket.Stash
}

func (bv *BucketView) Get(w http.ResponseWriter, r *http.Request) {
	names := bv.Stash.GetBucketNames()
	bv.RESTView.JSONResponse(w, r, names)
}

func (bv *BucketView) Post(w http.ResponseWriter, r *http.Request) {
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
	bv.Stash.Add(payload.Key)
	msg := fmt.Sprintf("Bucket %s has been added.", payload.Key)
	bv.RESTView.Response(w, r, msg)
}
