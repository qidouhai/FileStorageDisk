package db

import (
	"database/sql"
	"fmt"

	mydb "FileStorageDisk/db/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type TableFile struct {
	FileSha1 string
	FileName sql.NullString
	FileAddr sql.NullString
	FileSize sql.NullInt64
}

// OnFileUploadOK: 文件上传完成，保存meta到数据库
func OnFileUploadOK(fsha1, filename, fileaddr string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare(
		"INSERT ignore INTO tbl_file(`file_sha1`, `file_name`, `file_size`, " +
			"`file_addr`, `status`)VALUES(?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}

	defer stmt.Close()

	ret, err := stmt.Exec(fsha1, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	// 判断是否已经插入过相同fsha1的记录, 插入相同的会直接忽略
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 { // sql 执行成功了，但没有产生新的记录，抛个warning
			fmt.Printf("File with hash: %s has been upload before", fsha1)
		}
		return true
	}
	return false
}


func GetFileMeta(fhash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"SELECT file_sha1, file_name, file_addr, file_size "+
		"FROM tbl_file where file_sha1=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tf := TableFile{}
	// Scan可以把数据库取出的字段值赋值给指定的数据结构
	err = stmt.QueryRow(fhash).Scan(&tf.FileSha1, &tf.FileName, &tf.FileAddr, &tf.FileSize)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &tf, nil
}