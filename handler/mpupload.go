package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	rPool "FileStorageDisk/cache/redis"
	"FileStorageDisk/util"
)

// 分块上传的初始化信息
type MultipartUploadInfo struct {
	FileHash string
	FileSize int
	UploadId string
	ChunkSize int
	ChunkCount int
}

// InitialMultipartUploadHandler: 初始化分块上传
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehas")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}

	// 2. 获取redis的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash: filehash,
		FileSize: filesize,
		UploadId: username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize: 5 * 1024 * 1024,  // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}
	
	// 4. 将初始化的信息写入Redis缓存
	rConn.Do("HSET", "MP_" + upInfo.UploadId, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_" + upInfo.UploadId, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_" + upInfo.UploadId, "filesize", upInfo.FileSize)
	
	// 5. 将响应初始化数据返回到客户端
	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())
}