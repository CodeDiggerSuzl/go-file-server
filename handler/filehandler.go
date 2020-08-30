package handler

import (
	"encoding/json"
	"fmt"
	"go-file-server/meta"
	"go-file-server/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	// FILEHASH file hash value.
	FILEHASH = "filehash"
	// FILENAME file name
	FILENAME = "filename"
	// OPERATION operation
	OPERATION = "op"
	// OpTypeRename operation type
	OpTypeRename = "0"
)

// UploadHandler file upload handler.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Return upload page
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			fmt.Printf("Error during ioutil.ReadFile %s", err.Error())
			_, _ = io.WriteString(w, "Internal error ")
			return
		}
		_, _ = w.Write(data)
		// Get file upload stream
	} else if r.Method == http.MethodPost {
		// Get the file and save to local dir
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("File to parse file from request: %s", err.Error())
			return
		}
		defer file.Close()
		currFileDir, _ := os.Getwd()
		fileLocation := currFileDir + "/tmp/" + header.Filename
		// Generate the file meta info
		fileMeta := meta.Filemeta{
			FileName: header.Filename,
			Location: fileLocation,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"), // TODO Special time in Golang
		}

		newFile, err := os.Create(fileLocation)
		if err != nil {
			fmt.Printf("Error during os.create file: %s", err.Error())
			return
		}
		defer newFile.Close()
		// Set file size
		fileMeta.FileSize, err = io.Copy(newFile, file)

		if err != nil {
			fmt.Printf("Error during io.Copy file: %s", err.Error())
			return
		}
		_, _ = newFile.Seek(0, 0) // TODO seek method ?
		// Set fileSha1
		fileMeta.FileSha1 = utils.FileSha1(newFile)
		// Add fileMeta info the the fileMetaMap
		// meta.UpdateFileMeta(fileMeta)
		// change to save file Meta to Db rather than memory
		meta.UpdateFileMetaDb(fileMeta)
		// Redirect to new page
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSucHandler upload file.
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler get meta data of a file.
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	// returns a array and get the first one
	fHash := r.Form[FILEHASH][0]
	// fMeta := meta.GetFileMeta(fHash)
	// getFromDb
	fMeta, err := meta.GetFileMetaFromDb(fHash)

	// Encoding to json
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// FileQueryHandler query file by file hash
// func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm() // TODO
// 	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))

// }

// DownloadByFileHashHandler down load file by it's hash.
func DownloadByFileHashHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	fHash := r.Form.Get(FILEHASH)
	fMeta := meta.GetFileMeta(fHash)
	file, err := os.Open(fMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	// TODO if file is too huge, should use segment stream way to read a bit file and continue to read
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Add header to make the browser to download the file
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Basics_of_HTTP/MIME_types
	w.Header().Set("content-type", "application/octect-stream")
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Disposition
	w.Header().Set("content-disposition", "attachment; filename=\""+fMeta.FileName+"\"")
	_, _ = w.Write(data)
}

// UpdateMetaInfoHandle change filename by file hash.
func UpdateMetaInfoHandle(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	opType := r.Form.Get(OPERATION)
	//  0 Means: rename
	if opType != OpTypeRename {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	filehash := r.Form.Get(FILEHASH)
	newFName := r.Form.Get(FILENAME)

	currFileMeta := meta.GetFileMeta(filehash)
	currFileMeta.FileName = newFName
	meta.UpdateFileMeta(currFileMeta)
	data, err := json.Marshal(currFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// DelFileHandler delete file handler.
func DelFileHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	fileSha1 := r.Form.Get(FILEHASH)
	fileMeta := meta.GetFileMeta(fileSha1)
	// Del from disk FIRST
	err := os.Remove(fileMeta.Location)
	if err != nil {
		fmt.Printf("Err during del the file from disk: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Del the fileMetaInfo map
	meta.DelFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
}
