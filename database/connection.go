package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var dbCon *sql.DB

const (
	maxRetries    = 50
	retryInterval = 20 * time.Second
)

func InitDB() error {

	var err error

	dbConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/messenger_db?parseTime=true",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	// Retry loop to establish database connection
	for i := 1; i <= maxRetries; i++ {
		//dbCon, err = sql.Open("mysql", "newuser:password@tcp(127.0.0.1:3306)/messenger_db?parseTime=true")
		dbCon, err = sql.Open("mysql", dbConnStr)
		if err != nil {
			fmt.Printf("Attempt %d: Error connecting to database: %v\n", i, err)
			time.Sleep(retryInterval)
			continue
		}

		err = dbCon.Ping()
		if err != nil {
			fmt.Printf("Attempt %d: Error pinging database: %v\n", i, err)
			dbCon.Close()
			time.Sleep(retryInterval)
			continue
		}

		fmt.Println("Database connection established successfully!")
		break
	}

	return nil
}

func GetDB() *sql.DB {
	return dbCon
}
