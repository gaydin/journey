package templates

import (
	"bytes"
	"context"
	"net/http"

	"github.com/gaydin/journey/date"
	"github.com/gaydin/journey/feeds"
	"github.com/gaydin/journey/store"
	"github.com/gaydin/journey/structure"
	"github.com/gaydin/journey/structure/methods"
)

func ShowIndexRss(ctx context.Context, writer http.ResponseWriter) error {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	// 15 posts in rss for now
	posts, err := store.DB.RetrievePostsForIndex(ctx, 15, 0)
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: methods.Blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(writer)
	return err
}

func ShowTagRss(ctx context.Context, writer http.ResponseWriter, slug string) error {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	tag, err := store.DB.RetrieveTagBySlug(ctx, slug)
	if err != nil {
		return err
	}
	// 15 posts in rss for now
	posts, err := store.DB.RetrievePostsByTag(ctx, tag.Id, 15, 0)
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: methods.Blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(writer)
	return err
}

func ShowAuthorRss(ctx context.Context, writer http.ResponseWriter, slug string) error {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	author, err := store.DB.RetrieveUserBySlug(ctx, slug)
	if err != nil {
		return err
	}
	// 15 posts in rss for now
	posts, err := store.DB.RetrievePostsByUser(ctx, author.Id, 15, 0)
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: methods.Blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(writer)
	return err
}

func createFeed(values *structure.RequestData) *feeds.Feed {
	now := date.GetCurrentTime()
	feed := &feeds.Feed{
		Title:       string(values.Blog.Title),
		Description: string(values.Blog.Description),
		Link:        &feeds.Link{Href: string(values.Blog.Url)},
		Updated:     now,
		Image: &feeds.Image{
			Url:   string(values.Blog.Url) + string(values.Blog.Logo),
			Title: string(values.Blog.Title),
			Link:  string(values.Blog.Url),
		},
		Url: string(values.Blog.Url) + "/rss/",
	}
	for i := 0; i < len(values.Posts); i++ {
		if values.Posts[i].Id != 0 {
			// Make link
			var buffer bytes.Buffer
			buffer.WriteString(values.Blog.Url)
			buffer.WriteString("/")
			buffer.WriteString(values.Posts[i].Slug)
			item := &feeds.Item{
				Title:       string(values.Posts[i].Title),
				Description: string(values.Posts[i].Html),
				Link:        &feeds.Link{Href: buffer.String()},
				Id:          string(values.Posts[i].Uuid),
				Author:      &feeds.Author{Name: string(values.Posts[i].Author.Name), Email: ""},
				Created:     *values.Posts[i].Date,
			}
			// If the post has a cover image, add it to the item
			image := string(values.Posts[i].Image)
			if image != "" {
				item.Image = &feeds.Image{
					Url: string(values.Blog.Url) + image,
				}
			}
			feed.Items = append(feed.Items, item)
		}
	}

	return feed
}
