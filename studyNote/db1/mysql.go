package main

import (
	"database/sql"
	"fmt"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8")
	checkError(err)

	// insert the data
	stmt, err := db.Prepare("INSERT userinfo SET username=?, departname=?, created=?")
	checkError(err)

	res, err := stmt.Exec("astexie", "R & D department", "2012-12-09")
	checkError(err)

	id, err := res.LastInsertId()
	checkError(err)

	fmt.Println(id)

	// update the data
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkError(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkError(err)

	affect, err := res.RowsAffected()
	checkError(err)

	fmt.Println(affect)

	// search the data
	rows, err := db.Query("SELECT * FROM userinfo")
	checkError(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkError(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	// delete the data
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkError(err)

	res, err = stmt.Exec(id)
	checkError(err)

	affect, err = res.RowsAffected()
	checkError(err)

	fmt.Println(affect)

	_ = db.Close()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
