package message

const (
	LoginMesType       = "LoginMes"
	LoginResMesType    = "LoginResMes"
	RegisterMesType    = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
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
