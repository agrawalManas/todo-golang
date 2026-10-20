package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(dbUrl string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	var config *pgxpool.Config
	var err error

	config, err = pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Println("Could not connect to database", err)
		return nil, err
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Println("Could not setup a new connection to DB:", err)
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Println("Could not establish a stable connection to DB:", err)
		pool.Close()
		return nil, err
	}

	return pool, nil
}
