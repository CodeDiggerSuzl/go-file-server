package redisOps

import "testing"

func Test_getRedisConn(t *testing.T) {
	conn := GetRedisConn()
	if conn == nil {
		t.Fail()
	}
}
