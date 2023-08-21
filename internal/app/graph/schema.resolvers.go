package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"

	"github.com/tarkue/tolpi-backend/internal/app/database"
	"github.com/tarkue/tolpi-backend/internal/app/graph/model"
	usercontext "github.com/tarkue/tolpi-backend/internal/app/userContext"
)

var db = database.New()
var ActualTolpi = &model.Tolpi{}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context) (*model.User, error) {
	userId := usercontext.ForContext(ctx).ID
	User := db.CreateUser(userId)

	return User, nil
}

// SetCountry is the resolver for the setCountry field.
func (r *mutationResolver) SetCountry(ctx context.Context, country string) (*model.User, error) {
	userId := usercontext.ForContext(ctx).ID
	User := db.UpdateUserCountry(userId, country)

	return User, nil
}

// CreateTolpi is the resolver for the createTolpi field.
func (r *mutationResolver) CreateTolpi(ctx context.Context, input model.NewTolpi) (*model.Tolpi, error) {
	user := usercontext.ForContext(ctx)
	tolpi := db.CreateTolpi(&input, user.ID)

	ActualTolpi = tolpi

	return tolpi, nil
}

// Tolpies is the resolver for the Tolpies field.
func (r *queryResolver) Tolpies(ctx context.Context, country string) ([]*model.Tolpi, error) {
	return db.GetLastTolpies(country), nil
}

// User is the resolver for the User field.
func (r *queryResolver) User(ctx context.Context, userID string) (*model.User, error) {
	return db.FindUserById(userID), nil
}

// Tolpies is the resolver for the Tolpies field.
func (r *subscriptionResolver) Tolpies(ctx context.Context) (<-chan []*model.Tolpi, error) {
	ch := make(chan []*model.Tolpi)
	userID := usercontext.ForContext(ctx).ID

	go func() {
		t := []*model.Tolpi{}
		var usersId []string
		users := db.GetSubscribes(userID)
		if len(users) > 0 {
			for i := 0; i < len(users); i++ {
				usersId = append(usersId, users[i].UserID)
			}
		} else {
			return
		}

		for {
			if !TolpiContains(t, ActualTolpi) && ActualTolpi.Text != "" {
				t = append(t, ActualTolpi)
			}
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				if !TolpiContains(t, ActualTolpi) && ActualTolpi.Text != "" {
					if Contains(usersId, ActualTolpi.User.UserID) {
						ch <- append(t, ActualTolpi)
					}
				}
			}
		}
	}()
	return ch, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func TolpiContains(a []*model.Tolpi, x *model.Tolpi) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
