package admin

import (
	"context"
	"errors"
	"strings"

	"github.com/gaydin/journey/admin/generated"
)

func (h *Handler) AdminV1APIPostsNumberGet(ctx context.Context, params generated.AdminV1APIPostsNumberGetParams) (generated.AdminV1APIPostsNumberGetRes, error) {
	if params.Number < 1 {
		return nil, errors.New("invalid page number")
	}

	postsPerPage := int64(15)
	posts, err := h.db.RetrievePostsForApi(ctx, postsPerPage, (int64(params.Number)-1)*postsPerPage)
	if err != nil {
		return nil, err
	}

	respPosts := make(generated.AdminV1APIPostsNumberGetOKApplicationJSON, 0, len(posts))

	for i := range posts {
		tags := make([]string, 0, len(posts[i].Tags))
		for index := range posts[i].Tags {
			tags = append(tags, posts[i].Tags[index].Name)
		}

		respPosts = append(respPosts, generated.Post{
			ID:              generated.NewOptInt64(posts[i].Id),
			Title:           generated.NewOptString(posts[i].Title),
			Slug:            generated.NewOptString(posts[i].Slug),
			Markdown:        generated.NewOptString(string(posts[i].Markdown)),
			HTML:            generated.NewOptString(string(posts[i].Html)),
			IsFeatured:      generated.NewOptBool(posts[i].IsFeatured),
			IsPage:          generated.NewOptBool(posts[i].IsPage),
			IsPublished:     generated.NewOptBool(posts[i].IsPublished),
			Image:           generated.NewOptString(posts[i].Image),
			MetaDescription: generated.NewOptString(posts[i].MetaDescription),
			Date:            generated.NewOptDate(*posts[i].Date),
			Tags:            generated.NewOptString(strings.Join(tags, ",")),
		})
	}

	return &respPosts, nil
}
