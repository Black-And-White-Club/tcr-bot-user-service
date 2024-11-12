// graph/resolver.go

package graph

import (
	"context"
	"log"

	"github.com/Black-And-White-Club/tcr-bot-user-service/graph/model"
	"github.com/Black-And-White-Club/tcr-bot-user-service/service"
)

// Resolver struct definition
type Resolver struct {
	UserService service.UserService
}

// GetUser  resolver
func (r *Resolver) GetUser(ctx context.Context, discordID string) (*model.User, error) {
	user, err := r.UserService.GetUserByDiscordID(ctx, discordID) // Updated method name
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return nil, err
	}
	return user, nil
}

// CreateUser  resolver
func (r *Resolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	user, err := r.UserService.CreateUser(ctx, input) // Correctly calling the method
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}
	return user, nil
}
