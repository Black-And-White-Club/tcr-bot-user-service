// service/user_service.go

package service

import (
	"context"
	"fmt"

	"github.com/Black-And-White-Club/tcr-bot-user-service/graph/model"
	"github.com/jackc/pgx/v5"
)

// UserService interface defines methods for user operations
type UserService interface {
	GetUserByDiscordID(ctx context.Context, discordID string) (*model.User, error)
	CreateUser(ctx context.Context, input model.UserInput) (*model.User, error)
}

// UserServiceImpl is the concrete implementation of UserService
type UserServiceImpl struct {
	Client PGClient
}

// NewUser Service creates a new UserService
func NewUserService(client PGClient) *UserServiceImpl {
	return &UserServiceImpl{Client: client}
}

// CreateUser creates a new user in PostgreSQL
func (us *UserServiceImpl) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	// Validate input
	if input.DiscordID == "" || input.Name == "" {
		return nil, fmt.Errorf("DiscordID and Name are required")
	}

	// Check if the user already exists
	user, err := us.Client.GetUserByDiscordID(ctx, input.DiscordID)
	if err != nil && err != pgx.ErrNoRows { // Only proceed if error is not "no rows found"
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}

	if user != nil { // If user exists, return an error
		return nil, fmt.Errorf("user with Discord ID %s already exists", input.DiscordID)
	}

	// User does not exist, so we proceed to create a new user
	newUser := &model.User{
		DiscordID: input.DiscordID,
		Name:      input.Name,
	}

	if err := us.Client.CreateUser(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// GetUser ByDiscordID retrieves a user by Discord ID
func (us *UserServiceImpl) GetUserByDiscordID(ctx context.Context, discordID string) (*model.User, error) {
	// Validate input
	if discordID == "" {
		return nil, fmt.Errorf("DiscordID is required")
	}

	user, err := us.Client.GetUserByDiscordID(ctx, discordID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return user, nil
}
