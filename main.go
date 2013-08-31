package main

import (
	// "fmt"
	"net/http"

	"github.com/marconi/jsonstash/handlers"
)

func main() {
	http.Handle("/", handlers.NewRestHandler())
	http.ListenAndServe(":8000", nil)
}
