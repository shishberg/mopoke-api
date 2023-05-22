package mongomop

import (
	"context"

	"github.com/juju/errors"
	"github.com/shishberg/mopoke-api/db"
	"github.com/shishberg/mopoke-api/mop"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type conn struct {
	db      *mongo.Database
	tickets *mongo.Collection
	rel     *mongo.Collection
}

func New(db *mongo.Database) db.MopokeDB {
	return &conn{
		db:      db,
		tickets: db.Collection("tickets"),
		rel:     db.Collection("rel"),
	}
}

func (c *conn) Get(ctx context.Context, name string) (mop.Entity, error) {
	var t ticket
	if err := c.tickets.FindOne(ctx, bson.M{"name": name}).Decode(&t); err != nil {
		return mop.Entity{}, errors.NewNotFound(err, name)
	}
	return t.entity(), nil
}

func (c *conn) Insert(ctx context.Context, e mop.Entity) (string, error) {
	if e.ID != "" {
		return "", errors.BadRequestf("ID not permitted")
	}
	result, err := c.tickets.InsertOne(ctx, entityToTicket(e))
	if err != nil {
		return "", errors.Annotate(err, "failed to insert")
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (c *conn) Update(ctx context.Context, name string, e mop.Entity) error {
	t := entityToTicket(e)
	r, err := c.tickets.UpdateOne(ctx, ticket{Name: name}, bson.D{{"$set", t}})
	if err != nil {
		return errors.Annotate(err, "failed to update")
	}
	if r.MatchedCount == 0 {
		return errors.NotFoundf("%s", name)
	}
	return nil
}

func (c *conn) Delete(ctx context.Context, name string) error {
	r, err := c.tickets.DeleteOne(ctx, ticket{Name: name})
	if err != nil {
		return errors.Annotate(err, "failed to delete")
	}
	if r.DeletedCount == 0 {
		return errors.NotFoundf("%s", name)
	}
	return nil
}
