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

func (h *Handler) AdminV1APIPostPost(ctx context.Context, reqPost *generated.Post) (generated.AdminV1APIPostPostRes, error) {
	userId := security.GetUserID(ctx)

	var postSlug string
	if reqPost.Slug.Value != "" { // Check if user has submitted a custom slug
		postSlug = slug.Generate(ctx, h.db, reqPost.Slug.Value, "posts")
	} else {
		postSlug = slug.Generate(ctx, h.db, reqPost.Title.Value, "posts")
	}

	currentTime := date.GetCurrentTime()

	html := conversion.GenerateHtmlFromMarkdown([]byte(reqPost.Markdown.Value))

	post := structure.Post{
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
		Author:          &structure.User{Id: userId},
	}

	if err := methods.SavePost(ctx, h.url, h.db, &post); err != nil {
		return nil, err
	}

	return &generated.AdminV1APIPostPostOK{}, nil
}
