package main

import (
	"TcpDemo/chatroom/server/model"
	"fmt"
	"net"
	"time"
)

//
//func readPkg(conn net.Conn)(mess message.Message,err error){
//	buf := make([]byte,1024 * 4)
//	fmt.Println("读取客户端发送的数据")
//	//conn.Read 在conn没有被关闭的情况下，才会阻塞
//	//如果客户端关闭了，conn就不会阻塞
//	read, err := conn.Read(buf[:4])
//
//	if err!=nil || read!=4{
//		fmt.Println("conn.read err=",err)
//		//err = errors.New("read pkg header error")
//		return
//	}
//	//根据buf[:4]转化成一个uint32类型
//	var pkgLen uint32
//	pkgLen = binary.BigEndian.Uint32(buf[:4])
//
//	//根据pkglen读取消息内容
//	n, err := conn.Read(buf[:pkgLen])
//	if n!=int(pkgLen) || err!=nil{
//		fmt.Println("conn.Read fail err=",err)
//		err = errors.New("read pkg body error")
//		return
//	}
//
//	//把pkgLen 反序列化 --》message.Message{}
//	err = json.Unmarshal(buf[:pkgLen],&mess)
//	if err!=nil{
//		fmt.Println("json.Unmarshal fail err=",err)
//		err = errors.New("json unmarshal error")
//		return
//	}
//	return
//}
//
//func writePkg(conn net.Conn,data []byte)(err error){
//	//先发送一个长度给对方
//	var pkgLen uint32
//	pkgLen = uint32(len(data))
//	var buf [4]byte
//	binary.BigEndian.PutUint32(buf[0:4],pkgLen)
//	//发送长度
//	write, err := conn.Write(buf[:4])
//	if write!=4||err!=nil{
//		fmt.Println("send fail err=",err)
//		return
//	}
//	//发送data本身
//	write, err = conn.Write(data)
//	if write!=int(pkgLen)||err!=nil{
//		fmt.Println("send fail err=",err)
//		return
//	}
//	return
//}

//
////编写一个函数serverProcessLogin函数，专门处理登录的请求
//func serverProcessLogin(conn net.Conn,mes *message.Message)(err error){
//	//核心代码
//	//1、先从mes中取出mes.Data，并直接反序列化成LoginMes
//	var loginMes message.LoginMes
//	err = json.Unmarshal([]byte(mes.Data), &loginMes)
//	if err!=nil{
//		fmt.Println("json unmarshal fail err=",err)
//		return
//	}
//	//1、先声明一个resMes
//	var resMes message.Message
//	resMes.Type = message.LoginResMesType
//
//	//2、再声明一个loginResMes，并完成赋值
//	var loginResMes message.LoginResMes
//
//	//如果用户的id=100，密码=123456，认为合法，否则不合法
//	if loginMes.UserId==100 && loginMes.UserPwd=="123456"{
//		//合法
//		loginResMes.Code = 200
//
//	}else {
//		//不合法
//		loginResMes.Code = 500 //500状态码表示该用户不存在
//		loginResMes.Error = "该用户不存在，请注册再使用"
//	}
//
//	//3、将loginResMes序列号
//	marshal, err := json.Marshal(loginResMes)
//	if err != nil{
//		fmt.Println("序列号失败，err=",err)
//		return
//	}
//	//4、将data赋值给resMes
//	resMes.Data = string(marshal)
//
//	//5、对resMes进行序列号，准备发送
//	data ,err := json.Marshal(resMes)
//	if err!=nil{
//		fmt.Println("序列号失败，err=",err)
//		return
//	}
//
//	//6、发送data，我们将其封装到writePkg
//	err = writePkg(conn,data)
//	return
//}

////编写一个serverProcessMes函数
////功能：根据客户端发送消息种类不同，决定调用哪个函数来处理
//func serverProcessMes(conn net.Conn,mes *message.Message)(err error){
//	switch mes.Type {
//	case message.LoginMesType:
//		//处理登录逻辑
//		err = serverProcessLogin(conn,mes)
//	case message.LoginResMesType:
//		//处理登录返回的逻辑
//		default:
//		fmt.Println("消息类型不存在，无法处理。。。")
//	}
//	return
//}

func process(conn net.Conn) {
	//这里需要延时关闭conn
	defer conn.Close()

	//调用总控，创建一个总控
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端与服务器端，通讯携程错误，err=", err)
		return
	}
}

//这里我们编写一个函数，完成对UserDao的初始化任务
func initUserDao() {
	//这里的pool本身就是一个全局的变量
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	//当服务器启动时，我就就去初始化我们的redis连接池
	initPool("10.26.14.98:52385", 16, 0, 300*time.Second)
	initUserDao()
	fmt.Println("服务器在8889端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net listen err=", err)
		return
	}
	//一旦监听成功，就等待客户端连接服务器
	for {
		fmt.Println("等待客户端来连接服务器")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen accept err = ", err)
		}
		//一旦连接成功，则启动一个协程和客户端保持通讯
		go process(conn)
	}

}
