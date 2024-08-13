package pkg

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewDBConn(dbDriver, dbSource string) *sql.DB {
	log.Println("dbDriver: ", dbDriver)
	log.Println("dbSource: ", dbSource)
	db, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("error establish connection db")
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}
