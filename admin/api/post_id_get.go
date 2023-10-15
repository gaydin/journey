package admin

import (
	"context"
	"strings"

	"github.com/gaydin/journey/admin/generated"
)

func (h *Handler) AdminV1APIPostPostIdGet(ctx context.Context, params generated.AdminV1APIPostPostIdGetParams) (generated.AdminV1APIPostPostIdGetRes, error) {
	post, err := h.db.RetrievePostByID(ctx, params.PostId)
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0, len(post.Tags))
	for index := range post.Tags {
		tags = append(tags, post.Tags[index].Name)
	}

	apiPost := &generated.Post{
		ID:              generated.NewOptInt64(post.Id),
		Title:           generated.NewOptString(post.Title),
		Slug:            generated.NewOptString(post.Slug),
		Markdown:        generated.NewOptString(string(post.Markdown)),
		HTML:            generated.NewOptString(string(post.Html)),
		IsFeatured:      generated.NewOptBool(post.IsFeatured),
		IsPage:          generated.NewOptBool(post.IsPage),
		IsPublished:     generated.NewOptBool(post.IsPublished),
		Image:           generated.NewOptString(post.Image),
		MetaDescription: generated.NewOptString(post.MetaDescription),
		Date:            generated.NewOptDate(*post.Date),
		Tags:            generated.NewOptString(strings.Join(tags, ",")),
	}

	return apiPost, nil
}
