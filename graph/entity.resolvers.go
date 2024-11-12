// graph/entity.resolvers.go

package graph

import (
	"context"
	"fmt"

	"github.com/Black-And-White-Club/tcr-bot-user-service/graph/model"
)

// FindUser ByDiscordID is the resolver for the findUser ByDiscordID field.
func (r *entityResolver) FindUserByDiscordID(ctx context.Context, discordID string) (*model.User, error) {
	// Call the UserService's GetUser ByDiscordID method to retrieve the user
	user, err := r.UserService.GetUserByDiscordID(ctx, discordID) // Updated method name
	if err != nil {
		return nil, fmt.Errorf("failed to find user with Discord ID %s: %v", discordID, err)
	}
	return user, nil
}

// Entity returns EntityResolver implementation.
func (r *Resolver) Entity() EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
