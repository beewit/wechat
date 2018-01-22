package mp

import (
	"net/http"
)

type ErrorHandler interface {
	ServeError(http.ResponseWriter, *http.Request, error)
}

type ErrorHandlerFunc func(http.ResponseWriter, *http.Request, error)

func (fn ErrorHandlerFunc) ServeError(w http.ResponseWriter, r *http.Request, err error) {
	fn(w, r, err)
}

var DefaultErrorHandler = ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request, err error) {})
