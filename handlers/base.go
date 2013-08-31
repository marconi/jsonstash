package handlers

import (
	"github.com/marconi/jsonstash/bucket"
)

type BaseHandler struct {
	RestHandler *RestHandler
	Stash       *bucket.Stash
}
