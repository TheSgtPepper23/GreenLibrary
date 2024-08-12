package db

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestSQLCreateCollection(t *testing.T) {

	conn, err := pgxpool.New(context.Background(), "postgres://jim:JimGreen16@localhost:5432/GreenLibrary?sslmode=disable&pool_max_conns=10")

	if err != nil {
		t.Log(err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	sqlCollection := NewSQLCollectionContext(conn)
	if err != nil {
		t.Log(err.Error())
		os.Exit(1)
	}

	_, err = sqlCollection.SQLRetrieveCollections()

	if err != nil {
		t.Log(err.Error())
		t.Error(err)
	}
}
