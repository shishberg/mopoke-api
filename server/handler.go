package server

import (
	"encoding/json"
	"net/http"

	"github.com/juju/errors"
)

type idResponse struct {
	ID string `json:"id"`
}

type JSONServer interface {
	ServeJSON(w http.ResponseWriter, r *http.Request) (any, error)
}

type JSONHandler struct {
	JSONServer
}

func (jh JSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := jh.JSONServer.ServeJSON(w, r)
	if err != nil {
		var code int
		switch {
		// TODO: add more errors
		case errors.Is(err, errors.NotFound):
			code = http.StatusNotFound
		case errors.Is(err, errors.MethodNotAllowed):
			code = http.StatusMethodNotAllowed
		default:
			code = http.StatusInternalServerError
		}
		// TODO: structured error
		http.Error(w, err.Error(), code)
		return
	}

	if resp == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	out, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(out)
}
