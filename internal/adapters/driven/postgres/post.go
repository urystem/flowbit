package postgres

import (
	"context"
	"fmt"

	"1337b04rd/internal/domain"
)

func (db *poolDB) SelectActivePosts(ctx context.Context) ([]domain.PostNonContent, error) {
	const query string = `
	SELECT post_id, title, has_image
		FROM posts
		WHERE archived = FALSE
		ORDER BY post_time DESC`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.PostNonContent
	for rows.Next() {
		var p domain.PostNonContent

		if err := rows.Scan(&p.ID, &p.Title, &p.HasImage); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return posts, nil
}

func (db *poolDB) InsertPost(ctx context.Context, post *domain.InsertPost) (uint64, error) {
	const query string = `
		INSERT INTO posts (user_id, author, title, post_content, has_image)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING post_id;`
	var postID uint64
	err := db.Pool.QueryRow(ctx, query,
		post.Uuid,
		post.Name,
		post.Subject,
		post.Content,
		post.HasImage).Scan(&postID)

	return postID, err
}

func (db *poolDB) DeletePost(ctx context.Context, id uint64) error {
	const query string = `DELETE FROM posts WHERE post_id = $1`

	// Выполняем запрос
	ct, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	// Проверим, была ли строка удалена
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("delete post: no post with id %d", id)
	}

	return nil
}

func (db *poolDB) SelectArchivePosts(ctx context.Context) ([]domain.PostNonContent, error) {
	const query string = `
	SELECT post_id, title, has_image
		FROM posts
		WHERE archived = TRUE
		ORDER BY post_time DESC`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.PostNonContent
	for rows.Next() {
		var p domain.PostNonContent

		if err := rows.Scan(&p.ID, &p.Title, &p.HasImage); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return posts, nil
}

func (db *poolDB) GetPost(ctx context.Context, id uint64) (*domain.PostX, error) {
	const query = `
        SELECT 
            title, post_content, has_image
        FROM 
            posts
        WHERE 
            post_id = $1`
	
	postX := new(domain.PostX)
	
	err := db.QueryRow(ctx, query, id).Scan(&postX.Title, &postX.Content, &postX.HasImage)
	if err != nil {
		return nil, err
	}
	
	return postX, nil
}
