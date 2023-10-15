package admin

import (
	"context"
	"net/http"

	"github.com/gaydin/journey/admin/generated"
)

func (h *Handler) AdminV1APIAuthLogoutPost(_ context.Context) (generated.AdminV1APIAuthLogoutPostRes, error) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "deleted",
		Path:     "/admin/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	return &generated.AdminV1APIAuthLogoutPostOK{
		SetCookie: generated.NewOptString(cookie.String()),
	}, nil
}
