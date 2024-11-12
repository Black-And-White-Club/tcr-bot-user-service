package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/Black-And-White-Club/tcr-bot-user-service/graph/model"
	"github.com/Black-And-White-Club/tcr-bot-user-service/mocks" // Ensure this import is present
)

func TestEntityResolver_FindUserByDiscordID(t *testing.T) {
	mockPGClient, _, err := mocks.NewPGClientMock() // Create the mock PG client
	if err != nil {
		t.Fatalf("failed to create mock PG client: %v", err)
	}
	defer mockPGClient.Close(context.Background()) // Ensure to close the mock client after the test

	// Set up the mock user service
	mockUserService := &MockUserService{
		GetUserByDiscordIDFunc: func(ctx context.Context, discordID string) (*model.User, error) {
			if discordID == "validID" {
				return &model.User{DiscordID: discordID, Name: "Test User"}, nil
			}
			return nil, errors.New("user not found")
		},
	}

	resolver := &Resolver{UserService: mockUserService}
	entityResolver := &entityResolver{resolver}

	tests := []struct {
		discordID string
		want      *model.User
		wantErr   bool
	}{
		{"validID", &model.User{DiscordID: "validID", Name: "Test User"}, false},
		{"invalidID", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.discordID, func(t *testing.T) {
			got, err := entityResolver.FindUserByDiscordID(context.Background(), tt.discordID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUser ByDiscordID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.DiscordID != tt.want.DiscordID {
				t.Errorf("FindUser ByDiscordID() = %v, want %v", got, tt.want)
			}
		})
	}
}
