package admin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/authentication"
	"github.com/gaydin/journey/date"
)

func (h *Handler) AdminV1APIAuthLoginPost(ctx context.Context, req *generated.ParamsAuthLogin) (generated.AdminV1APIAuthLoginPostRes, error) {
	if !authentication.LoginIsCorrect(ctx, h.db, req.Login, req.Password) {
		return &generated.Unauthorized{}, nil
	}

	authCookie, err := h.newAuthCookie(ctx, req.Login)
	if err != nil {
		return &generated.AdminV1APIAuthLoginPostOK{}, err
	}

	return &generated.AdminV1APIAuthLoginPostOK{
		SetCookie: generated.NewOptString(authCookie),
	}, nil
}

func (h *Handler) newAuthCookie(ctx context.Context, login string) (string, error) {
	value := map[string]string{
		"name": login,
	}

	encoded, err := authentication.CookieHandler.Encode("session", value)
	if err != nil {
		return "", fmt.Errorf("encode %s, err: %w", encoded, err)
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: encoded,
		Path:  "/admin/",
		//HttpOnly: true,
	}

	user, err := h.db.RetrieveUserByName(ctx, login)
	if err != nil {
		return "", fmt.Errorf("couldn't get user ID %w", err)
	}

	err = h.db.UpdateLastLogin(ctx, date.GetCurrentTime(), user.Id)
	if err != nil {
		return "", fmt.Errorf("couldn't update last login date of a user %w", err)
	}

	return cookie.String(), nil
}
