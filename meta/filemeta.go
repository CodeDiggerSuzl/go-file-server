package meta

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

func DelFileMeta(fileSha1 string) {
	delete(fileMetaMap, fileSha1)
}
