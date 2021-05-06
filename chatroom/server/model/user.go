package model

//定义一个用户的结构体
type User struct {
	//确定字段信息
	//为了序列化和反序列化成功，我们必须保证
	//用户信息的字符串的key与结构体的字段对
	//应的tag保持一致，否则会失败
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}
