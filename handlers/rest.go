package handlers

import (
	"net/http"
	"reflect"
	"strings"

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
