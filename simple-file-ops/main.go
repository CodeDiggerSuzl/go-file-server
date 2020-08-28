package main

import (
	"fmt"
	"go-file-server/simple-file-ops/handler"
	"net/http"
)

func main() {
	// file upload handler
	http.HandleFunc("/file/upload/", handler.UploadHandler)
	// upload success handler
	http.HandleFunc("/file/upload/suc/", handler.UploadSucHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Fail to start server, the error is %s", err.Error())
	}
}
