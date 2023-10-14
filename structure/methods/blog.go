package methods

import (
	"context"
	"log"

	"github.com/gaydin/journey/configuration"
	"github.com/gaydin/journey/date"
	"github.com/gaydin/journey/slug"
	"github.com/gaydin/journey/store"
	"github.com/gaydin/journey/structure"
)

// Global blog - thread safe and accessible by all requests
var Blog *structure.Blog

var assetPath = "/assets/"

func UpdateActiveTheme(activeTheme string, userId int64) error {
	err := store.DB.UpdateActiveTheme(context.Background(), activeTheme, date.GetCurrentTime(), userId)
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

func GenerateBlog(ctx context.Context) error {
	// Write lock the global blog
	if Blog != nil {
		Blog.Lock()
		defer Blog.Unlock()
	}
	// Generate blog from db
	blog, err := store.DB.RetrieveBlog(ctx)
	if err != nil {
		return err
	}
	// Add parameters that are not saved in db
	blog.Url = configuration.Config.Url
	blog.AssetPath = assetPath
	// Create navigation slugs
	for index, _ := range blog.NavigationItems {
		blog.NavigationItems[index].Slug = slug.Generate(ctx, store.DB, blog.NavigationItems[index].Label, "navigation")
	}
	Blog = blog
	return nil
}
