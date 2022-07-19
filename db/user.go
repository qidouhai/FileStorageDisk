package db

import (
	mydb "FileStorageDisk/db/mysql"
	"fmt"
)

// UserSignup: 通过用户名和密码完成user表的注册
func UserSignup(username, password string) bool {
	// 插入一条用户数据
	stmt, err := mydb.DBConn().Prepare("INSERT ignore INTO tbl_user(user_name, user_pwd)" + 
						  "VALUES(?,?)")
	if err != nil {
		fmt.Println("Failed to insert, err: " + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to insert, err: " + err.Error())
		return false
	}

	// 校验是否注册成功（重复注册也当做注册失败）
	if rows, err := ret.RowsAffected(); nil == err && rows > 0 {
		return true
	}
	return false
}


// UserSignin: 判断密码是否一致
func UserSignin(username, encPasswd string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found: " + username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	pwd := string(pRows[0]["user_pwd"].([]byte))
	if len(pRows) > 0 && pwd == encPasswd {
		return true
	}
	return false
}

// UpdateToken : 刷新用户登录的token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`,`user_token`) values (?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username string
	Email string
	Phone string
	SignupAt string
}

func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare("SELECT user_name, signup_at FROM tbl_user WHERE user_name=? LIMIT 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	
	return user, nil
}