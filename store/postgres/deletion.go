package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) DeletePostTagsForPostID(ctx context.Context, postID int64) error {
	return s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		const stmtDeletePostTagsByPostID = "DELETE FROM posts_tags WHERE post_id = $1"
		return execError(tx.Exec(ctx, stmtDeletePostTagsByPostID, postID))
	})
}

func (s *Storage) DeletePostByID(ctx context.Context, postID int64) error {
	return s.Tx(ctx, func(ctx context.Context, tx pgx.Tx, s *Storage) error {
		const stmtDeletePostByID = "DELETE FROM posts WHERE id = $1"
		return execError(tx.Exec(ctx, stmtDeletePostByID, postID))
	})
}
