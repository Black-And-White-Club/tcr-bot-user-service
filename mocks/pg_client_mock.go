package mocks

import (
	"context"

	"github.com/Black-And-White-Club/tcr-bot-user-service/graph/model"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

// PGClientMock is a mock implementation of PGClient using pgxmock
type PGClientMock struct {
	mock pgxmock.PgxConnIface
}

// NewPGClientMock initializes a new PGClientMock
func NewPGClientMock() (*PGClientMock, pgxmock.PgxConnIface, error) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		return nil, nil, err
	}
	return &PGClientMock{mock: mock}, mock, nil
}

// CreateUser is a mock implementation of the CreateUser method
func (m *PGClientMock) CreateUser(ctx context.Context, user *model.User) error {
	// This function itself doesn't set expectations; they are set in the test cases.
	return nil
}

// GetUserByDiscordID is a mock implementation of the GetUserByDiscordID method
func (m *PGClientMock) GetUserByDiscordID(ctx context.Context, discordID string) (*model.User, error) {
	// Here, we return specific values to simulate the database responses.
	if discordID == "validID" {
		return &model.User{DiscordID: discordID, Name: "Test User"}, nil
	}
	return nil, pgx.ErrNoRows
}

// Close is a mock implementation of the Close method
func (m *PGClientMock) Close(ctx context.Context) error {
	return m.mock.Close(ctx)
}

// MockUser Service is a mock implementation of UserService
type MockUserService struct {
	PGClientMock *PGClientMock
}

// CreateUser  mocks the CreateUser  method of UserService
func (m *MockUserService) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	user := &model.User{DiscordID: input.DiscordID, Name: input.Name}
	if err := m.PGClientMock.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUser ByDiscordID mocks the GetUser ByDiscordID method of UserService
func (m *MockUserService) GetUserByDiscordID(ctx context.Context, discordID string) (*model.User, error) {
	return m.PGClientMock.GetUserByDiscordID(ctx, discordID)
}
