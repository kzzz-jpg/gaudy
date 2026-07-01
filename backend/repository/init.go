package repository

import (
	"database/sql"
	"fmt"
	"os"
)

func InitPostgreDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), "disable")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	//test
	_, err = db.Exec("DROP TABLE IF EXISTS guas")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS guas (
		    gua_id BIGSERIAL PRIMARY KEY,
		    title TEXT NOT NULL,
		    people TEXT[],
		    people_str TEXT,
		    content TEXT,
		    created_at TIMESTAMP DEFAULT now(),
		    search_vector tsvector GENERATED ALWAYS AS (
		        to_tsvector('simple',coalesce(title,'') || ' ' || coalesce(people_str,'') || ' ' || coalesce(content,'')) 
		    ) STORED
		)
    `)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitPostgreDBWithZhparser() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), "disable")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// DROP 會把舊表刪除 每次重新啟動資料都會消失 deploy 的時候記得把drop拿掉
	zhsql := `
		CREATE EXTENSION IF NOT EXISTS zhparser
    `
	if _, err := db.Exec(zhsql); err != nil {
		return nil, err
	}
	zhsql = `
		DROP TABLE IF EXISTS guas
    `
	if _, err := db.Exec(zhsql); err != nil {
		return nil, err
	}
	zhsql = `
		DROP TEXT SEARCH CONFIGURATION IF EXISTS zhcfg
    `
	if _, err := db.Exec(zhsql); err != nil {
		return nil, err
	}
	zhsql = `
		CREATE TEXT SEARCH CONFIGURATION zhcfg (PARSER = zhparser)
    `
	if _, err := db.Exec(zhsql); err != nil {
		return nil, err
	}
	zhsql = `
		ALTER TEXT SEARCH CONFIGURATION zhcfg ADD MAPPING FOR n,v,a,i,e,l WITH simple
    `
	if _, err := db.Exec(zhsql); err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS guas (
		    gua_id BIGSERIAL PRIMARY KEY,
		    title TEXT NOT NULL,
		    people TEXT[],
		    people_str TEXT,
		    content TEXT,
		    created_at TIMESTAMP DEFAULT now(),
		    search_vector tsvector GENERATED ALWAYS AS (
		        to_tsvector('zhcfg',coalesce(title,'') || ' ' || coalesce(people_str,'') || ' ' || coalesce(content,'')) 
		    ) STORED
		)
    `)
	if err != nil {
		return nil, err
	}
	return db, nil
}
