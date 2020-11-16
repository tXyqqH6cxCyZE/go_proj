package dao

import (
	"fmt"
	"testing"
)

/**
	单元测试：
	如果一个测试函数的函数名不是以Test开头，那么在使用go test命令时
	默认不会执行，不过我们可以设置该函数时一个子测试函数，可以在其他测试
	函数里通过t.Run方法来执行子测试函数
 */

func TestMain(m *testing.M){
	fmt.Println("---------------------测试开始---------------------")
	m.Run()
	fmt.Println("---------------------测试结束---------------------")
}


func TestUser(t *testing.T){
	t.Run("正在测试添加用户：", testAddUser)
	t.Run("正在测试获取一个用户", testGetUserById)
}
func testAddUser(t *testing.T) {
	fmt.Println("测试添加用户：")
	user := &User{
		Username: "dnw",
		Password: "123456",
		Email: "admin3@pjx.com",
	}

	// 将user添加到数据库
	user.AddUser()
}

func testGetUserById(t *testing.T){
	u := &User{}
	user, _ := u.GetUserByID(1)
	fmt.Println("用户的信息是：", *user)
}
