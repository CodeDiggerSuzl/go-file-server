package main

import (
	"fmt"
	"go-file-server/handler"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	// static file server
	pwd, _ := os.Getwd()
	http.Handle("/static/", http.FileServer(http.Dir(filepath.Join(pwd, "./"))))
	// File upload handler.
	http.HandleFunc("/file/upload", handler.UploadHandler)
	// Upload success handler.
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	// Get file meta info.
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadByFileHashHandler)
	http.HandleFunc("/file/update", handler.UpdateMetaInfoHandle)
	http.HandleFunc("/file/del", handler.DelFileHandler)

	http.HandleFunc("/user/signup", handler.SignUpHandler)
	http.HandleFunc("/user/signin", handler.UserSignInHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Fail to start server, the error is: %s", err.Error())
	}
}
