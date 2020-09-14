package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func (resp *RespMsg) JSONBytes() []byte {
	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return r
}

// Obj to json string
func (resp *RespMsg) JSONString() string {
	return string(resp.JSONBytes())
}

// GenSimpleRespStream return byte with code and mes only
func GenSimpleRespStream(code int, msg string) []byte {
	return []byte(fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg))
}

// GenSimpleRespString return string with code and msg
func GenSimpleRespString(code int, msg string) string {
	return fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg)
}
