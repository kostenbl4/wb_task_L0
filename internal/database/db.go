package database

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error){
	
	db, err := sql.Open("postgres", addr)
	
	if err != nil{
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)


	if err := db.Ping(); err != nil{
		return nil, err
	}

	return db, nil
}