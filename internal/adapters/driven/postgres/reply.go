package postgres

// import (
// 	"context"
// 	"fmt"

// 	"marketflow/internal/domain"

// 	"github.com/jackc/pgx/v5"
// )

// func (db *poolDB) InsertReply(ctx context.Context, reply *domain.InsertReply) (uint64, error) {
// 	const query = `
// 		WITH parent AS (
// 			SELECT post_id
// 			FROM comments
// 			WHERE comment_id = $1
// 		)
// 		INSERT INTO comments (
// 			post_id,
// 			parent_comment_id,
// 			user_id,
// 			comment_content,
// 			has_image
// 		)
// 		SELECT
// 			post_id,
// 			$1,
// 			$2,
// 			$3,
// 			$4
// 		FROM parent
// 		RETURNING comment_id;
// 	`
// 	var commentID uint64
// 	err := db.QueryRow(
// 		ctx, query,
// 		reply.ReplyToID,
// 		reply.User,
// 		reply.Content,
// 		reply.HasImage,
// 	).Scan(&commentID)

// 	if err == pgx.ErrNoRows {
// 		return 0, fmt.Errorf("parent comment %d not found", reply.ReplyToID)
// 	}
// 	if err != nil {
// 		return 0, err
// 	}

// 	return commentID, nil
// }
