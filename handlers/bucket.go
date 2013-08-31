package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type BucketPayload struct {
	Key string
}

type ValuePayload struct {
	Key   string
	Value string
}

type BucketListHandler struct {
	*BaseHandler
}

func NewBucketListHandler(rh *RestHandler) *BucketListHandler {
	handler := &BucketListHandler{
		BaseHandler: &BaseHandler{
			RestHandler: rh,
			Stash:       rh.Stash,
		},
	}
	return handler
}

func (blh *BucketListHandler) Get(w http.ResponseWriter, r *http.Request) {
	names := blh.Stash.GetBucketNames()
	blh.RestHandler.JSONResponse(w, r, names)
}

func (blh *BucketListHandler) Post(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body: ", err)
		err = errors.New("Invalid payload.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload := new(BucketPayload)
	if err := json.Unmarshal(b, payload); err != nil {
		log.Println("Error parsing body: ", err)
		err = errors.New("Invalid payload.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if payload.Key == "" {
		errStr := "Empty bucket key."
		log.Println(errStr)
		err = errors.New(errStr)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// we have a valid payload
	if _, err := blh.Stash.Add(payload.Key); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := fmt.Sprintf("Bucket %s has been added.", payload.Key)
	blh.RestHandler.Response(w, r, msg)
}

type BucketHandler struct {
	*BaseHandler
}

func NewBucketHandler(rh *RestHandler) *BucketHandler {
	handler := &BucketHandler{
		BaseHandler: &BaseHandler{
			RestHandler: rh,
			Stash:       rh.Stash,
		},
	}
	return handler
}

func (bh *BucketHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(":key")
	bucket, err := bh.Stash.Get(key)
	if err != nil {
		log.Println("Invalid bucket key: ", key)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bh.RestHandler.JSONResponse(w, r, bucket.GetAll())
}

func (bh *BucketHandler) Post(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(":key")
	bucket, err := bh.Stash.Get(key)
	if err != nil {
		log.Println("Invalid bucket key: ", key)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body: ", err)
		err = errors.New("Invalid payload.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload := new(ValuePayload)
	if err := json.Unmarshal(b, payload); err != nil {
		log.Println("Error parsing body: ", err)
		err = errors.New("Invalid payload.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if payload.Key == "" || payload.Value == "" {
		errStr := "Invalid payload."
		log.Println(errStr)
		err = errors.New(errStr)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// we have a valid payload
	bucket.Add(payload.Key, payload.Value)
	bh.RestHandler.Response(w, r, "Payload added to bucket.")
}
