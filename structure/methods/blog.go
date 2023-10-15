package methods

import (
	"context"
	"encoding/json"
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

func UpdateBlog(ctx context.Context, blogURL string, db store.Database, b *structure.Blog, userId int64) error {
	// Marshal navigation items to json string
	navigation, err := json.Marshal(b.NavigationItems)
	if err != nil {
		return err
	}

	err = db.UpdateSettings(ctx, b.Title, b.Description, b.Logo, b.Cover, b.PostsPerPage, b.ActiveTheme, string(navigation), date.GetCurrentTime(), userId)
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

func UpdateActiveTheme(activeTheme string, userId int64) error {
	err := store.DB.UpdateActiveTheme(context.Background(), activeTheme, date.GetCurrentTime(), userId)
	if err != nil {
		return err
	}
	// Generate new global blog
	err = GenerateBlog(context.Background(), configuration.Config.Url, store.DB)
	if err != nil {
		log.Panic("Error: couldn't generate blog data:", err)
	}
	return nil
}

func GenerateBlog(ctx context.Context, blogURL string, db store.Database) error {
	// Write lock the global blog
	if Blog != nil {
		Blog.Lock()
		defer Blog.Unlock()
	}
	// Generate blog from db
	blog, err := db.RetrieveBlog(ctx)
	if err != nil {
		return err
	}
	// Add parameters that are not saved in db
	blog.Url = blogURL
	blog.AssetPath = assetPath
	// Create navigation slugs
	for index, _ := range blog.NavigationItems {
		blog.NavigationItems[index].Slug = slug.Generate(ctx, db, blog.NavigationItems[index].Label, "navigation")
	}
	Blog = blog
	return nil
}
