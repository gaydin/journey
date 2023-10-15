package admin

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/filenames"
	"github.com/gaydin/journey/logger"
)

func (h *Handler) AdminV1APIImagesNumberGet(ctx context.Context, params generated.AdminV1APIImagesNumberGetParams) (generated.AdminV1APIImagesNumberGetRes, error) {
	images := make([]string, 0)
	// Walk all files in images folder
	err := filepath.Walk(filenames.ImagesFilepath, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && (strings.EqualFold(filepath.Ext(filePath), ".jpg") || strings.EqualFold(filepath.Ext(filePath), ".jpeg") || strings.EqualFold(filepath.Ext(filePath), ".gif") || strings.EqualFold(filepath.Ext(filePath), ".png") || strings.EqualFold(filepath.Ext(filePath), ".svg")) {
			// Rewrite to file path on server
			filePath = strings.Replace(filePath, filenames.ImagesFilepath, "/images", 1)
			// Make sure to always use "/" as path separator (to make a valid url that we can use on the blog)
			filePath = filepath.ToSlash(filePath)
			// Prepend file to slice (thus reversing the order)
			images = append([]string{filePath}, images...)
		}
		return nil
	})

	if err != nil {
		logger.FromContext(ctx).Info("apiImagesHandler walk error", logger.Error(err))
	}

	if len(images) == 0 {
		return &generated.AdminV1APIImagesNumberGetOK{}, nil
	}

	imagesPerPage := int64(15)
	start := (params.Number * imagesPerPage) - imagesPerPage
	end := params.Number * imagesPerPage

	if start > int64((len(images))-1) {
		return &generated.AdminV1APIImagesNumberGetOK{}, nil
	}

	if end > int64(len(images)) {
		end = int64(len(images))
	}

	return &generated.AdminV1APIImagesNumberGetOK{Images: images[start:end]}, nil
}
