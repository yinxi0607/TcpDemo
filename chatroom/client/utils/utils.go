package utils

import (
	"TcpDemo/chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

//这里将这些方法关联到结构体中

type Transfer struct {
	//分析它应该有哪些字段
	Conn net.Conn
	Buf  [8096]byte //这是传输时使用的缓存

}

func (this *Transfer) ReadPkg() (mess message.Message, err error) {
	//buf := make([]byte,1024 * 4)
	fmt.Println("读取客户端发送的数据")
	//conn.Read 在conn没有被关闭的情况下，才会阻塞
	//如果客户端关闭了，conn就不会阻塞
	read, err := this.Conn.Read(this.Buf[:4])

	if err != nil || read != 4 {
		fmt.Println("conn.read err=", err)
		//err = errors.New("read pkg header error")
		return
	}
	//根据buf[:4]转化成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])

	//根据pkglen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Read fail err=", err)
		err = errors.New("read pkg body error")
		return
	}

	//把pkgLen 反序列化 --》message.Message{}
	err = json.Unmarshal(this.Buf[:pkgLen], &mess)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		err = errors.New("json unmarshal error")
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	//先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	//发送长度
	write, err := this.Conn.Write(this.Buf[:4])
	if write != 4 || err != nil {
		fmt.Println("send fail err=", err)
		return
	}
	//发送data本身
	write, err = this.Conn.Write(data)
	if write != int(pkgLen) || err != nil {
		fmt.Println("send fail err=", err)
		return
	}
	return
}
