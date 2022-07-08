package meta

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

// 更新或新增文件元信息
func UpdateFileMeta(fm FileMeta) {
	fileMetas[fm.FileSha1] = fm
}

// 通过文件的sha1值获取对应的文件元信息
func GetFileMeta(fhash string) FileMeta {
	return fileMetas[fhash]
}