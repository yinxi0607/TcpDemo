package main

import (
	"TcpDemo/chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//写一个函数，完成登录
func login(userId int, userPwd string) (err error) {
	//下一个就要开始订协议
	//fmt.Println("userId = %d ,userPwd = %s",userId,userPwd)
	//return nil
	//1、链接到服务器
	dial, err := net.Dial("tcp", "localhost:8889")
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
	mes, err := readPkg(dial)
	if err != nil {
		fmt.Println("readPkg err=", err)
		return
	}
	//将mes.Data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		fmt.Println("登录成功")
	} else {
		fmt.Println(loginResMes.Error)
	}
	return

}
