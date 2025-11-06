package postgres

import (
	"context"
	"fmt"
)

func (p *poolDB) CleanerTest(ctx context.Context) error {
	tx, err := p.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	const (
		delAverages = `DELETE 
		FROM exchange_averages 
		WHERE source = 'test';`
		delBackup = `DELETE 
		FROM exchange_backup 
		WHERE source = 'test';`
	)
	if _, err := tx.Exec(ctx, delAverages); err != nil {
		return fmt.Errorf("delete from exchange_averages: %w", err)
	}

	if _, err := tx.Exec(ctx, delBackup); err != nil {
		return fmt.Errorf("delete from exchange_backup: %w", err)
	}
	return tx.Commit(ctx)
}
