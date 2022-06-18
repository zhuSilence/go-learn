package main

import (
	"fmt"
	"github.com/zhuSilence/go-learn/zinx/ziface"
	"github.com/zhuSilence/go-learn/zinx/znet"
)

/**
基于 Zinx 框架开发的服务器用于程序
*/

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error ", err)
	}
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("call back ping error ", err)
	}
}

func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error ", err)
	}
}

func main() {
	// 1. create a server with zinx
	s := znet.NewServer("[Zinx V0.4]")
	// 2. 添加 router
	s.AddRouter(&PingRouter{})
	// 3. start server
	s.Server()
}
