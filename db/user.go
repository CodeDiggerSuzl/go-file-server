package db

import (
	"fmt"
	"go-file-server/db/mysql"
)

// UserSign username and passwd
func UserSign(username, password string) bool {
	stmt, err := mysql.GetDbConn().Prepare("INSERT ignore into tbl_user(`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Printf("error during insert user %s", err.Error())
		return false
	}
	defer stmt.Close()
	result, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Printf("error during exec insert user %s", err.Error())
		return false
	}
	if rowsAffected, err := result.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}
