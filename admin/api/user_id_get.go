package admin

import (
	"context"

	"github.com/gaydin/journey/admin/api/security"
	"github.com/gaydin/journey/admin/generated"
)

func (h *Handler) AdminV1APIUseridGet(ctx context.Context) (generated.AdminV1APIUseridGetRes, error) {
	return &generated.AdminV1APIUseridGetOK{
		ID: generated.NewOptInt64(security.GetUserID(ctx)),
	}, nil
}
