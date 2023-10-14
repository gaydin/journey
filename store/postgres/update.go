package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) UpdatePost(ctx context.Context, postID int64, title string, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, metaDescription string, image string, updatedAt time.Time, updatedBy int64) error {
	currentPost, err := s.RetrievePostByID(ctx, postID)
	if err != nil {
		return err
	}

	status := "draft"
	if published {
		status = "published"
	}

	return s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		// If the updated post is published for the first time, add publication date and user
		if published && !currentPost.IsPublished {
			const stmtUpdatePostPublished = "UPDATE posts SET title = $1, slug = $2, markdown = $3, html = $4, featured = $5, page = $6, status = $7, meta_description = $8, image = $9, updated_at = $10, updated_by = $11, published_at = $12, published_by = $13 WHERE id = $14"
			return execError(tx.Exec(ctx, stmtUpdatePostPublished, title, slug, markdown, html, featured, isPage, status, metaDescription, image, updatedAt, updatedBy, updatedAt, updatedBy, postID))
		}
		const stmtUpdatePost = "UPDATE posts SET title = $1, slug = $2, markdown = $3, html = $4, featured = $5, page = $6, status = $7, meta_description = $8, image = $9, updated_at = $10, updated_by = $11 WHERE id = $12"
		return execError(tx.Exec(ctx, stmtUpdatePost, title, slug, markdown, html, featured, isPage, status, metaDescription, image, updatedAt, updatedBy, postID))
	})
}

func (s *Storage) UpdateSettings(
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
) error {
	return s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		const stmtUpdateSettings = "UPDATE settings SET value = $1, updated_at = $2, updated_by = $3 WHERE key = $4"

		// Title
		_, err := tx.Exec(ctx, stmtUpdateSettings, title, updatedAt, updatedBy, "title")
		if err != nil {
			return fmt.Errorf("exec stmtUpdateSettings title: %w", err)
		}

		// Description
		_, err = tx.Exec(ctx, stmtUpdateSettings, description, updatedAt, updatedBy, "description")
		if err != nil {
			return fmt.Errorf("exec stmtUpdateSettings description: %w", err)
		}

		// Logo
		_, err = tx.Exec(ctx, stmtUpdateSettings, logo, updatedAt, updatedBy, "logo")
		if err != nil {
			return fmt.Errorf("exec stmtUpdateSettings logo: %w", err)
		}

		// Cover
		_, err = tx.Exec(ctx, stmtUpdateSettings, cover, updatedAt, updatedBy, "cover")
		if err != nil {
			return fmt.Errorf("exec stmtUpdateSettings cover: %w", err)
		}

		// PostsPerPage
		_, err = tx.Exec(ctx, stmtUpdateSettings, strconv.FormatInt(postsPerPage, 10), updatedAt, updatedBy, "postsPerPage")
		if err != nil {
			return fmt.Errorf("exec stmtUpdateSettings postsPerPage: %w", err)
		}

		// ActiveTheme
		_, err = tx.Exec(ctx, stmtUpdateSettings, activeTheme, updatedAt, updatedBy, "activeTheme")
		if err != nil {
			return fmt.Errorf("exec stmtUpdateSettings activeTheme: %w", err)
		}

		// Navigation
		_, err = tx.Exec(ctx, stmtUpdateSettings, navigation, updatedAt, updatedBy, "navigation")
		if err != nil {
			return fmt.Errorf("exec stmtUpdateSettings navigation: %w", err)
		}
		return nil
	})
}

func (s *Storage) UpdateActiveTheme(ctx context.Context, activeTheme string, updatedAt time.Time, updatedBy int64) error {
	return s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		const stmtUpdateSettings = "UPDATE settings SET value = $1, updated_at = $2, updated_by = $3 WHERE key = $4"
		return execError(tx.Exec(ctx, stmtUpdateSettings, activeTheme, updatedAt, updatedBy, "activeTheme"))
	})
}

func (s *Storage) UpdateUser(ctx context.Context, userID int64, name string, slug string, email string, image string, cover string, bio string, website string, location string, updatedAt time.Time, updatedBy int64) error {
	return s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		const stmtUpdateUser = "UPDATE users SET name = $1, slug = $2, email = $3, image = $4, cover = $5, bio = $6, website = $7, location = $8, updated_at = $9, updated_by = $10 WHERE id = $11"
		return execError(tx.Exec(ctx, stmtUpdateUser, name, slug, email, image, cover, bio, website, location, updatedAt, updatedBy, userID))
	})
}

func (s *Storage) UpdateLastLogin(ctx context.Context, logInDate time.Time, userID int64) error {
	const stmtUpdateLastLogin = "UPDATE users SET last_login = $1 WHERE id = $2"
	return s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		return execError(tx.Exec(ctx, stmtUpdateLastLogin, logInDate, userID))
	})
}

func (s *Storage) UpdateUserPassword(ctx context.Context, userID int64, password string, updatedAt time.Time, updatedBy int64) error {
	const stmtUpdateUserPassword = "UPDATE users SET password = $1, updated_at = $2, updated_by = $3 WHERE id = $4" // #nosec G101
	return s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		return execError(tx.Exec(ctx, stmtUpdateUserPassword, password, updatedAt, updatedBy, userID))
	})
}
