package admin

import (
	"context"
	"errors"
	"net/http"

	"github.com/gaydin/journey/admin/api/security"
	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/authentication"
	"github.com/gaydin/journey/date"
	"github.com/gaydin/journey/logger"
	"github.com/gaydin/journey/structure"
	"github.com/gaydin/journey/structure/methods"
)

func (h *Handler) AdminV1APIUserPatch(ctx context.Context, jsonUser *generated.User) (generated.AdminV1APIUserPatchRes, error) {
	userId := security.GetUserID(ctx)

	// Make sure user id is over 0
	if jsonUser.ID.Value < 1 {
		return nil, errors.New("wrong user id")
	} else if userId != jsonUser.ID.Value {
		// Make sure the authenticated user is only changing his/her own data.
		//TODO: Make sure the user is admin when multiple users have been introduced
		return nil, errors.New("you don't have permission to change this data")
	}

	// Get old user data to compare
	tempUser, err := h.db.RetrieveUser(ctx, jsonUser.ID.Value)
	if err != nil {
		return nil, err
	}

	// Make sure user email is provided
	if jsonUser.Email.Value == "" {
		jsonUser.Email.SetTo(tempUser.Email)
	}

	// Make sure name is provided
	if jsonUser.Name.Value == "" {
		jsonUser.Name.SetTo(tempUser.Name)
	}

	// Make sure user slug is provided
	if jsonUser.Slug.Value == "" {
		jsonUser.Slug.SetTo(tempUser.Slug)
	}

	// Check if new name is already taken
	if jsonUser.Name.Value != tempUser.Name {
		_, err = h.db.RetrieveUserByName(ctx, jsonUser.Name.Value)
		if err == nil {
			// The new user name is already taken. Assign the old name.
			// TODO: Return error that will be displayed in the admin interface.
			jsonUser.Name.SetTo(tempUser.Name)
		}
	}

	// Check if new slug is already taken
	if jsonUser.Slug.Value != tempUser.Slug {
		_, err = h.db.RetrieveUserBySlug(ctx, jsonUser.Slug.Value)
		if err == nil {
			// The new user slug is already taken. Assign the old slug.
			// TODO: Return error that will be displayed in the admin interface.
			jsonUser.Slug.SetTo(tempUser.Slug)
		}
	}

	user := structure.User{
		Id:       jsonUser.ID.Value,
		Name:     jsonUser.Name.Value,
		Slug:     jsonUser.Slug.Value,
		Email:    jsonUser.Email.Value,
		Image:    jsonUser.Image.Value,
		Cover:    jsonUser.Cover.Value,
		Bio:      jsonUser.Bio.Value,
		Website:  jsonUser.Website.Value,
		Location: jsonUser.Location.Value,
	}

	err = methods.UpdateUser(ctx, h.db, &user, userId)
	if err != nil {
		return nil, err
	}

	if jsonUser.Password.Value != "" && (jsonUser.Password.Value == jsonUser.PasswordRepeated.Value) { // Update password if a new one was submitted
		encryptedPassword, err := authentication.EncryptPassword(jsonUser.Password.Value)
		if err != nil {
			return nil, err
		}

		err = h.db.UpdateUserPassword(ctx, user.Id, encryptedPassword, date.GetCurrentTime(), jsonUser.ID.Value)
		if err != nil {
			return nil, err
		}
	}

	// Check if the user name was changed. If so, update the session cookie to the new user name.
	if jsonUser.Name.Value != tempUser.Name {
		value := map[string]string{
			"name": jsonUser.Name.Value,
		}

		var cookieValue string

		if encoded, err := authentication.CookieHandler.Encode("session", value); err == nil {
			cookie := &http.Cookie{
				Name:  "session",
				Value: encoded,
				Path:  "/admin/",
			}

			cookieValue = cookie.String()
		}

		err = h.db.UpdateLastLogin(ctx, date.GetCurrentTime(), userId)
		if err != nil {
			logger.FromContext(ctx).Error("couldn't update last login date of a user:", logger.Error(err))
		}

		return &generated.AdminV1APIUserPatchOK{
			SetCookie: generated.NewOptString(cookieValue),
		}, nil
	}

	return &generated.AdminV1APIUserPatchOK{}, nil
}
