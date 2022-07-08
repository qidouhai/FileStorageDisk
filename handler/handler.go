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

// 接口： 通过文件sha1值获取文件元信息
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