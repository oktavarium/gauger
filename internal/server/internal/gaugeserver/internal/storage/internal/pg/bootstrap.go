package pg

import (
	"context"
	"fmt"
)

var createTablesQuery string = `
CREATE TABLE IF NOT EXISTS
gauge
(
	id SERIAL PRIMARY KEY,
	name TEXT,
	value DOUBLE PRECISION
);

CREATE TABLE IF NOT EXISTS
counter
(
	id SERIAL PRIMARY KEY,
 	name TEXT,
 	value INTEGER
);
`

func (s *storage) bootstrap(ctx context.Context) error {
	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error occured on opening tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, createTablesQuery)
	if err != nil {
		return fmt.Errorf("error occured on creating tables: %w", err)
	}

	return tx.Commit()
}
