package admin

import (
	"context"
	"errors"

	"github.com/gaydin/journey/admin/api/security"
	"github.com/gaydin/journey/admin/generated"
)

func (h *Handler) AdminV1APIUserIDGet(ctx context.Context, params generated.AdminV1APIUserIDGetParams) (generated.AdminV1APIUserIDGetRes, error) {
	sessionUserId := security.GetUserID(ctx)
	//TODO validation
	if params.ID != sessionUserId {
		// Make sure the authenticated user is only accessing his/her own data.
		//TODO: Make sure the user is admin when multiple users have been introduced
		return nil, errors.New("you don't have permission to access this data")
	}

	user, err := h.db.RetrieveUser(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &generated.User{
		ID:       generated.NewOptInt64(user.Id),
		Name:     generated.NewOptString(user.Name),
		Slug:     generated.NewOptString(user.Slug),
		Email:    generated.NewOptString(user.Email),
		Image:    generated.NewOptString(user.Image),
		Cover:    generated.NewOptString(user.Cover),
		Bio:      generated.NewOptString(user.Bio),
		Website:  generated.NewOptString(user.Website),
		Location: generated.NewOptString(user.Location),
	}, nil
}
