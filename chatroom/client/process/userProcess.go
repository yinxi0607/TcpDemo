package process

import (
	"TcpDemo/chatroom/client/utils"
	"TcpDemo/chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	//暂时不需要字段。。。
}

//给关联一个用户登录的方法

//写一个函数，完成登录
func (this *UserProcess) Login(userId int, userPwd string) (err error) {
	//下一个就要开始订协议
	//fmt.Println("userId = %d ,userPwd = %s",userId,userPwd)
	//return nil
	//1、链接到服务器
	dial, err := net.Dial("tcp", "localhost:8889")
	tf := &utils.Transfer{
		Conn: dial,
	}
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//延时关闭
	defer dial.Close()

	//2、准备通过conn发送消息给服务器
	var mess message.Message
	mess.Type = message.LoginMesType

	//3、创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4、将loginMes序列化
	marshal, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json marshal err=", err)
		return
	}

	//5、把data赋给 mes.Data字段
	mess.Data = string(marshal)

	//6、将mess进行序列化
	marshal, err = json.Marshal(mess)
	if err != nil {
		fmt.Println("json marshal err=", err)
		return
	}

	//7、到这个时候，data就是我们要发送的消息
	//7.1 先把data的长度发送给服务器
	//先获取到data的长度--》转成一个表示长度的data切片
	var pkgLen uint32
	pkgLen = uint32(len(marshal))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//发送长度
	write, err := dial.Write(buf[:4])
	if write != 4 || err != nil {
		fmt.Println("send fail err=", err)
		return
	}
	//fmt.Printf("客户端发送消息的长度=%d,内容是%s\n",len(marshal),marshal)

	//7.2发送消息本身
	_, err = dial.Write(marshal)
	if err != nil {
		fmt.Println("send data fail err=", err)
		return
	}

	//这里还需要处理服务器端返回的消息

	mes, err := tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg err=", err)
		return
	}
	//将mes.Data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		//fmt.Println("登录成功")

		//可以显示当前在线用户数，遍历loginResMes.UsersId
		fmt.Println("当前在线用户列表如下：")
		for _, v := range loginResMes.UsersId {
			fmt.Println("用户id：\t", v)
		}
		fmt.Println("\n\n")

		//这里我们需要在客户端启动一个协程
		//该协程保持与服务器端的通讯，如果服务器端有数据推送给客户端
		//则接收并显示在客户端的终端
		go serverProcessMes(dial)

		//1、显示我们的登录成功的菜单【循环】
		for {
			ShowMenu()
		}

	} else {
		fmt.Println(loginResMes.Error)
	}
	return

}

//注册函数
func (this *UserProcess) Register(userId int, userPwd, userName string) (err error) {
	//1、链接到服务器
	dial, err := net.Dial("tcp", "localhost:8889")
	tf := &utils.Transfer{
		Conn: dial,
	}
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//延时关闭
	defer dial.Close()
	//2、准备通过conn发送消息给服务器
	var mess message.Message
	mess.Type = message.RegisterMesType

	//3、创建一个LoginMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName
	//4、将registerMes序列化
	marshal, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json marshal err=", err)
		return
	}

	//5、把data赋给 mes.Data字段
	mess.Data = string(marshal)

	//6、将mess进行序列化
	marshal, err = json.Marshal(mess)
	if err != nil {
		fmt.Println("json marshal err=", err)
		return
	}

	//发送data给服务器端
	err = tf.WritePkg(marshal)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
	}

	mes, err := tf.ReadPkg()
	if err != nil {
		fmt.Println("readPKg(conn) err=", err)
		return
	}
	//将mes的data部分进行反序列化RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功，请重新登录")
	} else {
		fmt.Println(registerResMes.Error)
	}
	return
}
