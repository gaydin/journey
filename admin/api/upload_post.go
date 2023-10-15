package admin

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/gaydin/journey/admin/generated"
	"github.com/gaydin/journey/date"
	"github.com/gaydin/journey/filenames"
)

func (h *Handler) AdminV1APIUploadPost(_ context.Context, req generated.OptAdminV1APIUploadPostReq) (generated.AdminV1APIUploadPostRes, error) {
	files, ok := req.Get()
	if !ok {
		return nil, nil
	}

	// Slice to hold all paths to the files
	allFilePaths := make(generated.AdminV1APIUploadPostOKApplicationJSON, 0)
	// Copy each part to destination.
	for _, part := range files.Multiplefiles {
		// If part.Name() is empty, skip this iteration.
		if part.Name == "" {
			continue
		}

		// Folder structure: year/month/randomname
		currentDate := date.GetCurrentTime()
		filePath := filepath.Join(filenames.ImagesFilepath, currentDate.Format("2006"), currentDate.Format("01"))

		if err := os.MkdirAll(filePath, 0777); err != nil {
			return nil, err
		}

		dst, err := os.Create(filepath.Join(filePath, strconv.FormatInt(currentDate.Unix(), 10)+"_"+uuid.New().String()+filepath.Ext(part.Name)))
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, part.File); err != nil {
			return nil, err
		}

		// Rewrite to file path on server
		filePath = strings.Replace(dst.Name(), filenames.ImagesFilepath, "/images", 1)
		// Make sure to always use "/" as path separator (to make a valid url that we can use on the blog)
		filePath = filepath.ToSlash(filePath)
		allFilePaths = append(allFilePaths, filePath)
	}

	return &allFilePaths, nil
}
