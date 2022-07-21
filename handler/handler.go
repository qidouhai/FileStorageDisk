package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	dblayer "FileStorageDisk/db"
	"FileStorageDisk/meta"
	"FileStorageDisk/util"
)

// UploadHandler: 文件上传接口
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传的HTML页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			fmt.Println("Internal server error")
			return
		}
		io.WriteString(w, string(data))

	} else if r.Method == "POST" {
		// 1.获取文件句柄、文件头、错误（如果有）
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data, err:%s\n", err.Error())
			return
		}
		defer file.Close()

		// 5. 设置文件元信息
		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: "/tmp/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 2.创建一个本地文件接收当前文件流
		// newFile, err := os.Create("/tmp/" + header.Filename)
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file, err:%s\n", err.Error())
			return
		}
		// 3. 将内存中的文件拷贝到newFile的buffer区
		// _, err = io.Copy(newFile, file)
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
			return
		}

		// 6. 更新FileMeta
		newFile.Seek(0, 0) // 把文件句柄的位置移到开始位置
		fileMeta.FileSha1 = util.FileSha1(newFile)
		// meta.UpdateFileMeta(fileMeta)
		// 持久化到数据库
		_ = meta.UpdateFileMetaToDB(fileMeta)
		r.ParseForm()
		username := r.Form.Get("username")
		ok := dblayer.OnUserFileUploadOK(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if ok {
			// 4. 向客户端返回成功信息/或重定向到一个成功页面
			// http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		}else {
			w.Write([]byte("Upload Failed."))
		}
	}
}

// TryFastUploadHandler：尝试秒传接口
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 1. 解析请求参数
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 2. 从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaFromDB(filehash)
	// fmt.Printf("fileMeta: %v\n", fileMeta)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. 查不到则返回秒传失败
	if fileMeta == nil {
		resp := util.RespMsg{
			Code: -1,
			Msg: "秒传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONBytes())
		return
	}
	// 4. 上传过则将文件信息写入用户文件表，返回成功
	ok := dblayer.OnUserFileUploadOK(username, filehash, filename, int64(filesize))
	if ok {
		resp := util.RespMsg{
			Code: 0,
			Msg: "秒传成功！",
		}
		w.Write(resp.JSONBytes())
		return
	}

	resp := util.RespMsg{
		Code: -2,
		Msg: "秒传失败，请稍后重传！",
	}
	w.Write(resp.JSONBytes())
	return
}

// QueryFileHandler: 查询批量文件信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	limitCnt, _ := strconv.Atoi(r.Form.Get("limit")) 
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// fmt.Printf("userFiles: %v\n", userFiles)

	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

// UploadSucHandler: 上传已完成页面
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler：通过文件sha1值获取文件元信息的接口
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	// 解析客户端发送请求的参数
	r.ParseForm()

	fh := r.Form["filehash"][0] // 默认第0个
	// fm := meta.GetFileMeta(fh)
	fm, err := meta.GetFileMetaFromDB(fh)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// 转为Json 字符串形式返回给客户端
	contentBytes, err := json.Marshal(fm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(contentBytes)
}

// DownloadHandler: 下载文件接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 拿到客户端发送来的sha1值
	fsha1 := r.Form.Get("filehash")
	// 获取元信息对象
	fm := meta.GetFileMeta(fsha1)
	// 从指定位置读入文件到内存，然后返回给客户端
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// 加载到内存(文件较小时可使用ioutil一次性全部加载到内存；
	// 文件较大时应要考虑实现流的形式)
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 加上http的响应头，让浏览器识别出来，然后就可以当成一个文件的下载
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// UpdateFileMetaHandler: 修改文件接口（重命名）
func UpdateFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 3个参数：待操作类型、fsha1值、新文件名
	opType := r.Form.Get("op") // 0 表示重命名操作
	fsha1 := r.Form.Get("filehash")
	newFilename := r.Form.Get("filename")

	// 暂时仅支持重名命操作
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// POST 请求
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 修改当前文件名
	curFileMeta := meta.GetFileMeta(fsha1)
	curFileMeta.FileName = newFilename
	meta.UpdateFileMeta(curFileMeta)

	// 转成json字符串形式，返回给客户端
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// DeleteFileHandler: 删除文件的接口
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	// 删除文件在"云端"的物理位置
	os.Remove(fm.Location)

	// 删除对应文件元信息的索引
	meta.RemoveFileMeta(fsha1)

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Delete successfully!")
}
