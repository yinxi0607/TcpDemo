package model

import (
	"TcpDemo/chatroom/common/message"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

//我们在服务器启动后，就初始化一个userDao实例
//把他做成全局的变量，在需要和redis操作时，就直接使用即可
var (
	MyUserDao *UserDao
)

//定义一个UserDao结构体
//完成User结构体的各种操作

type UserDao struct {
	pool *redis.Pool
}

//使用工厂模式，获取到一个UserDao的实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

//思考一下，UserDao要完成什么方法

//1、根据用户id，返回一个User实例+err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	res, err := redis.String(conn.Do("HGET", "users", id))
	if err != nil {
		if err == redis.ErrNil {
			//表示users中没有找到对应的id
			err = ERROR_USER_NOTEXISTS
		}
		return nil, err
	}
	user = &User{}
	//这里我们需要把res反序列化成User实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json unmarshal err=", err)
		return nil, err
	}
	return
}

//完成登录的校验Login
//1、Login完成用户的验证
//2、如果用户的id和pwd都正确，则返回一个user实例
//3、如果用户的id或pwd有错误，则返回对应的错误信息
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	//先从UserDao的链接池中取出一根链接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return nil, err
	}
	//这步证明用户获取到了
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return nil, err
	}
	return user, nil
}

func (this *UserDao) Register(user *message.User) (err error) {
	//先从UserDao的连接池中取出一根链接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return err
	}
	//这时证明这个用户还没有注册过，则可以完成注册
	marshal, err := json.Marshal(user)
	if err != nil {
		return err
	}
	//入库
	_, err = conn.Do("HSET", "users", user.UserId, marshal)
	if err != nil {
		fmt.Println("保存注册用户错误 err=", err)
		return err
	}
	return
}
