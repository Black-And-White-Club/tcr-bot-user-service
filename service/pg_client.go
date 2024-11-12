// service/pg_client.go

package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Black-And-White-Club/tcr-bot-user-service/graph/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PGClient interface defines methods for database operations
type PGClient interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByDiscordID(ctx context.Context, discordID string) (*model.User, error)
	Close(ctx context.Context) error
}

// PGClientImpl is the implementation of the PGClient interface
type PGClientImpl struct {
	Pool *pgxpool.Pool
}

// NewPGClient creates a new PGClient
func NewPGClient(dataSourceName string) (*PGClientImpl, error) {
	config, err := pgxpool.ParseConfig(dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &PGClientImpl{Pool: pool}, nil
}

// GetUser ByDiscordID retrieves a user by Discord ID
func (pg *PGClientImpl) GetUserByDiscordID(ctx context.Context, discordID string) (*model.User, error) {
	var user model.User
	err := pg.Pool.QueryRow(ctx, "SELECT discord_id, name FROM users WHERE discord_id = $1", discordID).Scan(&user.DiscordID, &user.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil if user is not found
		}
		log.Printf("Error retrieving user: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// CreateUser  creates a new user in PostgreSQL
func (pg *PGClientImpl) CreateUser(ctx context.Context, user *model.User) error {
	_, err := pg.Pool.Exec(ctx, "INSERT INTO users (discord_id, name) VALUES ($1, $2)", user.DiscordID, user.Name)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Close closes the database connection pool
func (pg *PGClientImpl) Close(ctx context.Context) error {
	pg.Pool.Close()
	log.Println("PostgreSQL connection pool closed")
	return nil
}
