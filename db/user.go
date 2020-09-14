package db

import (
	"fmt"
	"go-file-server/db/mysql"
	"go-file-server/db/redisOps"

	"github.com/gomodule/redigo/redis"
)

const (
	redisTokenPrefix = "auth:user:"
)

// UserSign username and password
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

// UserSignin user sign in db ops
func UserSignin(userName, encpwd string) bool {
	stmt, err := mysql.GetDbConn().Prepare("SELECT `user_pwd` FROM tbl_user WHERE user_name = ? LIMIT 1")
	if err != nil {
		fmt.Printf("error during query user info %s", err.Error())
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(userName)
	if err != nil {
		fmt.Printf("error during query user row:%s", err.Error())
		return false
	} else if rows == nil {
		fmt.Printf("user signin got nothing")
		return false
	}
	prepareRows := mysql.PrepareRows(rows)
	fmt.Println(string(prepareRows[0]["user_pwd"].([]byte)))
	if len(prepareRows) > 0 && string(prepareRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// UpdateUserToken update user token
func UpdateUserToken(username, token string) bool {
	stmt, err := mysql.GetDbConn().Prepare("INSERT INTO tbl_user_token (`user_name`,`user_token`)values (?,?)")
	if err != nil {
		fmt.Printf("error during user user token :%s", err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Printf("err during exec: %s", err.Error())
		return false
	}
	return true
}

// User Model of user
type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

// GetUserToken get user token by name
func GetUserToken(username string) (string, error) {
	token := ""
	stmt, err := mysql.GetDbConn().Prepare("Select `user_token` from tbl_user_token where `user_name` = ? limit 1")
	if err != nil {
		fmt.Printf("error during query user token %s", err.Error())
		return token, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&token)
	if err != nil {
		return token, err
	}
	return token, nil
}

// QueryUserInfo query user info
func QueryUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mysql.GetDbConn().Prepare("select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Printf("err during query user info %s", err.Error())
		return user, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Printf("err during query user info:%s", err.Error())
		return user, err
	}
	return user, nil
}

// SetUserTokenIntoRedis set user token with expire time 30*60 seconds, aka 30 minutes
func SetUserTokenIntoRedis(username, token string) (bool, error) {
	conn := redisOps.GetRedisConn()
	resp, err := conn.Do("SET", redisTokenPrefix+username, token, "EX", "1800")
	defer conn.Close()
	if err != nil {
		fmt.Printf("error during set user token into redisOps:%s", err.Error())
		return false, err
	}
	if resp != "OK" {
		return false, nil
	}
	return true, nil
}

func GetUserTokenFromRedis(username string) (string, error) {
	conn := redisOps.GetRedisConn()
	resp, err := redis.String(conn.Do("GET", redisTokenPrefix+username))
	defer conn.Close()
	if err != nil {
		fmt.Printf("error during get user token from redisOps %s", err.Error())
		return resp, err
	}
	return resp, nil
}
