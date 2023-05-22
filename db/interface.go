package db

import (
	"context"

	"github.com/shishberg/mopoke-api/mop"
)

type MopokeDB interface {
	Get(ctx context.Context, name string) (mop.Entity, error)
	Insert(ctx context.Context, e mop.Entity) (string, error)
	Update(ctx context.Context, name string, e mop.Entity) error
	Delete(ctx context.Context, name string) error
}
