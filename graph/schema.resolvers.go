// !test
package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"context"
	"fmt"

	"github.com/Black-And-White-Club/tcr-bot-user-service/graph/model"
)

// CreateUser  is the resolver for the createUser  field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	// Call the UserService's CreateUser  method to create a new user
	user, err := r.UserService.CreateUser(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	return user, nil
}

// GetUser  is the resolver for the getUser  field.
func (r *queryResolver) GetUser(ctx context.Context, discordID string) (*model.User, error) {
	// Call the UserService's GetUser ByDiscordID method to retrieve the user
	user, err := r.UserService.GetUserByDiscordID(ctx, discordID) // Updated method name
	if err != nil {
		return nil, fmt.Errorf("failed to get user with Discord ID %s: %v", discordID, err)
	}
	return user, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }