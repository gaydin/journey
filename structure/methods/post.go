package methods

import (
	"context"
	"log"

	"github.com/gaydin/journey/date"
	"github.com/gaydin/journey/store"
	"github.com/gaydin/journey/structure"
)

func SavePost(ctx context.Context, db store.Database, p *structure.Post) error {
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
	postId, err := store.DB.InsertPost(context.Background(), p.Title, p.Slug, p.Markdown, p.Html, p.IsFeatured, p.IsPage, p.IsPublished, p.MetaDescription, p.Image, *p.Date, p.Author.Id)
	if err != nil {
		return err
	}
	// Insert postTags
	for _, tagId := range tagIds {
		err = store.DB.InsertPostTag(context.Background(), postId, tagId)
		if err != nil {
			return err
		}
	}
	// Generate new global blog
	err = GenerateBlog(context.Background())
	if err != nil {
		log.Panic("Error: couldn't generate blog data:", err)
	}
	return nil
}

func UpdatePost(p *structure.Post) error {
	tagIds := make([]int64, 0)
	// Insert tags
	for _, tag := range p.Tags {
		// Tag slug might already be in database
		tagId, err := store.DB.RetrieveTagIDBySlug(context.Background(), tag.Slug)
		if err != nil {
			// Tag is probably not in database yet
			tagId, err = store.DB.InsertTag(context.Background(), tag.Name, tag.Slug, date.GetCurrentTime(), p.Author.Id)
			if err != nil {
				return err
			}
		}
		if tagId != 0 {
			tagIds = append(tagIds, tagId)
		}
	}
	// Update post
	err := store.DB.UpdatePost(context.Background(), p.Id, p.Title, p.Slug, p.Markdown, p.Html, p.IsFeatured, p.IsPage, p.IsPublished, p.MetaDescription, p.Image, *p.Date, p.Author.Id)
	if err != nil {
		return err
	}
	// Delete old postTags
	err = store.DB.DeletePostTagsForPostID(context.Background(), p.Id)
	// Insert postTags
	if err != nil {
		return err
	}
	for _, tagId := range tagIds {
		err = store.DB.InsertPostTag(context.Background(), p.Id, tagId)
		if err != nil {
			return err
		}
	}
	// Generate new global blog
	err = GenerateBlog(context.Background())
	if err != nil {
		log.Panic("Error: couldn't generate blog data:", err)
	}
	return nil
}

func DeletePost(postId int64) error {
	err := store.DB.DeletePostByID(context.Background(), postId)
	if err != nil {
		return err
	}
	// Generate new global blog
	err = GenerateBlog(context.Background())
	if err != nil {
		log.Panic("Error: couldn't generate blog data:", err)
	}
	return nil
}
