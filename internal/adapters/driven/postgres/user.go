package postgres

import (
	"context"
	"fmt"

	"marketflow/internal/domain"

	"github.com/google/uuid"
)

func (db *poolDB) InsertUser(ctx context.Context, ses *domain.Session) error {
	const query string = `
	INSERT INTO users (session_id, name, avatar_url)
	VALUES ($1, $2, $3)`
	// ON CONFLICT (session_id) DO NOTHING;

	_, err := db.Pool.Exec(ctx, query, ses.Uuid, ses.Name, ses.AvatarURL)
	return err
}

func (db *poolDB) DeleteUser(ctx context.Context, sessionID uuid.UUID) error {
	const query = `DELETE FROM users WHERE session_id = $1;`

	ct, err := db.Pool.Exec(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("delete user: no user with session_id %s", sessionID)
	}

	return nil
}
