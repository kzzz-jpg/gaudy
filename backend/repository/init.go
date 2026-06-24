package repository

import (
	"database/sql"
	"fmt"
	"os"
)

func Initdb() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), "disable")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS guas (
		    gua_id BIGSERIAL PRIMARY KEY,
		    title TEXT NOT NULL,
		    people TEXT[],
		    content TEXT,
		    created_at TIMESTAMP DEFAULT now()
		)
    `)
	if err != nil {
		return nil, err
	}
	return db, nil
}
