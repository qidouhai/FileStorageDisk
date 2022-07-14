package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql 数据库驱动
)

// 声明mysql连接对象
var db *sql.DB

// DB 初始化
func init() {
	db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1)/fileserver?charset=utf8")
	db.SetConnMaxIdleTime(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

// // DBConn: 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}