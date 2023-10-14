package methods

import (
	"context"

	"github.com/gaydin/journey/date"
	"github.com/gaydin/journey/store"
	"github.com/gaydin/journey/structure"
)

func SaveUser(ctx context.Context, db store.Database, u *structure.User, hashedPassword string, createdBy int64) error {
	userId, err := db.InsertUser(ctx, u.Name, u.Slug, hashedPassword, u.Email, u.Image, u.Cover, date.GetCurrentTime(), createdBy)
	if err != nil {
		return err
	}
	err = db.InsertRoleUser(ctx, u.Role, userId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(ctx context.Context, db store.Database, u *structure.User, updatedById int64) error {
	err := db.UpdateUser(ctx, u.Id, u.Name, u.Slug, u.Email, u.Image, u.Cover, u.Bio, u.Website, u.Location, date.GetCurrentTime(), updatedById)
	if err != nil {
		return err
	}
	return nil
}
