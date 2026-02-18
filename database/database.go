package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	// Build connection URL
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://localhost:5432/hiimtaylor_go_development?sslmode=disable"
	}

	var err error 
	// Attempt connection
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}

	// Verify connection
	err = Pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("Connected to database")
}

func Close() {
	Pool.Close()
}