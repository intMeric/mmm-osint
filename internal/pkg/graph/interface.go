package graph

import "context"

type Graph interface {
	CreateNode(ctx context.Context, node *Node) error
	CreateRelation(ctx context.Context, relation *Relation) error
	GetNode(ctx context.Context, id string) (*Node, error)
	NodeExists(ctx context.Context, id string) (bool, error)
	Close(ctx context.Context) error
}