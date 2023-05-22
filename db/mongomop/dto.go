package mongomop

import (
	"github.com/shishberg/mopoke-api/mop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ticket struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description,omitempty"`
}

func (t ticket) entity() mop.Entity {
	return mop.Entity{
		ID:          t.ID.Hex(),
		Name:        t.Name,
		Title:       t.Title,
		Description: t.Description,
	}
}

func entityToTicket(e mop.Entity) ticket {
	return ticket{
		Name:        e.Name,
		Title:       e.Title,
		Description: e.Description,
	}
}
