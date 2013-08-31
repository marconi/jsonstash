package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/marconi/jsonstash/bucket"
	"github.com/marconi/jsonstash/utils"
)

type BucketPayload struct {
	Key string
}

type RESTView struct {
	Router *mux.Router
	Stash  *bucket.Stash
}

func NewRestView() *mux.Router {
	rv := &RESTView{
		Router: mux.NewRouter(),
		Stash:  bucket.NewStash(),
	}

	// mount views
	bucketView := &BucketView{RESTView: rv, Stash: rv.Stash}
	rv.AddView("/buckets", bucketView, []string{"GET", "POST"})

	return rv.Router
}

func (rv *RESTView) AddView(path string, view interface{}, methods []string) {
	// factory function that returns a handler wrapper
	wrapper := func(method *reflect.Value) func(w http.ResponseWriter, r *http.Request) {
		f := func(w http.ResponseWriter, r *http.Request) {
			wVal := reflect.ValueOf(w)
			rVal := reflect.ValueOf(r)
			method.Call([]reflect.Value{wVal, rVal})
		}
		return f
	}

	// add handler for each method
	for _, m := range methods {
		m = strings.Title(strings.ToLower(m))
		method := reflect.ValueOf(view).MethodByName(m)
		rv.Router.HandleFunc(path, wrapper(&method)).Methods(m)
	}
}

func (rv *RESTView) JSONResponse(w http.ResponseWriter, r *http.Request, data []string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(utils.ToJSON(data)))
}

func (rv *RESTView) Response(w http.ResponseWriter, r *http.Request, data string) {
	w.Write([]byte(data))
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

func main() {
	http.Handle("/", NewRestView())
	http.ListenAndServe(":8000", nil)
}
