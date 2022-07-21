package main

import (
	"FileStorageDisk/handler"
	"fmt"
	"net/http"
)

func main() {
	// static configure
	http.Handle("/static/", 
			http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	// filelayer
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))
	http.HandleFunc("/file/dowload", handler.HTTPInterceptor(handler.DownloadHandler))
	http.HandleFunc("/file/update", handler.HTTPInterceptor(handler.UpdateFileMetaHandler))
	http.HandleFunc("/file/delete", handler.HTTPInterceptor(handler.DeleteFileHandler))
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(
											handler.TryFastUploadHandler))

	// userlayer
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	// 
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server, err:%s\n", err.Error())
	}
}