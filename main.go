package main

import (
	"log"
	"fmt"
	"net/http"
	masterserver "tic-tac-toe/cmd/master_server"
)

func main() {
	//table creation and init
	tables := make([]string,0)
	users_table := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	tables = append(tables, users_table)
	sessions_table := `
	CREATE TABLE IF NOT EXISTS sessions (
    	session_id UUID PRIMARY KEY,
    	user_id INTEGER REFERENCES users(id),
    	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    	expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    	user_agent TEXT,
    	ip_address TEXT
	);
	`
	tables = append(tables, sessions_table)
	err := masterserver.Create_tables(tables)
	if err != nil {
		err_stmt := fmt.Sprintf("Failed to create tables!! with error: %s", err)
		log.Println(err_stmt)
	} else {
		log.Println("Database connection established and table created successfully")
	}

	//routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := []byte("Hello Darthman")
		w.Write(resp)
	})
	
	//auth
	http.HandleFunc("/auth/register", masterserver.Register)
	http.HandleFunc("/auth/login", masterserver.Login)
	http.HandleFunc("/auth/logout", masterserver.Logout)
	
	//server
	log.Println("Server starting on port 6969")
	log.Fatal(http.ListenAndServe(":6969", nil))
}
