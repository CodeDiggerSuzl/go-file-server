package mysql

import (
	"database/sql"
	"fmt"
	"log"
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

func PrepareRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	//
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		checkErr(err)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
