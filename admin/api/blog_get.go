package admin

import (
	"context"

	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/structure/methods"
	"github.com/gaydin/journey/templates"
)

func (h *Handler) AdminV1APIBlogGet(_ context.Context) (generated.AdminV1APIBlogGetRes, error) {
	// Read lock the global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()

	navigationItems := make([]generated.Navigation, 0, len(methods.Blog.NavigationItems))
	for i := range methods.Blog.NavigationItems {
		navigationItems = append(navigationItems, generated.Navigation{
			Label: generated.NewOptString(methods.Blog.NavigationItems[i].Label),
			URL:   generated.NewOptString(methods.Blog.NavigationItems[i].Url),
		})
	}

	return &generated.Blog{
		URL:             generated.NewOptString(methods.Blog.Url),
		Title:           generated.NewOptString(methods.Blog.Title),
		Description:     generated.NewOptString(methods.Blog.Description),
		Logo:            generated.NewOptString(methods.Blog.Logo),
		Cover:           generated.NewOptString(methods.Blog.Cover),
		PostsPerPage:    generated.NewOptInt64(methods.Blog.PostsPerPage),
		Themes:          templates.GetAllThemes(),
		ActiveTheme:     generated.NewOptString(methods.Blog.ActiveTheme),
		NavigationItems: navigationItems,
	}, nil
}
