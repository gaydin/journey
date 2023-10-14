package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gaydin/journey/date"
	"github.com/gaydin/journey/store/postgres/migrate"
	"github.com/gaydin/journey/structure"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(dsn string) (*Storage, error) {
	pgxPool, err := connect(dsn)
	if err != nil {
		return nil, err
	}

	store := &Storage{
		pool: pgxPool,
	}

	if err = store.checkBlogSettings(context.Background()); err != nil {
		return nil, err
	}

	return store, nil
}

func connect(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("ParseConfig err: %w", err)
	}

	pgxPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("ConnectConfig err: %w", err)
	}

	if err := migrate.Migrate(pgxPool); err != nil {
		return nil, fmt.Errorf("migrate fail: %w", err)
	}

	return pgxPool, nil
}

// Function to check and insert any missing blog settings into the database (settings could be missing if migrating from Ghost).
func (s *Storage) checkBlogSettings(ctx context.Context) error {
	const stmtRetrieveBlog = "SELECT value FROM settings WHERE key = $1"

	var tempBlog structure.Blog
	err := s.pool.QueryRow(ctx, stmtRetrieveBlog, "title").Scan(&tempBlog.Title)
	if err != nil {
		err = s.insertSettingString(ctx, "title", "My Blog", "blog", date.GetCurrentTime())
		if err != nil {
			return err
		}
	}

	err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "description").Scan(&tempBlog.Description)
	if err != nil {
		err = s.insertSettingString(ctx, "description", "Just another Blog", "blog", date.GetCurrentTime())
		if err != nil {
			return err
		}
	}

	var email []byte
	err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "email").Scan(&email)
	if err != nil {
		err = s.insertSettingString(ctx, "email", "", "blog", date.GetCurrentTime())
		if err != nil {
			return err
		}
	}

	err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "logo").Scan(&tempBlog.Logo)
	if err != nil {
		err = s.insertSettingString(ctx, "logo", "/public/images/blog-logo.jpg", "blog", date.GetCurrentTime())
		if err != nil {
			return err
		}
	}

	err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "cover").Scan(&tempBlog.Cover)
	if err != nil {
		err = s.insertSettingString(ctx, "cover", "/public/images/blog-cover.jpg", "blog", date.GetCurrentTime())
		if err != nil {
			return err
		}
	}

	var postsPerPage string
	if err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "postsPerPage").Scan(&postsPerPage); errors.Is(err, pgx.ErrNoRows) {
		err = s.insertSettingString(ctx, "postsPerPage", "5", "blog", date.GetCurrentTime())
		if err != nil {
			return err
		}
	}

	err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "activeTheme").Scan(&tempBlog.ActiveTheme)
	if err != nil {
		err = s.insertSettingString(ctx, "activeTheme", "promenade", "theme", date.GetCurrentTime())
		if err != nil {
			return err
		}
	}

	var navigation []byte
	err = s.pool.QueryRow(ctx, stmtRetrieveBlog, "navigation").Scan(&navigation)
	if err != nil {
		err = s.insertSettingString(ctx, "navigation", "[{\"label\":\"Home\", \"url\":\"/\"}]", "blog", date.GetCurrentTime())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) Tx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx, s *Storage) error) error {
	transaction, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx %w", err)
	}

	defer func() {
		if err != nil {
			_ = transaction.Rollback(ctx)
		} else {
			err = transaction.Commit(ctx)
		}
	}()

	err = fn(ctx, transaction, s)
	return err
}

func execError(_ pgconn.CommandTag, err error) error {
	return err
}
