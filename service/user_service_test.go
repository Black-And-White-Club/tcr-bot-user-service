package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Black-And-White-Club/tcr-bot-user-service/graph/model"
	"github.com/Black-And-White-Club/tcr-bot-user-service/mocks"
	"github.com/Black-And-White-Club/tcr-bot-user-service/service"
	"github.com/jackc/pgx/v5" // Import pgx package for ErrNoRows
	"github.com/pashagolub/pgxmock/v4"
)

// MockUser Service is a mock implementation of the UserService interface
type MockUserService struct {
	GetUserFunc    func(ctx context.Context, discordID string) (*model.User, error)
	CreateUserFunc func(ctx context.Context, input model.UserInput) (*model.User, error)
}

// GetUser  is the mock implementation of the GetUser  method
func (m *MockUserService) GetUser(ctx context.Context, discordID string) (*model.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, discordID)
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

// Resolver is the struct that contains the UserService
type Resolver struct {
	UserService *MockUserService
}

func TestResolver_CreateUser(t *testing.T) {
	mockUserService := &MockUserService{
		GetUserFunc: func(ctx context.Context, discordID string) (*model.User, error) {
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

func TestUserServiceImpl_GetUserByDiscordID(t *testing.T) {
	mockClient, mock, err := mocks.NewPGClientMock()
	if err != nil {
		t.Fatalf("failed to create mock client: %v", err)
	}
	defer mockClient.Close(context.Background()) // Pass context here

	userService := service.NewUserService(mockClient)

	// Pre-populate the mock client with a user
	mock.ExpectQuery("SELECT discord_id, name FROM users").
		WithArgs("12345").
		WillReturnRows(pgxmock.NewRows([]string{"discord_id", "name"}).AddRow("12345", "Test User"))

	tests := []struct {
		name      string
		discordID string
		wantErr   bool
	}{
		{
			name:      "Successful User Retrieval",
			discordID: "12345",
			wantErr:   false,
		},
		{
			name:      "User  Not Found",
			discordID: "notfound",
			wantErr:   true,
		},
		{
			name:      "Missing DiscordID",
			discordID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr && tt.name == "User  Not Found" {
				// Simulate user not found
				mock.ExpectQuery("SELECT discord_id, name FROM users").
					WithArgs(tt.discordID).
					WillReturnError(pgx.ErrNoRows) // Use pgx.ErrNoRows }

				_, err := userService.GetUserByDiscordID(context.Background(), tt.discordID)
				if (err != nil) != tt.wantErr {
					t.Errorf("User  ServiceImpl.GetUser ByDiscordID() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
