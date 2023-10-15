package admin

import (
	"context"

	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/structure/methods"
)

func (h *Handler) AdminV1APIPostPostIdDelete(ctx context.Context, params generated.AdminV1APIPostPostIdDeleteParams) (generated.AdminV1APIPostPostIdDeleteRes, error) {
	return nil, methods.DeletePost(ctx, h.url, h.db, params.PostId)
}
