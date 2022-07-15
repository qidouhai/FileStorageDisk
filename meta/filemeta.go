package meta

import (
	mydb "FileStorageDisk/db"
)

type FileMeta struct {
	FileSha1 string  // 文件唯一标识 
	FileName string
	FileSize int64
	Location string  // 文件在云端的路径
	UploadAt string // 文件上传时间
}

// 保存每个上传文件的元信息
var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta: 更新或新增文件元信息
func UpdateFileMeta(fm FileMeta) {
	fileMetas[fm.FileSha1] = fm
}

// UpdateFileMetaDB: 新增/更新文件元信息到数据库
func UpdateFileMetaToDB(fm FileMeta) bool {
	return mydb.OnFileUploadOK(fm.FileSha1, fm.FileName, fm.Location, fm.FileSize)
}

// GetFileMeta: 通过文件的sha1值获取对应的文件元信息
func GetFileMeta(fhash string) FileMeta {
	return fileMetas[fhash]
}

// GetFileMetaDB：从mysql获取元信息
func GetFileMetaFromDB(fhash string) (FileMeta, error) {
	tf, err := mydb.GetFileMeta(fhash)
	if err != nil {
		return FileMeta{}, err
	}

	fm := FileMeta{
		FileSha1: tf.FileSha1,
		FileName: tf.FileName.String,
		Location: tf.FileAddr.String,
		FileSize: tf.FileSize.Int64,
	}
	return fm, nil
}

// 
func RemoveFileMeta(fsha1 string) {
	delete(fileMetas, fsha1)
}