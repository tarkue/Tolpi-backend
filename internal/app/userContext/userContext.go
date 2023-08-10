package usercontext

import "context"

type contextKey struct {
	name string
}

type UserContext struct {
	ID string
}

var UserCtxKey = &contextKey{"user"}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *UserContext {
	raw, _ := ctx.Value(UserCtxKey).(*UserContext)
	return raw
}
