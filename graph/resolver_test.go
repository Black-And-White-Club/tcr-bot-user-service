package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/Black-And-White-Club/tcr-bot-user-service/graph/model"
)

// MockUser Service is a mock implementation of the UserService interface
type MockUserService struct {
	GetUserByDiscordIDFunc func(ctx context.Context, discordID string) (*model.User, error)
	CreateUserFunc         func(ctx context.Context, input model.UserInput) (*model.User, error)
}

// GetUser ByDiscordID is the mock implementation of the GetUser ByDiscordID method
func (m *MockUserService) GetUserByDiscordID(ctx context.Context, discordID string) (*model.User, error) {
	if m.GetUserByDiscordIDFunc != nil {
		return m.GetUserByDiscordIDFunc(ctx, discordID)
	}
	return nil, nil
}

// CreateUser  is the mock implementation of the CreateUser  method
func (m *MockUserService) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, input)
	}
	return nil, nil
}

func TestResolver_GetUser(t *testing.T) {
	mockUserService := &MockUserService{
		GetUserByDiscordIDFunc: func(ctx context.Context, discordID string) (*model.User, error) {
			if discordID == "existingID" {
				return &model.User{DiscordID: discordID, Name: "Existing User"}, nil
			}
			return nil, errors.New("user not found")
		},
	}

	resolver := &Resolver{UserService: mockUserService}

	tests := []struct {
		name      string
		discordID string
		want      *model.User
		wantErr   bool
	}{
		{"User  Found", "existingID", &model.User{DiscordID: "existingID", Name: "Existing User"}, false},
		{"User  Not Found", "nonExistingID", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolver.UserService.GetUserByDiscordID(context.Background(), tt.discordID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resolver.GetUser ByDiscordID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.DiscordID != tt.want.DiscordID {
				t.Errorf("Resolver.GetUser ByDiscordID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_CreateUser(t *testing.T) {
	mockUserService := &MockUserService{
		GetUserByDiscordIDFunc: func(ctx context.Context, discordID string) (*model.User, error) {
			if discordID == "existingID" {
				return &model.User{DiscordID: discordID, Name: "Existing User"}, nil
			}
			return nil, errors.New("user not found")
		},
		CreateUserFunc: func(ctx context.Context, input model.UserInput) (*model.User, error) {
			if input.DiscordID == "newID" {
				return &model.User{DiscordID: input.DiscordID, Name: input.Name}, nil
			}
			return nil, errors.New("failed to create user")
		},
	}

	resolver := &Resolver{UserService: mockUserService}

	tests := []struct {
		name    string
		input   model.UserInput
		want    *model.User
		wantErr bool
	}{
		{"Create_User_Successfully", model.UserInput{DiscordID: "newID", Name: "New User"}, &model.User{DiscordID: "newID", Name: "New User"}, false},
		{"Create_User_Already_Exists", model.UserInput{DiscordID: "existingID", Name: "Existing User"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolver.UserService.CreateUser(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resolver.CreateUser () error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.DiscordID != tt.want.DiscordID {
				t.Errorf("Resolver.CreateUser () = %v, want %v", got, tt.want)
			}
		})
	}
}
