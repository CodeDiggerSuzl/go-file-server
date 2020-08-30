package mysql

import (
	"database/sql"
	"fmt"
	"os"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB

func init() {
	var err error
	dbConn, err = sql.Open("mysql", "root:mysql@suz1@tcp(127.0.0.1:3306)/db_file_server")
	dbConn.SetMaxOpenConns(1000)
	if err != nil {
		fmt.Printf("DB conn error %s", err.Error())
		panic(err.Error())
	}
	err = dbConn.Ping()
	if err != nil {
		fmt.Printf("DB ping error %s", err.Error())
		os.Exit(1)
	}
}

// GetDbConn Get db connection
func GetDbConn() *sql.DB {
	return dbConn
}
