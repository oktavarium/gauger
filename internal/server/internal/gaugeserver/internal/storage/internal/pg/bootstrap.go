package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var createTablesQuery string = `
CREATE TABLE IF NOT EXISTS
gauge
(
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE,
	value DOUBLE PRECISION
);

CREATE TABLE IF NOT EXISTS
counter
(
	id SERIAL PRIMARY KEY,
 	name TEXT UNIQUE,
 	value BIGINT
);
`

func (s *storage) bootstrap(ctx context.Context) (err error) {
	fmt.Println("!")
	tx, err := s.BeginTx(ctx, pgx.TxOptions{})
	fmt.Println("!!")
	if err != nil {
		fmt.Println("!!!")
		return fmt.Errorf("error occured on opening tx: %w", err)
	}
	fmt.Println("!!!!")
	defer func() {
		if err != nil {
			fmt.Println("!!!!!")
			err = tx.Rollback(ctx)
		} else {
			fmt.Println("!!!!!!!!!")
			err = tx.Commit(ctx)
		}

	}()
	fmt.Println("!!!!!!")
	_, err = tx.Exec(ctx, createTablesQuery)
	if err != nil {
		fmt.Println("!!!!!!!")
		return fmt.Errorf("error occured on creating tables: %w", err)
	}
	fmt.Println("!!!!!!!!")
	return nil
}
