package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

//创建一个 user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动监听
	go user.ListenMessage()

	return user
}

func (this *User) ListenMessage() {

	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	//用户上线，加入用户
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	// 广播用户上线
	this.server.BroadCast(this, "已上线")
}

func (this *User) Offline() {
	//用户下线，删除用户
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	// 广播用户下线
	this.server.BroadCast(this, "下线")
}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		//查询在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + "在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//rename 用户名
		newName := strings.Split(msg, "|")[1]
		//判断 name 是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名已经存在\n")
		} else {
			this.server.mapLock.Lock()
			this.server.OnlineMap[newName] = this
			delete(this.server.OnlineMap, this.Name)
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("您已经更新用户名为：" + newName + "\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式 to|张三|消息内容
		//1. 获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("消息格式不正确，请使用\"to|张三|你好格式\"")
			return
		}
		//2. 根据用户名，找到对方的 User 对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("该用户不存在\n")
			return
		}

		//3. 获取消息内容，通过对方的User 对象将消息发送出去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("消息不能为空\n")
			return
		}
		remoteUser.SendMsg(this.Name + "对您说：" + content + "\n")
	} else {
		this.server.BroadCast(this, msg)
	}
}
