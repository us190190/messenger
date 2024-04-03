package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var dbCon *sql.DB

func InitDB() error {
	var err error
	//dbCon, err = sql.Open("mysql", "newuser:password@tcp(127.0.0.1:3306)/messenger_db?parseTime=true")
	dbCon, err = sql.Open("mysql", "app-user::g$G}V-Fus-nVE{S@tcp(35.232.174.70:3306)/messenger_db?parseTime=true")
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = dbCon.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func GetDB() *sql.DB {
	return dbCon
}
