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

type RESTView struct {
	Router *mux.Router
	Stash  *bucket.Stash
}

func NewRestView() *mux.Router {
	rv := &RESTView{
		Router: mux.NewRouter(),
		Stash:  bucket.NewStash(),
	}

	rv.Router.HandleFunc("/buckets", rv.GetListBuckets).Methods("GET")
	rv.Router.HandleFunc("/buckets", rv.PostListBuckets).Methods("POST")

	return rv.Router
}

func (rv *RESTView) GetListBuckets(w http.ResponseWriter, r *http.Request) {
	names := rv.Stash.GetBucketNames()
	rv.JSONResponse(w, r, names)
}

func (rv *RESTView) PostListBuckets(w http.ResponseWriter, r *http.Request) {
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
	rv.Stash.Add(payload.Key)
	msg := fmt.Sprintf("Bucket %s has been added.", payload.Key)
	rv.Response(w, r, msg)
	return
}

func (rv *RESTView) JSONResponse(w http.ResponseWriter, r *http.Request, data []string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(utils.ToJSON(data)))
}

func (rv *RESTView) Response(w http.ResponseWriter, r *http.Request, data string) {
	w.Write([]byte(data))
}

type BucketPayload struct {
	Key string
}

func main() {
	http.Handle("/", NewRestView())
	http.ListenAndServe(":8000", nil)
}
