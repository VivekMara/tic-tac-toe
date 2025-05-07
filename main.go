package main

import (
	"log"
	"net/http"
	"os"
	"tic-tac-toe/cmd/helpers"
	masterserver "tic-tac-toe/cmd/master_server"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// DB initialization
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	stmt := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	database, err := helpers.New(dbUser, dbPass, port, dbName, stmt)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Conn.Close()

	log.Println("Database connection established and table created successfully")

	//routes
	http.HandleFunc("/", helpers.Root)
	http.HandleFunc("/auth/v1/register", masterserver.Register)

	//server
	log.Println("Server starting on port 6969")
	log.Fatal(http.ListenAndServe(":6969", nil))
}
