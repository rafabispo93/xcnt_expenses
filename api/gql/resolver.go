package gql

import (
	"context"

	"./state"
)

//go:generate gqlgen

type eventProducer interface {
	Produce(ctx context.Context, payload ...interface{}) error
}

// Producers encapsulates the EventProducers the Resolver needs.
type Producers struct {
	ExpenseUpdated          eventProducer
	ExpenseInserted 				eventProducer
}

// Resolver is the root resolver for GraphQL handling.
type Resolver struct {
	producers Producers
	store     state.Store
}

// NewResolver instantiates a new root resolver.
func NewResolver(producers Producers, store state.Store) *Resolver {
	return &Resolver{
		producers: producers,
		store:     store,
	}
}

type mutationResolver struct{ *Resolver }

// Mutation returns the root resolver for all mutations.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type queryResolver struct{ *Resolver }

// Query returns the root resolver for all queries.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
