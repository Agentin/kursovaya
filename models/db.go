package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "trainer"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "secure_password"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "phishing_training_portal"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}
	log.Println("Connected to PostgreSQL")
	return nil
}
