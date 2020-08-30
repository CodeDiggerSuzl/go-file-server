package db

import (
	"database/sql"
	"fmt"
	"go-file-server/db/mysql"
)

// OnFileUploadFinished save to db after file uploaded
func OnFileUploadFinished(fileHash, fileName, fileAddr string, fileSize int64) bool {
	stmt, err := mysql.GetDbConn().Prepare("INSERT INTO tbl_file(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) VALUES(?,?,?,?,1)")
	if err != nil {
		fmt.Printf("insert error: %s", err)
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		fmt.Printf("Error Exec %s", err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil {
		// sql exec success, but nothing happend
		if rf <= 0 {
			fmt.Printf("file with filehash %s has ben upload before", fileHash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// GetFileMetaFromDb get file info form db
func GetFileMetaFromDb(fileSha1 string) (*TableFile, error) {
	stmt, err := mysql.GetDbConn().Prepare("SELECT file_sha1,file_addr,file_name,file_size FROM tbl_file WHERE file_sha1 = ? AND status = 1 LIMIT 1")
	if err != nil {
		fmt.Printf("err during prepare stmt: %s", err.Error())
		return nil, err
	}
	defer stmt.Close()
	tf := TableFile{}
	err = stmt.QueryRow(fileSha1).Scan(&tf.FileHash, &tf.FileAddr, &tf.FileAddr, &tf.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Query file info for nothing ")
			return nil, nil
		} else {
			fmt.Printf("error during queryRow: %s", err.Error())
			return nil, err
		}
	}
	return &tf, nil
}
