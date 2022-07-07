package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
		// 2.创建一个本地文件接收当前文件流
		newFile, err := os.Create("/tmp/" + header.Filename)
		if err != nil {
			fmt.Printf("Failed to create file, err:%s\n", err.Error())
			return
		}
		// 3. 将内存中的文件拷贝到newFile的buffer区
		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
			return
		}
		// 4. 向客户端返回成功信息/或重定向到一个成功页面
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSucHandler: 上传已完成页面
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}