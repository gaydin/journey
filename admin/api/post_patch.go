package admin

import (
	"context"

	"github.com/gaydin/journey/admin/api/security"
	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/conversion"
	"github.com/gaydin/journey/date"
	"github.com/gaydin/journey/slug"
	"github.com/gaydin/journey/structure"
	"github.com/gaydin/journey/structure/methods"
)

func (h *Handler) AdminV1APIPostPatch(ctx context.Context, reqPost *generated.Post) (generated.AdminV1APIPostPatchRes, error) {
	userID := security.GetUserID(ctx)

	// Get current slug of post
	post, err := h.db.RetrievePostByID(ctx, reqPost.ID.Value)
	if err != nil {
		return nil, err
	}

	var postSlug string
	if reqPost.Slug.Value != post.Slug { // Check if user has submitted a custom slug
		postSlug = slug.Generate(ctx, h.db, reqPost.Slug.Value, "posts")
	} else {
		postSlug = post.Slug
	}

	html := conversion.GenerateHtmlFromMarkdown([]byte(reqPost.Markdown.Value))

	currentTime := date.GetCurrentTime()
	post = &structure.Post{
		Id:              reqPost.ID.Value,
		Title:           reqPost.Title.Value,
		Slug:            postSlug,
		Markdown:        []byte(reqPost.Markdown.Value),
		Html:            html,
		IsFeatured:      reqPost.IsFeatured.Value,
		IsPage:          reqPost.IsPage.Value,
		IsPublished:     reqPost.IsPublished.Value,
		MetaDescription: reqPost.MetaDescription.Value,
		Image:           reqPost.Image.Value,
		Date:            &currentTime,
		Tags:            methods.GenerateTagsFromCommaString(ctx, h.db, reqPost.Tags.Value),
		Author:          &structure.User{Id: userID},
	}

	if err := methods.UpdatePost(ctx, h.url, h.db, post); err != nil {
		return nil, err
	}

	return &generated.AdminV1APIPostPatchOK{}, err
}
