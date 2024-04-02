package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var dbCon *sql.DB

func InitDB() {
	var err error
	dbCon, err = sql.Open("mysql", "newuser:password@tcp(127.0.0.1:3306)/messenger_db?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	err = dbCon.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() *sql.DB {
	return dbCon
}
