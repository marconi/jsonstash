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
	blh.RestHandler.ListJSONResponse(w, r, names)
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
	bkey := r.URL.Query().Get(":bkey")

	// TODO: implement querying by range

	bucket, err := bh.Stash.Get(bkey)
	if err != nil {
		log.Println("Invalid bucket key: ", bkey)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bh.RestHandler.ListJSONResponse(w, r, bucket.GetAll())
}

func (bh *BucketHandler) Post(w http.ResponseWriter, r *http.Request) {
	bkey := r.URL.Query().Get(":bkey")
	bucket, err := bh.Stash.Get(bkey)
	if err != nil {
		log.Println("Invalid bucket key: ", bkey)
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
	if err := bucket.Add(payload.Key, payload.Value); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bh.RestHandler.Response(w, r, "Payload added to bucket.")
}

func (bh *BucketHandler) Delete(w http.ResponseWriter, r *http.Request) {
	bkey := r.URL.Query().Get(":bkey")
	if err := bh.Stash.Delete(bkey); err != nil {
		log.Println("Unable to delete bucket: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bh.RestHandler.Response(w, r, "Bucket deleted successfully.")
}

type ValueHandler struct {
	*BaseHandler
}

func NewValueHandler(rh *RestHandler) *ValueHandler {
	handler := &ValueHandler{
		BaseHandler: &BaseHandler{
			RestHandler: rh,
			Stash:       rh.Stash,
		},
	}
	return handler
}

func (vh *ValueHandler) Get(w http.ResponseWriter, r *http.Request) {
	bkey := r.URL.Query().Get(":bkey")
	vkey := r.URL.Query().Get(":vkey")
	bucket, err := vh.Stash.Get(bkey)
	if err != nil {
		log.Println("Invalid bucket key: ", bkey)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	value, err := bucket.Get(vkey)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vh.RestHandler.JSONResponse(w, r, value)
}

func (vh *ValueHandler) Put(w http.ResponseWriter, r *http.Request) {
	bkey := r.URL.Query().Get(":bkey")
	vkey := r.URL.Query().Get(":vkey")

	bucket, err := vh.Stash.Get(bkey)
	if err != nil {
		log.Println("Invalid bucket key: ", bkey)
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
	if payload.Value == "" {
		errStr := "Invalid payload."
		log.Println(errStr)
		err = errors.New(errStr)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := bucket.Update(vkey, payload.Value); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vh.RestHandler.JSONResponse(w, r, "Value has been updated successfully.")
}

func (vh *ValueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	bkey := r.URL.Query().Get(":bkey")
	vkey := r.URL.Query().Get(":vkey")
	bucket, err := vh.Stash.Get(bkey)
	if err != nil {
		log.Println("Invalid bucket key: ", bkey)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := bucket.Delete(vkey); err != nil {
		log.Println("Unable to delete: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vh.RestHandler.JSONResponse(w, r, "Value has been deleted successfully.")
}
