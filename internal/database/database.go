package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(url string) (*sql.DB, error) {

	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Fail fast if the database is not reachable

	// this context is used to set a timeout for the ping operation
	// every time we ping the database, we want to make sure it doesn't take too long
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	log.Println("Successfully connected to the database")
	return db, nil
}
