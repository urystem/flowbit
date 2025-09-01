package postgres

import (
	"context"
	"fmt"

	"marketflow/internal/domain"
)

func (db *poolDB) GetComments(ctx context.Context, postID uint64) ([]domain.Comment, error) {
	const query = `
		SELECT 
			c.comment_id,
			u.name,
			u.avatar_url,
			c.parent_comment_id,
			c.comment_content,
			c.has_image,
			c.comment_time
		FROM comments c
		JOIN users u ON c.user_id = u.session_id
		WHERE c.post_id = $1
		ORDER BY c.comment_time ASC`

	rows, err := db.Pool.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var c domain.Comment
		err := rows.Scan(
			&c.CommentID,
			&c.UserName,
			&c.AvatarURL,
			&c.ReplyToID,
			&c.Content,
			&c.HasImage,
			&c.DataTime,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (db *poolDB) InsertComment(ctx context.Context, comment *domain.InsertComment) (uint64, error) {
	const query = `
		INSERT INTO comments (
			post_id,
			user_id,
			comment_content,
			has_image
		)
		VALUES ($1, $2, $3, $4)
		RETURNING comment_id;`

	var commentID uint64
	err := db.QueryRow(
		ctx, query,
		comment.PostID,
		comment.User,
		comment.Content,
		comment.HasImage,
	).Scan(&commentID)
	return commentID, err
}

func (db *poolDB) DeleteComment(ctx context.Context, commentID uint64) error {
	const query = `
		DELETE FROM comments
		WHERE comment_id = $1;
	`

	cmdTag, err := db.Exec(ctx, query, commentID)
	if err != nil {
		return fmt.Errorf("DeleteComment: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteComment: no comment with id %d found", commentID)
	}

	return nil
}
