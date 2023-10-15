package admin

import (
	"context"
	"errors"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	"github.com/gaydin/journey/admin/api/security"
	"github.com/gaydin/journey/admin/generated"
)

func (h *Handler) NewError(_ context.Context, err error) *generated.ErrorStatusCode {
	var ae *security.AuthError
	if errors.As(err, &ae) || errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied) {
		return &generated.ErrorStatusCode{
			StatusCode: http.StatusUnauthorized,
			Response: generated.Error{
				ErrorMessage: generated.NewOptString(err.Error()),
			},
		}
	}

	return &generated.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: generated.Error{
			ErrorMessage: generated.NewOptString(err.Error()),
		},
	}
}
