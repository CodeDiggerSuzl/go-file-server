package db

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Test_getUserToken(t *testing.T) {
	userName := "chanman"
	token, err := GetUserToken(userName)
	if token != "509b90ce5cd23b6238db46ac1b7229a85f5c5057" || err != nil {
		log.Fatal("got the wrong thing")
	}
}
func Test_timeStamp(t *testing.T) {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	t.Logf("%T", ts)
	t.Log(ts)
}

func Test_setTokenIntoRedis(t *testing.T) {
	ok, err := SetUserTokenIntoRedis("rand_user_name", "rand_value")
	if !ok {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}

func Test_getUserTokenFromRedis(t *testing.T) {
	tokenFromRedis, err := GetUserTokenFromRedis("no_exist_name")
	if tokenFromRedis == "" {
		t.Fail()
	}
	if err != nil {
		log.Print(err)
		t.Fail()
	}
}
