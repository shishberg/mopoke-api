package server

import (
	"net/http"
	"strings"

	"github.com/juju/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Ticket struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
}

func HandleEntity(db *mongo.Client) http.Handler {
	return JSONHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
		if r.Method != http.MethodGet {
			return nil, errors.MethodNotAllowed
		}
		name := strings.TrimPrefix(r.URL.Path, "/")

		tickets := db.Database("mopoke").Collection("tickets")
		var t Ticket
		if err := tickets.FindOne(r.Context(), bson.M{"name": name}).Decode(&t); err != nil {
			return nil, errors.NewNotFound(err, name)
		}
		return t, nil
	})
}
