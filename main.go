package main

import (
	"net/http"

	"github.com/marconi/jsonstash/handlers"
)

func main() {
	http.Handle("/", handlers.NewRestView())
	http.ListenAndServe(":8000", nil)
}
