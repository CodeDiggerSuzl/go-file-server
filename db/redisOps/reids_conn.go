package redisOps

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var redisConn redis.Conn

func init() {
	var err error
	redisConn, err = redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Printf("err dailing redisOps %s", err.Error())
		panic(err.Error())
	}
	var result interface{}
	result, err = redisConn.Do("PING")
	if err != nil {
		fmt.Printf("err ping redisOps %s", err.Error())
		panic(err.Error())
	}
	if result != "PONG" {
		panic("error during ping")
	}
}
func GetRedisConn() redis.Conn {
	return redisConn
}
