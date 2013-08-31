package handlers

import (
	// "fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/pat"
	"github.com/marconi/jsonstash/bucket"
	"github.com/marconi/jsonstash/utils"
)

type RestHandler struct {
	Router *pat.Router
	Stash  *bucket.Stash
}

func NewRestHandler() *pat.Router {
	rh := &RestHandler{
		Router: pat.New(),
		Stash:  bucket.NewStash(),
	}

	// mount handlers
	rh.MountHandler("/buckets/{key:[-_0-9A-Za-z]+}",
		NewBucketHandler(rh),
		[]string{"GET", "POST", "DELETE"})

	rh.MountHandler("/buckets",
		NewBucketListHandler(rh),
		[]string{"GET", "POST"})

	return rh.Router
}

func (rh *RestHandler) MountHandler(path string, view interface{}, methods []string) {
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
		routerMethod := reflect.ValueOf(rh.Router).MethodByName(m)
		handlerMethod := reflect.ValueOf(view).MethodByName(m)

		pathVal := reflect.ValueOf(path)
		handlerVal := reflect.ValueOf(wrapper(&handlerMethod))
		routerMethod.Call([]reflect.Value{pathVal, handlerVal})
	}
}

func (rh *RestHandler) JSONResponse(w http.ResponseWriter, r *http.Request, data []string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(utils.ToJSON(data)))
}

func (rh *RestHandler) Response(w http.ResponseWriter, r *http.Request, data string) {
	w.Write([]byte(data))
}
