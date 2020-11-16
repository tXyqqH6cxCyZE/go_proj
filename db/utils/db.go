package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/**
	数据库工具包：
	init() : 建立数据库连接
 */
var(
	Db *sql.DB
	err error
)

func init(){
	// Open函数可能只是验证其参数，而不创建与数据库的连接。
	// 如果要检查数据源的名称是否合法，应调用返回值的Ping方法
	Db, err = sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}


}


