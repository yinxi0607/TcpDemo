package main

import (
	"TcpDemo/chatroom/client/process"
	"fmt"
)

//定义两个变量，一个表示用户id，一个表示用户的密码
var userId int
var userPwd string
var userName string

func main() {
	//接收用户的选址
	var key int
	//判断是否还继续显示菜单
	var loop = true

	for loop {
		fmt.Println("------------------欢迎登陆多人聊天系统-----------------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择（1-3）")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			fmt.Println("请输入用户的id")
			fmt.Scanf("%d\n", &userId)

			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n", &userPwd)
			//完成登录
			//1、创建一个UserProcess的实例
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
			//loop = false
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id：")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户昵称：")
			fmt.Scanf("%s\n", &userName)
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)
			//2、调用
			//loop = false
		case 3:
			fmt.Println("退出系统")
			//loop = false
		default:
			fmt.Println("您的输入有误，请重新输入")

		}
	}
	////根据用户的输入，显示新的提示信息
	//if key==1{
	//
	//
	//	//因为使用了新的程序结构，我们创建
	//	//先把登录的函数，写到另外一个文件，比如login.go
	//	//这里我们需要重新调用
	//
	//	//login(userId, userPwd)
	//	//if err != nil{
	//	//	fmt.Println("登录失败")
	//	//}else {
	//	//	fmt.Println("登录成功")
	//	//}
	//}else if key==2 {
	//	fmt.Println("进行用户注册")
	//}
}
