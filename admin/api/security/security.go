package security

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/authentication"
	"github.com/gaydin/journey/store"
)

func New(db store.Database) *Handler {
	return &Handler{
		db: db,
	}
}

type Handler struct {
	db store.Database
}

type AuthError struct {
	generated.Unauthorized
	m string
}

func (ae *AuthError) Error() string {
	return ae.m
}

func (sh *Handler) HandleCookieAuth(ctx context.Context, _ string, t generated.CookieAuth) (context.Context, error) {
	userName := authentication.GetUserNameString(t.APIKey)
	if userName == "" {
		return nil, &AuthError{m: "not logged in"}
	}

	user, err := sh.db.RetrieveUserByName(ctx, userName)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, errors.New("unknown user")
	case err != nil:
		return nil, err
	}

	return context.WithValue(ctx, contextUserIDKey{}, user.Id), nil
}

type contextUserIDKey struct{}

func GetUserID(ctx context.Context) int64 {
	userName, ok := ctx.Value(contextUserIDKey{}).(int64)
	if !ok {
		return 0
	}

	return userName
}
