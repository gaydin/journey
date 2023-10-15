package methods

import (
	"context"
	"log"

	"github.com/gaydin/journey/date"
	"github.com/gaydin/journey/store"
	"github.com/gaydin/journey/structure"
)

func SavePost(ctx context.Context, blogURL string, db store.Database, p *structure.Post) error {
	tagIds := make([]int64, 0)
	// Insert tags
	for _, tag := range p.Tags {
		// Tag slug might already be in database
		tagId, err := db.RetrieveTagIDBySlug(ctx, tag.Slug)
		if err != nil {
			// Tag is probably not in database yet
			tagId, err = db.InsertTag(ctx, tag.Name, tag.Slug, date.GetCurrentTime(), p.Author.Id)
			if err != nil {
				return err
			}
		}
		if tagId != 0 {
			tagIds = append(tagIds, tagId)
		}
	}
	// Insert post
	postId, err := db.InsertPost(ctx, p.Title, p.Slug, p.Markdown, p.Html, p.IsFeatured, p.IsPage, p.IsPublished, p.MetaDescription, p.Image, *p.Date, p.Author.Id)
	if err != nil {
		return err
	}
	// Insert postTags
	for _, tagId := range tagIds {
		err = db.InsertPostTag(ctx, postId, tagId)
		if err != nil {
			return err
		}
	}
	// Generate new global blog
	err = GenerateBlog(ctx, blogURL, db)
	if err != nil {
		log.Panic("Error: couldn't generate blog data:", err)
	}
	return nil
}

func UpdatePost(ctx context.Context, blogURL string, db store.Database, p *structure.Post) error {
	tagIds := make([]int64, 0)
	// Insert tags
	for _, tag := range p.Tags {
		// Tag slug might already be in database
		tagId, err := db.RetrieveTagIDBySlug(ctx, tag.Slug)
		if err != nil {
			// Tag is probably not in database yet
			tagId, err = db.InsertTag(ctx, tag.Name, tag.Slug, date.GetCurrentTime(), p.Author.Id)
			if err != nil {
				return err
			}
		}
		if tagId != 0 {
			tagIds = append(tagIds, tagId)
		}
	}
	// Update post
	err := db.UpdatePost(ctx, p.Id, p.Title, p.Slug, p.Markdown, p.Html, p.IsFeatured, p.IsPage, p.IsPublished, p.MetaDescription, p.Image, *p.Date, p.Author.Id)
	if err != nil {
		return err
	}
	// Delete old postTags
	err = db.DeletePostTagsForPostID(ctx, p.Id)
	// Insert postTags
	if err != nil {
		return err
	}
	for _, tagId := range tagIds {
		err = db.InsertPostTag(ctx, p.Id, tagId)
		if err != nil {
			return err
		}
	}
	// Generate new global blog
	err = GenerateBlog(ctx, blogURL, db)
	if err != nil {
		log.Panic("Error: couldn't generate blog data:", err)
	}
	return nil
}

func DeletePost(ctx context.Context, blogURL string, db store.Database, postId int64) error {
	err := db.DeletePostByID(ctx, postId)
	if err != nil {
		return err
	}
	// Generate new global blog
	err = GenerateBlog(ctx, blogURL, db)
	if err != nil {
		log.Panic("Error: couldn't generate blog data:", err)
	}
	return nil
}
