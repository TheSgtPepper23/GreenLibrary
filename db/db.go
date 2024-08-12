package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnection() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("CONN_STRING"))

	if err != nil {
		return nil, err
	}

	return dbpool, nil
}
