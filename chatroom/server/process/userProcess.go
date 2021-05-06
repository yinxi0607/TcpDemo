package process2

import (
	"TcpDemo/chatroom/common/message"
	"TcpDemo/chatroom/server/model"
	"TcpDemo/chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加一个字段，表示该Conn是那个用户的
	UserId int
}

//编写一个函数serverProcessLogin函数，专门处理登录的请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//核心代码
	//1、先从mes中取出mes.Data，并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json unmarshal fail err=", err)
		return
	}
	//1、先声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	//2、再声明一个loginResMes，并完成赋值
	var loginResMes message.LoginResMes
	//我们需要到redis数据库去完成验证
	//1、使用model.MyUserDao到redis去验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Error = err.Error()
			loginResMes.Code = 500
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		fmt.Println(user, "登录成功")
		loginResMes.Code = 200
		//这里因为用户登录成功，我们就把改登录成功的用户放入到userMgr中
		//将登录成功的用户的userId，赋给this
		this.UserId = loginMes.UserId
		userMgr.AddOnLineUser(this)
		//将当前用户的id 放入到loginResMes.UsersId
		//遍历userMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
	}

	////如果用户的id=100，密码=123456，认为合法，否则不合法
	//if loginMes.UserId==100 && loginMes.UserPwd=="123456"{
	//	//合法
	//	loginResMes.Code = 200
	//
	//}else {
	//	//不合法
	//	loginResMes.Code = 500 //500状态码表示该用户不存在
	//	loginResMes.Error = "该用户不存在，请注册再使用"
	//}

	//3、将loginResMes序列号
	marshal, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("序列化失败，err=", err)
		return
	}
	//4、将data赋值给resMes
	resMes.Data = string(marshal)

	//5、对resMes进行序列号，准备发送
	data, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("序列化失败，err=", err)
		return
	}

	//6、发送data，我们将其封装到writePkg
	//因为分层模式，我们先创建一个transfer实例
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	return
}

//编写一个函数serverProcessRegister函数，专门处理注册的请求
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.unmarshal fail err=", err)
		return err
	}
	//先声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//我们需要到redis数据库去完成注册
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误"
		}
	} else {
		registerResMes.Code = 200
	}
	marshal, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return err
	}
	resMes.Data = string(marshal)

	data, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return err
	}

	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return err
}
