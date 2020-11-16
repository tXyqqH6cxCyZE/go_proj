package dao

import (
	"db/utils"
	"fmt"
)

/**
	数据库的增删改查
 */
type User struct {
	ID int
	Username string
	Password string
	Email string
}
// 通过预编译后stmt的AddUser()方法
func(user *User) AddUser() error {
	// 写sql语句
	sqlStr := "INSERT INTO users(username, password, email) VALUES(?,?,?)"

	// 预编译
	stmt, err := utils.Db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("预编译出现异常:", err)
		return err
	}

	// 执行
	_, erro := stmt.Exec(user.Username, user.Password, user.Email)
	if erro != nil {
		fmt.Println("执行出现异常：", erro)
		return erro
	}
	return nil
}

// 只通过DB的Exec方法AddUser
//func(user *User) AddUser() error {
//	sqlStr := "insert into users(username, password, email) values(?,?,?)"
//
//	_, erro := utils.Db.Exec(sqlStr, user.Username, user.Password, user.Email)
//	if erro != nil {
//		fmt.Println("执行出现异常：", erro)
//		return erro
//	}
//	return nil
//}

func(user *User) GetUserByID(userId int) (*User, error) {
	sqlStr := "SELECT * FROM users where id = ?"

	// 执行sql
	row := utils.Db.QueryRow(sqlStr, userId)

	// 声明3个变量
	var username string
	var password string
	var email string
	// 将各个字段中的值读到以上三个变量中
	err := row.Scan(&userId, &username, &password, &email)
	if err != nil {
		return nil, err
	}
	// 将3个变量的值赋给User结构体
	u := &User{
		ID: userId,
		Username: username,
		Password: password,
		Email: email,
	}

	return u, nil
}

func(user *User) GetUsers()([] *User, error) {
	sqlStr := "SELECT * FROM users"

	rows, err := utils.Db.Query(sqlStr)
	if err != nil {
		return nil, err
	}

	// 定义一个User切片
	var users []*User
	// 遍历
	defer rows.Close()	// defer 预先关闭rows
	for rows.Next(){
		// 声明四个变量
		var userID int
		var username string
		var password string
		var email string
		// 将各个字段中的值读到以上三个变量中
		err := rows.Scan(&userID, &username, &password, &email)
		if err != nil {
			return nil, err
		}
		// 将3个变量赋给User结构体
		u := &User{
			ID: userID,
			Username: username,
			Password: password,
			Email: email,
		}
		// 将u添加到users切片中
		users = append(users, u)
	}
	return users, nil
}