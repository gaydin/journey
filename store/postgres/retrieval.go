package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/gaydin/journey/structure"
)

func (s *Storage) RetrievePostByID(ctx context.Context, postID int64) (*structure.Post, error) {
	const stmtRetrievePostByID = `	
		SELECT id,
			   uuid,
			   title,
			   slug,
			   markdown,
			   html,
			   featured,
			   page,
			   status,
			   meta_description,
			   image,
			   author_id,
			   published_at
		FROM posts
		WHERE id = $1`

	row := s.pool.QueryRow(ctx, stmtRetrievePostByID, postID)
	return s.extractPost(ctx, row)
}

func (s *Storage) RetrievePostBySlug(ctx context.Context, slug string) (*structure.Post, error) {
	const stmtRetrievePostBySlug = `
		SELECT id,
			   uuid,
			   title,
			   slug,
			   markdown,
			   html,
			   featured,
			   page,
			   status,
			   meta_description,
			   image,
			   author_id,
			   published_at
		FROM posts
		WHERE slug = $1`

	row := s.pool.QueryRow(ctx, stmtRetrievePostBySlug, slug)

	return s.extractPost(ctx, row)
}

func (s *Storage) RetrievePostsByUser(ctx context.Context, userID int64, limit int64, offset int64) ([]structure.Post, error) {
	const stmtRetrievePostsByUser = `
		SELECT id,
			   uuid,
			   title,
			   slug,
			   markdown,
			   html,
			   featured,
			   page,
			   status,
			   meta_description,
			   image,
			   author_id,
			   published_at
		FROM posts
		WHERE page = false
		  AND status = 'published'
		  AND author_id = $1
		ORDER BY published_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.pool.Query(ctx, stmtRetrievePostsByUser, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts, err := s.extractPosts(ctx, rows)
	if err != nil {
		return nil, fmt.Errorf("extractPosts err: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err %w", err)
	}
	return *posts, nil
}

func (s *Storage) RetrievePostsByTag(ctx context.Context, tagID int64, limit int64, offset int64) ([]structure.Post, error) {
	const stmtRetrievePostsByTag = `
		SELECT posts.id,
			   posts.uuid,
			   posts.title,
			   posts.slug,
			   posts.markdown,
			   posts.html,
			   posts.featured,
			   posts.page,
			   posts.status,
			   posts.meta_description,
			   posts.image,
			   posts.author_id,
			   posts.published_at
		FROM posts,
			 posts_tags
		WHERE posts_tags.post_id = posts.id
		  AND posts_tags.tag_id = $1
		  AND page = false
		  AND status = 'published'
		ORDER BY posts.published_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.pool.Query(ctx, stmtRetrievePostsByTag, tagID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts, err := s.extractPosts(ctx, rows)
	if err != nil {
		return nil, fmt.Errorf("extractPosts err: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err %w", err)
	}

	return *posts, nil
}

func (s *Storage) RetrievePostsForIndex(ctx context.Context, limit int64, offset int64) ([]structure.Post, error) {
	const stmtRetrievePostsForIndex = `
		SELECT id,
			   uuid,
			   title,
			   slug,
			   markdown,
			   html,
			   featured,
			   page,
			   status,
			   meta_description,
			   image,
			   author_id,
			   published_at
		FROM posts
		WHERE page = false
		  AND status = 'published'
		ORDER BY published_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := s.pool.Query(ctx, stmtRetrievePostsForIndex, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts, err := s.extractPosts(ctx, rows)
	if err != nil {
		return nil, fmt.Errorf("extractPosts err: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err %w", err)
	}

	return *posts, nil
}

func (s *Storage) RetrievePostsForApi(ctx context.Context, limit int64, offset int64) ([]structure.Post, error) {
	const stmtRetrievePostsForAPI = `
		SELECT id,
			   uuid,
			   title,
			   slug,
			   markdown,
			   html,
			   featured,
			   page,
			   status,
			   meta_description,
			   image,
			   author_id,
			   published_at
		FROM posts
		ORDER BY id DESC
		LIMIT $1 OFFSET $2`

	rows, err := s.pool.Query(ctx, stmtRetrievePostsForAPI, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts, err := s.extractPosts(ctx, rows)
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err %w", err)
	}

	return *posts, nil
}

func (s *Storage) extractPosts(ctx context.Context, rows pgx.Rows) (*[]structure.Post, error) {
	posts := make([]structure.Post, 0)

	for rows.Next() {
		var (
			userID int64
			status string
			post   structure.Post
		)

		err := rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.Slug, &post.Markdown, &post.Html, &post.IsFeatured, &post.IsPage, &status, &post.MetaDescription, &post.Image, &userID, &post.Date)
		if err != nil {
			return nil, err
		}

		// If there was no publication date attached to the post, make its creation date the date of the post
		if post.Date == nil {
			post.Date, err = s.retrievePostCreationDateByID(ctx, post.Id)
			if err != nil {
				return nil, err
			}
		}

		// Evaluate status
		if status == "published" {
			post.IsPublished = true
		} else {
			post.IsPublished = false
		}

		// Retrieve user
		post.Author, err = s.RetrieveUser(ctx, userID)
		if err != nil {
			return nil, err
		}

		// Retrieve tags
		post.Tags, err = s.RetrieveTags(ctx, post.Id)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return &posts, nil
}

func (s *Storage) extractPost(ctx context.Context, row pgx.Row) (*structure.Post, error) {
	var (
		post   structure.Post
		userID int64
		status string
	)

	err := row.Scan(&post.Id, &post.Uuid, &post.Title, &post.Slug, &post.Markdown, &post.Html, &post.IsFeatured, &post.IsPage, &status, &post.MetaDescription, &post.Image, &userID, &post.Date)
	if err != nil {
		return nil, err
	}

	// If there was no publication date attached to the post, make its creation date the date of the post
	if post.Date == nil {
		post.Date, err = s.retrievePostCreationDateByID(ctx, post.Id)
		if err != nil {
			return nil, err
		}
	}

	// Evaluate status
	if status == "published" {
		post.IsPublished = true
	} else {
		post.IsPublished = false
	}

	// Retrieve user
	post.Author, err = s.RetrieveUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Retrieve tags
	post.Tags, err = s.RetrieveTags(ctx, post.Id)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *Storage) RetrieveNumberOfPosts(ctx context.Context) (int64, error) {
	const stmtRetrievePostsCount = `
		SELECT count(*)
		FROM posts
		WHERE page = false
			AND status = 'published'`

	var count int64
	if err := s.pool.QueryRow(ctx, stmtRetrievePostsCount).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Storage) RetrieveNumberOfPostsByUser(ctx context.Context, userID int64) (int64, error) {
	const stmtRetrievePostsCountByUser = `	
		SELECT count(*)
		FROM posts
		WHERE page = false
			AND status = 'published'
			AND author_id = $1`

	var count int64
	if err := s.pool.QueryRow(ctx, stmtRetrievePostsCountByUser, userID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Storage) RetrieveNumberOfPostsByTag(ctx context.Context, tagID int64) (int64, error) {
	const stmtRetrievePostsCountByTag = `
		SELECT count(*)
		FROM posts,
			 posts_tags
		WHERE posts_tags.post_id = posts.id
		  AND posts_tags.tag_id = $1
		  AND page = false
		  AND status = 'published'`

	var count int64
	if err := s.pool.QueryRow(ctx, stmtRetrievePostsCountByTag, tagID).Scan(&count); err != nil {
		return 0, fmt.Errorf("stmtRetrievePostsCountByTag %w", err)
	}

	return count, nil
}

func (s *Storage) retrievePostCreationDateByID(ctx context.Context, postID int64) (*time.Time, error) {
	const stmtRetrievePostCreationDateByID = `
		SELECT created_at
		FROM posts
		WHERE id = $1`

	var creationDate time.Time
	if err := s.pool.QueryRow(ctx, stmtRetrievePostCreationDateByID, postID).Scan(&creationDate); err != nil {
		return &creationDate, fmt.Errorf("stmtRetrievePostCreationDateByID %w", err)
	}

	return &creationDate, nil
}

func (s *Storage) RetrieveUser(ctx context.Context, userID int64) (*structure.User, error) {
	const stmtRetrieveUserByID = `
		SELECT id,
			   name,
			   slug,
			   email,
			   image,
			   cover,
			   bio,
			   website,
			   location
		FROM users
		WHERE id = $1`

	var (
		user     structure.User
		image    sql.NullString
		cover    sql.NullString
		bio      sql.NullString
		website  sql.NullString
		location sql.NullString
	)

	err := s.pool.QueryRow(ctx, stmtRetrieveUserByID, userID).Scan(
		&user.Id,
		&user.Name,
		&user.Slug,
		&user.Email,
		&image,
		&cover,
		&bio,
		&website,
		&location,
	)

	user.Image = image.String
	user.Cover = cover.String
	user.Bio = bio.String
	user.Website = website.String
	user.Location = location.String

	if err != nil {
		return nil, fmt.Errorf("stmtRetrieveUserByIDyt %w", err)
	}

	return &user, nil
}

func (s *Storage) RetrieveUserBySlug(ctx context.Context, slug string) (*structure.User, error) {
	const stmtRetrieveUserBySlug = `
		SELECT id,
			   name,
			   slug,
			   email,
			   image,
			   cover,
			   bio,
			   website,
			   location
		FROM users
		WHERE slug = $1`

	var user structure.User
	err := s.pool.QueryRow(ctx, stmtRetrieveUserBySlug, slug).Scan(
		&user.Id,
		&user.Name,
		&user.Slug,
		&user.Email,
		&user.Image,
		&user.Cover,
		&user.Bio,
		&user.Website,
		&user.Location,
	)

	if err != nil {
		return nil, fmt.Errorf("stmtRetrieveUserBySlug scan %w", err)
	}

	return &user, nil
}

// RetrieveUserByName Retrieve user by name
func (s *Storage) RetrieveUserByName(ctx context.Context, name string) (*structure.User, error) {
	const stmtRetrieveUserByName = `
		SELECT id,
			   name,
			   slug,
			   email,
			   image,
			   cover,
			   bio,
			   website,
			   location
		FROM users
		WHERE name = $1`

	var (
		user     structure.User
		image    sql.NullString
		cover    sql.NullString
		bio      sql.NullString
		website  sql.NullString
		location sql.NullString
	)

	err := s.pool.QueryRow(ctx, stmtRetrieveUserByName, name).Scan(
		&user.Id,
		&user.Name,
		&user.Slug,
		&user.Email,
		&image,
		&cover,
		&bio,
		&website,
		&location,
	)

	user.Image = image.String
	user.Cover = cover.String
	user.Bio = bio.String
	user.Website = website.String
	user.Location = location.String

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) RetrieveTags(ctx context.Context, postID int64) ([]structure.Tag, error) {
	const stmtRetrieveTags = `
		SELECT tag_id
		FROM posts_tags
		WHERE post_id = $1`

	rows, err := s.pool.Query(ctx, stmtRetrieveTags, postID)
	if err != nil {
		return nil, fmt.Errorf("query stmtRetrieveTags %w", err)
	}
	defer rows.Close()

	var tags []structure.Tag
	for rows.Next() {
		var tagID int64
		err := rows.Scan(&tagID)
		if err != nil {
			return nil, err
		}

		tag, err := s.RetrieveTag(ctx, tagID)
		// TODO: Error while receiving individual tag is ignored right now. Keep it this way?
		if err == nil {
			tags = append(tags, *tag)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err %w", err)
	}

	return tags, nil
}

func (s *Storage) RetrieveTag(ctx context.Context, tagID int64) (*structure.Tag, error) {
	const stmtRetrieveTagByID = `
		SELECT id, name, slug
		FROM tags
		WHERE id = $1`

	var tag structure.Tag
	err := s.pool.QueryRow(ctx, stmtRetrieveTagByID, tagID).Scan(&tag.Id, &tag.Name, &tag.Slug)
	if err != nil {
		return nil, fmt.Errorf("stmtRetrieveTagByID %w", err)
	}

	return &tag, nil
}

func (s *Storage) RetrieveTagBySlug(ctx context.Context, slug string) (*structure.Tag, error) {
	const stmtRetrieveTagBySlug = `
		SELECT id, name, slug
		FROM tags
		WHERE slug = $1`

	var tag structure.Tag
	err := s.pool.QueryRow(ctx, stmtRetrieveTagBySlug, slug).Scan(&tag.Id, &tag.Name, &tag.Slug)
	if err != nil {
		return nil, fmt.Errorf("stmtRetrieveTagBySlug %w", err)
	}

	return &tag, nil
}

func (s *Storage) RetrieveTagIDBySlug(ctx context.Context, slug string) (int64, error) {
	const stmtRetrieveTagIDBySlug = `
		SELECT id
		FROM tags
		WHERE slug = $1`

	var id int64
	if err := s.pool.QueryRow(ctx, stmtRetrieveTagIDBySlug, slug).Scan(&id); err != nil {
		return 0, fmt.Errorf("stmtRetrieveTagIDBySlug %w", err)
	}

	return id, nil
}

func (s *Storage) RetrieveHashedPasswordForUser(ctx context.Context, name string) ([]byte, error) {
	const stmtRetrieveHashedPasswordByName = `
		SELECT password
		FROM users
		WHERE name = $1`

	var hashedPassword []byte
	err := s.pool.QueryRow(ctx, stmtRetrieveHashedPasswordByName, name).Scan(&hashedPassword)
	if err != nil {
		return []byte{}, fmt.Errorf("stmtRetrieveHashedPasswordByName %w", err)
	}

	return hashedPassword, nil
}

func (s *Storage) RetrieveBlog(ctx context.Context) (*structure.Blog, error) {
	const stmtRetrieveBlog = `
		SELECT value
		FROM settings
		WHERE key = $1`

	var tempBlog structure.Blog

	err := s.pool.QueryRow(ctx, stmtRetrieveBlog, "title").Scan(&tempBlog.Title)
	if err != nil {
		return &tempBlog, err
	}

	err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "description").Scan(&tempBlog.Description)
	if err != nil {
		return &tempBlog, err
	}

	err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "logo").Scan(&tempBlog.Logo)
	if err != nil {
		return &tempBlog, err
	}

	err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "cover").Scan(&tempBlog.Cover)
	if err != nil {
		return &tempBlog, err
	}

	var postsPerPage string
	if err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "postsPerPage").Scan(&postsPerPage); err != nil {
		return &tempBlog, err
	}

	ppp, err := strconv.ParseInt(postsPerPage, 10, 0)
	if err != nil {
		return &tempBlog, err
	}
	tempBlog.PostsPerPage = ppp

	err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "activeTheme").Scan(&tempBlog.ActiveTheme)
	if err != nil {
		return &tempBlog, err
	}

	postCount, err := s.RetrieveNumberOfPosts(ctx)
	if err != nil {
		return &tempBlog, err
	}

	tempBlog.PostCount = postCount

	var navigation []byte
	err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "navigation").Scan(&navigation)
	if err != nil {
		return &tempBlog, err
	}

	tempBlog.NavigationItems, err = makeNavigation(navigation)
	if err != nil {
		return &tempBlog, err
	}

	return &tempBlog, err
}

func (s *Storage) RetrieveActiveTheme(ctx context.Context) (string, error) {
	const stmtRetrieveBlog = `
		SELECT value
		FROM settings
		WHERE key = $1`

	var activeTheme string
	if err := s.pool.QueryRow(ctx, stmtRetrieveBlog, "activeTheme").Scan(&activeTheme); err != nil {
		return "", err
	}

	return activeTheme, nil
}

func (s *Storage) RetrieveUsersCount(ctx context.Context) int {
	const stmtRetrieveUsersCount = `
		SELECT count(*)
		FROM users`

	var userCount int
	if err := s.pool.QueryRow(ctx, stmtRetrieveUsersCount).Scan(&userCount); err != nil {
		return -1
	}

	return userCount
}

func makeNavigation(navigation []byte) ([]structure.Navigation, error) {
	var navigationItems []structure.Navigation

	if err := json.Unmarshal(navigation, &navigationItems); err != nil {
		return navigationItems, err
	}

	return navigationItems, nil
}
