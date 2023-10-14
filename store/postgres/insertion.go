package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	stmtInsertPost     = "INSERT INTO posts (uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, created_at, created_by, updated_at, updated_by, published_at, published_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id"
	stmtInsertUser     = "INSERT INTO users (uuid, name, slug, password, email, image, cover, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id"
	stmtInsertRoleUser = "INSERT INTO roles_users (role_id, user_id) VALUES ($1, $2)"
	stmtInsertTag      = "INSERT INTO tags (uuid, name, slug, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	stmtInsertPostTag  = "INSERT INTO posts_tags (post_id, tag_id) VALUES ($1, $2)"
	stmtInsertSetting  = "INSERT INTO settings (uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
)

func (s *Storage) InsertPost(ctx context.Context, title string, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, metaDescription string, image string, createdAt time.Time, createdBy int64) (int64, error) {
	var postID int64
	return postID, s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		status := "draft"
		if published {
			status = "published"
		}

		var result pgx.Row
		if published {
			result = tx.QueryRow(
				ctx,
				stmtInsertPost,
				uuid.New().String(),
				title,
				slug,
				markdown,
				html,
				featured,
				isPage,
				status,
				metaDescription,
				image,
				createdBy,
				createdAt,
				createdBy,
				createdAt,
				createdBy,
				createdAt,
				createdBy,
			)
		} else {
			result = tx.QueryRow(
				ctx,
				stmtInsertPost,
				uuid.New().String(),
				title,
				slug,
				markdown,
				html,
				featured,
				isPage,
				status,
				metaDescription,
				image,
				createdBy,
				createdAt,
				createdBy,
				createdAt,
				createdBy,
				nil,
				nil,
			)
		}

		return result.Scan(&postID)
	})
}

func (s *Storage) InsertUser(ctx context.Context, name string, slug string, password string, email string, image string, cover string, createdAt time.Time, createdBy int64) (int64, error) {
	var userID int64
	return userID, s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		err := tx.QueryRow(ctx, stmtInsertUser, uuid.New().String(), name, slug, password, email, image, cover, createdAt, createdBy, createdAt, createdBy).Scan(&userID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Storage) InsertRoleUser(ctx context.Context, roleID int64, userID int64) error {
	return s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		return execError(tx.Exec(ctx, stmtInsertRoleUser, roleID, userID))
	})
}

func (s *Storage) InsertTag(ctx context.Context, name string, slug string, createdAt time.Time, createdBy int64) (int64, error) {
	var tagID int64
	return tagID, s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		return tx.QueryRow(ctx, stmtInsertTag, uuid.New().String(), name, slug, createdAt, createdBy, createdAt, createdBy).Scan(&tagID)
	})
}

func (s *Storage) InsertPostTag(ctx context.Context, postID int64, tagID int64) error {
	return s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		return execError(tx.Exec(ctx, stmtInsertPostTag, postID, tagID))
	})
}

func (s *Storage) insertSettingString(ctx context.Context, key string, value string, settingType string, createdAt time.Time) error {
	const defaultUserID = 1
	return s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		return execError(tx.Exec(ctx, stmtInsertSetting, uuid.New().String(), key, value, settingType, createdAt, defaultUserID, createdAt, defaultUserID))
	})
}
