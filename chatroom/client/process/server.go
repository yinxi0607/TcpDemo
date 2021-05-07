package process

import (
	"TcpDemo/chatroom/common/message"
	"TcpDemo/chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func ShowMenu() {
	fmt.Println("-------恭喜xxx登录成功--------")
	fmt.Println("-------1、显示在线用户列表----------")
	fmt.Println("-------2、发送消息----------")
	fmt.Println("-------3、信息列表----------")
	fmt.Println("-------4、退出系统----------")
	fmt.Println("请选择（1-4）：")
	var key int
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		fmt.Println("显示在线用户列表")
		outputOnlineUser()
	case 2:
		fmt.Println("发送消息")
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("选择了退出系统。。。。bye")
		os.Exit(0)
	default:
		fmt.Println("您输入有误，请重新输入")
	}

}

//和服务器端保持通讯
func serverProcessMes(conn net.Conn) {
	//创建一个transfer实例，不停的读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("服务器端出错，err=", err)
			return
		}
		//如果读取到消息，又是下一步处理逻辑
		fmt.Printf("mes=%v\n", mes)
		switch mes.Type {
		case message.NotifyUserStatusMesType: //通知有人上线了
			//处理
			//1、取出NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2、把这个用户的信息、状态保存到map[int]User中
			updateUserStatus(&notifyUserStatusMes)
		default:
			fmt.Println("服务器端返回了一个未知的消息类型")
		}
	}
}
