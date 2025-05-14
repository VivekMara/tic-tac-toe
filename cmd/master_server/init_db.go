package masterserver

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Create_tables(tables []string) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	db_name := os.Getenv("DB_NAME")
	db_user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")

	conn_str := fmt.Sprintf("postgres://%s:%s@127.0.0.1:59274/%s?sslmode=disable", db_user, pass, db_name)

	db, err = sql.Open("postgres", conn_str)

	if err != nil {
		return err
	}

	for i := range tables {
		_, err = db.Exec(tables[i])
		if err != nil {
			return err
		}
	}

	return nil
}