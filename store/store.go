package store

import (
	"context"
	"time"

	"github.com/gaydin/journey/configuration"
	"github.com/gaydin/journey/store/postgres"
	"github.com/gaydin/journey/structure"
)

type Database interface {
	DeletePostTagsForPostID(ctx context.Context, id int64) error
	DeletePostByID(ctx context.Context, id int64) error

	InsertPost(ctx context.Context, title string, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, metaDescription string, image string, createdAt time.Time, createdBy int64) (int64, error)
	InsertUser(ctx context.Context, name string, slug string, password string, email string, image string, cover string, createdAt time.Time, createdBy int64) (int64, error)
	InsertRoleUser(ctx context.Context, roleID, userId int64) error
	InsertTag(ctx context.Context, name string, slug string, createdAt time.Time, createdBy int64) (int64, error)
	InsertPostTag(ctx context.Context, postID int64, tagID int64) error

	RetrievePostByID(ctx context.Context, id int64) (*structure.Post, error)
	RetrievePostBySlug(ctx context.Context, slug string) (*structure.Post, error)
	RetrievePostsByUser(ctx context.Context, userID int64, limit int64, offset int64) ([]structure.Post, error)
	RetrievePostsByTag(ctx context.Context, tagID int64, limit int64, offset int64) ([]structure.Post, error)
	RetrievePostsForIndex(ctx context.Context, limit int64, offset int64) ([]structure.Post, error)
	RetrievePostsForApi(ctx context.Context, limit int64, offset int64) ([]structure.Post, error)
	RetrieveNumberOfPosts(ctx context.Context) (int64, error)
	RetrieveNumberOfPostsByUser(ctx context.Context, userID int64) (int64, error)
	RetrieveNumberOfPostsByTag(ctx context.Context, tagID int64) (int64, error)
	RetrieveUser(ctx context.Context, id int64) (*structure.User, error)
	RetrieveUserBySlug(ctx context.Context, slug string) (*structure.User, error)
	RetrieveUserByName(ctx context.Context, name string) (*structure.User, error)
	RetrieveTags(ctx context.Context, postID int64) ([]structure.Tag, error)
	RetrieveTag(ctx context.Context, tagID int64) (*structure.Tag, error)
	RetrieveTagBySlug(ctx context.Context, slug string) (*structure.Tag, error)
	RetrieveTagIDBySlug(ctx context.Context, slug string) (int64, error)
	RetrieveHashedPasswordForUser(ctx context.Context, name string) ([]byte, error)
	RetrieveBlog(ctx context.Context) (*structure.Blog, error)
	RetrieveActiveTheme(ctx context.Context) (string, error)
	RetrieveUsersCount(ctx context.Context) int

	UpdatePost(ctx context.Context, id int64, title string, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, metaDescription string, image string, updatedAt time.Time, updatedBy int64) error
	UpdateActiveTheme(ctx context.Context, activeTheme string, updatedAt time.Time, updatedBy int64) error
	UpdateUser(ctx context.Context, id int64, name string, slug string, email string, image string, cover string, bio string, website string, location string, updatedAt time.Time, updatedBy int64) error
	UpdateLastLogin(ctx context.Context, logInDate time.Time, userID int64) error
	UpdateUserPassword(ctx context.Context, id int64, password string, updatedAt time.Time, updatedBy int64) error
	UpdateSettings(
		ctx context.Context,
		title string,
		description string,
		logo string,
		cover string,
		postsPerPage int64,
		activeTheme string,
		navigation string,
		updatedAt time.Time,
		updatedBy int64,
	) error
}

func New(config *configuration.Configuration) (Database, error) {
	db, err := postgres.New(config.Database.DSN)
	if err != nil {
		return nil, err
	}

	DB = db

	return db, nil
}
