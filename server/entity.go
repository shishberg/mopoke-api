package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/juju/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Ticket struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description,omitempty"`
}

func HandleEntity(db *mongo.Client) http.Handler {
	return JSONHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
		name := strings.TrimPrefix(r.URL.Path, "/")
		tickets := db.Database("mopoke").Collection("tickets")

		switch r.Method {
		case http.MethodGet:
			return getEntity(r, name, tickets)
		case http.MethodPost:
			return postEntity(r, name, tickets)
		case http.MethodPut:
			return putEntity(r, name, tickets)
		default:
			return nil, errors.MethodNotAllowed
		}

	})
}

func getEntity(r *http.Request, name string, tickets *mongo.Collection) (*Ticket, error) {
	var t Ticket
	if err := tickets.FindOne(r.Context(), bson.M{"name": name}).Decode(&t); err != nil {
		return nil, errors.NewNotFound(err, name)
	}
	return &t, nil
}

type PostResponse struct {
	ID string
}

func postEntity(r *http.Request, name string, tickets *mongo.Collection) (any, error) {
	var t Ticket
	// TODO: allow t.ID non-empty?
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, errors.NewBadRequest(err, "failed to parse ticket JSON")
	}
	result, err := tickets.InsertOne(r.Context(), t)
	if err != nil {
		return nil, errors.Annotate(err, "failed to insert")
	}
	return PostResponse{ID: fmt.Sprint(result.InsertedID)}, nil
}

func putEntity(r *http.Request, name string, tickets *mongo.Collection) (any, error) {
	var t Ticket
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, errors.NewBadRequest(err, "failed to parse ticket JSON")
	}
	result, err := tickets.UpdateOne(r.Context(), Ticket{Name: name}, bson.D{{"$set", t}})
	if err != nil {
		return nil, errors.Annotate(err, "failed to update")
	}
	return PostResponse{ID: fmt.Sprint(result.UpsertedID)}, nil
}
