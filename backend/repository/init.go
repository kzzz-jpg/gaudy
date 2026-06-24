package repository

import (
	"database/sql"
	"fmt"
)

func Initdb() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", "localhost", "postgres", "cchs91193", "guadb", "disable")
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
