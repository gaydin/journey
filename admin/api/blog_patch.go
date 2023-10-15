package admin

import (
	"context"
	"strings"

	"github.com/gaydin/journey/admin/api/security"
	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/structure"
	"github.com/gaydin/journey/structure/methods"
)

func (h *Handler) AdminV1APIBlogPatch(ctx context.Context, blogReq *generated.Blog) (generated.AdminV1APIBlogPatchRes, error) {
	userId := security.GetUserID(ctx)

	// Make sure postPerPage is over 0
	if blogReq.PostsPerPage.Value < 1 {
		blogReq.PostsPerPage.SetTo(1)
	}

	// Remove blog url in front of navigation urls
	for index := range blogReq.NavigationItems {
		if strings.HasPrefix(blogReq.NavigationItems[index].URL.Value, blogReq.URL.Value) {
			blogReq.NavigationItems[index].URL.SetTo(strings.Replace(blogReq.NavigationItems[index].URL.Value, blogReq.URL.Value, "", 1))
			// If we removed the blog url, there should be a / in front of the url
			if !strings.HasPrefix(blogReq.NavigationItems[index].URL.Value, "/") {
				blogReq.NavigationItems[index].URL.SetTo("/" + blogReq.NavigationItems[index].URL.Value)
			}
		}
	}

	// Retrieve old blog settings for comparison
	blog, err := h.db.RetrieveBlog(ctx)
	if err != nil {
		return nil, err
	}

	navigationItems := make([]structure.Navigation, 0, len(blogReq.NavigationItems))
	for i := range blogReq.NavigationItems {
		navigationItems = append(navigationItems, structure.Navigation{
			Label: blogReq.NavigationItems[i].Label.Value,
			Url:   blogReq.NavigationItems[i].URL.Value,
		})
	}

	tempBlog := structure.Blog{
		Url:             h.url,
		Title:           blogReq.Title.Value,
		Description:     blogReq.Description.Value,
		Logo:            blogReq.Logo.Value,
		Cover:           blogReq.Cover.Value,
		AssetPath:       "/assets/",
		PostCount:       blog.PostCount,
		PostsPerPage:    blogReq.PostsPerPage.Value,
		ActiveTheme:     blogReq.ActiveTheme.Value,
		NavigationItems: navigationItems,
	}

	err = methods.UpdateBlog(ctx, h.url, h.db, &tempBlog, userId)
	if err != nil {
		return nil, err
	}

	// Check if active theme setting has been changed, if so, generate templates from new theme
	if tempBlog.ActiveTheme != blog.ActiveTheme {
		// TODO: new render
		//err = h.tpl.ParseNew(tempBlog.ActiveTheme)
		//if err != nil {
		// If there's an error while generating the new templates, the whole program must be stopped.
		//log.Println("Fatal error: Template data couldn't be generated from theme files: " + err.Error())
		//return nil, err
		//}
	}

	return &generated.AdminV1APIBlogPatchOK{}, nil
}
