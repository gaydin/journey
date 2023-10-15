package admin

import (
	"context"
	"fmt"

	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/authentication"
	"github.com/gaydin/journey/filenames"
	"github.com/gaydin/journey/slug"
	"github.com/gaydin/journey/structure"
	"github.com/gaydin/journey/structure/methods"
)

func (h *Handler) AdminV1APIAuthRegistrationPost(ctx context.Context, req *generated.ParamsAuthRegistration) (generated.AdminV1APIAuthRegistrationPostRes, error) {
	hashedPassword, err := authentication.EncryptPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("encrypt password %w", err)
	}

	user := structure.User{
		Name:  req.Login,
		Slug:  slug.Generate(ctx, h.db, req.Login, "users"),
		Email: req.Email,
		Image: filenames.DefaultUserImageFilename,
		Cover: filenames.DefaultUserCoverFilename,
		Role:  4,
	}

	err = methods.SaveUser(ctx, h.db, &user, hashedPassword, 1)
	if err != nil {
		return nil, fmt.Errorf("save user %w", err)
	}

	authCookie, err := h.newAuthCookie(ctx, req.Login)
	if err != nil {
		return &generated.AdminV1APIAuthRegistrationPostOK{}, err
	}

	return &generated.AdminV1APIAuthRegistrationPostOK{
		SetCookie: generated.NewOptString(authCookie),
	}, nil
}
