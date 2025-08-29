package postgres

import "context"

func (db *poolDB) Archiver(ctx context.Context) error {
	const sql = `
UPDATE posts
SET archived = true
WHERE archived = false
  AND (
    (
      -- нет комментариев, прошло более 10 минут с момента публикации
      NOT EXISTS (
        SELECT 1 FROM comments WHERE comments.post_id = posts.post_id
      )
      AND post_time < NOW() - interval '10 minutes'
    )
    OR
    (
      -- есть комментарии, и последний был более 15 минут назад
      EXISTS (
        SELECT 1 FROM comments WHERE comments.post_id = posts.post_id
      )
      AND (
        SELECT MAX(comment_time)
        FROM comments
        WHERE comments.post_id = posts.post_id
      ) < NOW() - interval '15 minutes'
    )
  );`
	_, err := db.Exec(ctx, sql)

	return err
}
