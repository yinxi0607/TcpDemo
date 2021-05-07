package model

import (
	"TcpDemo/chatroom/common/message"
	"net"
)

//因为在客户端，很多地方会使用到curUser，我们将其做一个全局
type CurUser struct {
	Conn net.Conn
	message.User
}
