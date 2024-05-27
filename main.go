package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"log"
)

func main() {
	// Connect to the PostgreSQL database
	conn, err := pgx.Connect(context.Background(), "postgresql://user:password@localhost:5432/merge")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	// Generate a new UUID
	newUUID := uuid.New()

	// Insert the UUID into the database
	_, err = conn.Exec(context.Background(), "INSERT INTO test (u) VALUES ($1)", newUUID)
	if err != nil {
		log.Fatalf("Failed to insert UUID: %v\n", err)
	}

	log.Println("UUID inserted successfully")
}
