package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// UploadHandler file upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// return upload page
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./simple-file-ops/static/view/index.html")
		if err != nil {
			fmt.Printf("Error during ioutil.ReadFile %s", err.Error())
			_, _ = io.WriteString(w, "Internal error ")
			return
		}
		_, _ = w.Write(data)
		// get file upload stream
	} else if r.Method == http.MethodPost {
		// get the file and save to local dir
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("File to parse file from request: %s", err.Error())
			return
		}
		defer file.Close()
		currFileDir, _ := os.Getwd()
		newFile, err := os.Create(currFileDir + "/tmp/" + header.Filename)
		if err != nil {
			fmt.Printf("Error during os.create file: %s", err.Error())
			return
		}
		defer newFile.Close()
		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Error during io.Copy file: %s", err.Error())
			return
		}
		http.Redirect(w, r, "/file/upload/suc/", http.StatusFound)
	}
}

// UploadSucHandler upload file
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Upload finished!")
}
