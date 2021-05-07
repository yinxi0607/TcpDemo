package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

//这里我们定义几个用户状态的常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息的类型
}

//定义两个消息。后面需要再增加
type LoginMes struct {
	UserId   int    `json:"userId"`   //用户id
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户名
}

type LoginResMes struct {
	Code    int    `json:"code"` // 返回的状态码 500表示用户未注册 200 表示登录成功
	UsersId []int  //增加保存用户id的切片
	Error   string `json:"error"` //返回的错误信息
}

type RegisterMes struct {
	User User //类型就是User结构体
}

type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码400表示该用户已经存在，200表示注册成功
	Error string `json:"error"` //返回错误信息
}

//为了配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"` //用户id
	Status int `json:"status"` //用户的状态
}

//增加一个smsMes  发送的消息
type SmsMes struct {
	Content string `json:"content"` //内容
	User           //匿名结构体
}
