package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	// 消息广播的 channel
	Message chan string
}

// NewServer 创建一个 Server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// ListenMessage 监听 Message 广播消息，
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("链接建立成功")
	user := NewUser(conn, this)
	//用户上线逻辑
	user.Online()
	//监听用户是否活跃 channel
	isLive := make(chan bool)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				//用户下线逻辑
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//提取用户的消息（去除\n）
			msg := string(buf[:n-1])
			//用户处理消息逻辑
			user.DoMessage(msg)

			//用户的任意消息，代表用户活跃
			isLive <- true
		}
	}()

	//当前 handler 阻塞
	for {
		select {
		case <-isLive:
			//当前用户活跃,重置定时器，这里可以不用添加逻辑，会进行下一个 case 的判断语句执行，重置定时器
		case <-time.After(time.Second * 300):
			//进入表示已经超时，进行强 T
			user.SendMsg("你被踢了，强制下线...\n")
			//关闭资源
			close(user.C)
			conn.Close()
			//退出 go 程
			return
		}

	}
}

func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	defer listener.Close()

	// 启动监听
	go this.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err: ", err)
			continue
		}

		//do
		go this.Handler(conn)
	}

	//close listen

}
