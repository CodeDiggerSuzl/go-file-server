package meta

import (
	"go-file-server/db"
)

// Filemeta meta info of file data
type Filemeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetaMap map[string]Filemeta

// init fileMetaMap, this function will run when this package is load.
func init() {
	fileMetaMap = make(map[string]Filemeta)
}

// UpdateFileMeta update the file map
func UpdateFileMeta(fmeta Filemeta) {
	fileMetaMap[fmeta.FileSha1] = fmeta
}

// GetFileMeta get file meta info by the sha1 value
func GetFileMeta(fileSha1 string) Filemeta {
	return fileMetaMap[fileSha1]
}

// DelFileMeta delete file by fileSha1
func DelFileMeta(fileSha1 string) {
	delete(fileMetaMap, fileSha1)
}

// UpdateFileMetaDb upload file info to db
func UpdateFileMetaDb(fileMeta Filemeta) bool {
	return db.OnFileUploadFinished(fileMeta.FileSha1, fileMeta.FileName, fileMeta.Location, fileMeta.FileSize)
}

func GetFileMetaFromDb(fileSha1 string) (*Filemeta, error) {
	tFile, err := db.GetFileMetaFromDb(fileSha1)
	if tFile == nil || err != nil {
		return nil, err
	}
	fMeta := Filemeta{
		FileSha1: tFile.FileHash,
		FileName: tFile.FileName.String,
		FileSize: tFile.FileSize.Int64,
		Location: tFile.FileAddr.String,
	}
	return &fMeta, err
}
