package server

import (
	"encoding/json"
	"net/http"

	"github.com/juju/errors"
)

type JSONHandler func(w http.ResponseWriter, r *http.Request) (any, error)

func NewJSONHandler(jh JSONHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := jh(w, r)
		if err != nil {
			var code int
			switch {
			// TODO: add more errors
			case errors.IsNotFound(err):
				code = http.StatusNotFound
			default:
				code = http.StatusInternalServerError
			}
			// TODO: structured error
			http.Error(w, err.Error(), code)
			return
		}

		out, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(out)
	})
}
