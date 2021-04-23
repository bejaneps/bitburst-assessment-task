// Code generated by sqlc. DO NOT EDIT.

package objects

import (
	"context"
)

type Querier interface {
	DeleteNotSeenObjects(ctx context.Context) ([]int32, error)
	InsertObjectOrUpdate(ctx context.Context, oID int32) (int32, error)
	InsertObjectsOrUpdate(ctx context.Context, oIds []int32) ([]int32, error)
	UpdateObject(ctx context.Context, oID int32) (int32, error)
}

var _ Querier = (*Queries)(nil)
