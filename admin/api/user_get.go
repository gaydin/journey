package admin

import (
	"context"

	"github.com/gaydin/journey/admin/api/security"
	"github.com/gaydin/journey/admin/generated"
)

func (h *Handler) AdminV1APIUserGet(ctx context.Context) (generated.AdminV1APIUserGetRes, error) {
	user, err := h.db.RetrieveUser(ctx, security.GetUserID(ctx))
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
