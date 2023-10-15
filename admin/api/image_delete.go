package admin

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/filenames"
)

func (h *Handler) AdminV1APIImageDelete(_ context.Context, req *generated.AdminV1APIImageDeleteReq) (generated.AdminV1APIImageDeleteRes, error) {
	err := filepath.Walk(filenames.ImagesFilepath, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Base(filePath) == filepath.Base(req.Filename) {
			err := os.Remove(filePath)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}
