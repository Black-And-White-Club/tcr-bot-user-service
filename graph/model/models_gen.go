// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

// Mutations available in the User Service.
type Mutation struct {
}

// Queries available in the User Service.
type Query struct {
}

// Represents a user in the system.
type User struct {
	DiscordID string `json:"discordID"`
	Name      string `json:"name"`
	TagNumber *int   `json:"tagNumber,omitempty"`
	Role      string `json:"role"`
}

func (User) IsEntity() {}

// Input type for creating a new user.
type UserInput struct {
	Name      string `json:"name"`
	DiscordID string `json:"discordID"`
}