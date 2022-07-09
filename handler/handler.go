package handler

import (
	"FileStorageDisk/meta"
	"FileStorageDisk/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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
		meta.UpdateFileMeta(fileMeta)

		// 4. 向客户端返回成功信息/或重定向到一个成功页面
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSucHandler: 上传已完成页面
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler：通过文件sha1值获取文件元信息的接口
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	// 解析客户端发送请求的参数
	r.ParseForm()

	fh := r.Form["filehash"][0]  // 默认第0个
	fm := meta.GetFileMeta(fh)

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
	w.Header().Set("content-disposition", "attachment;filename=\"" + fm.FileName + "\"")
	w.Write(data)
}

// UpdateFileMetaHandler: 修改文件接口（重命名）
func UpdateFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 3个参数：待操作类型、fsha1值、新文件名
	opType := r.Form.Get("op")  // 0 表示重命名操作
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

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	// 删除对应文件元信息的索引
	meta.RemoveFileMeta(fsha1)

	// 删除文件在"云端"的物理位置
	os.Remove(fm.Location)

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Delete successfully!")
}