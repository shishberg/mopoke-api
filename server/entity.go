package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/juju/errors"
	"github.com/shishberg/mopoke-api/db"
	"github.com/shishberg/mopoke-api/mop"
)

type handler struct {
	db db.MopokeDB
}

func NewHandler(db db.MopokeDB) http.Handler {
	return JSONHandler{&handler{db}}
}

func (h *handler) ServeJSON(w http.ResponseWriter, r *http.Request) (any, error) {
	log.Println(r.Method, r.URL)
	name := strings.TrimPrefix(r.URL.Path, "/")

	switch r.Method {
	case http.MethodGet:
		return h.getEntity(r, name)
	case http.MethodPost:
		return h.postEntity(r)
	case http.MethodPut:
		return nil, h.putEntity(r, name)
	case http.MethodDelete:
		return nil, h.deleteEntity(r, name)
	default:
		return nil, errors.MethodNotAllowed
	}
}

func (h *handler) getEntity(r *http.Request, name string) (mop.Entity, error) {
	entity, err := h.db.Get(r.Context(), name)
	return entity, errors.Trace(err)
}

func (h *handler) postEntity(r *http.Request) (any, error) {
	var e mop.Entity
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return "", errors.NewBadRequest(err, "failed to parse JSON body as entity")
	}
	id, err := h.db.Insert(r.Context(), e)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return idResponse{id}, nil
}

func (h *handler) putEntity(r *http.Request, name string) error {
	var e mop.Entity
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return errors.NewBadRequest(err, "failed to parse JSON body as entity")
	}
	err := h.db.Update(r.Context(), name, e)
	return errors.Trace(err)
}

func (h *handler) deleteEntity(r *http.Request, name string) error {
	return errors.Trace(h.db.Delete(r.Context(), name))
}
