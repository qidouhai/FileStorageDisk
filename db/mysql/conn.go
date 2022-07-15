package mysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"  // 导入mysql驱动
)

// 声明mysql连接对象
var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3307)/fileserver?charset=utf8")
	db.SetConnMaxIdleTime(1000)
	err := db.Ping()  // test connection
	if err != nil {
		log.Fatal("Failed to connect to mysql, err: " + err.Error())
	}
}

// DBConn返回数据库连接对象
func DBConn() *sql.DB {
	return db
}