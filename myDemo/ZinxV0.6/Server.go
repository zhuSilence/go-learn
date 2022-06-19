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

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	fmt.Println("recv msg from client msgId:", request.GetMsgID(), " msgData ", string(request.GetDate()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("sendMsg error", err)
	}

}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	fmt.Println("recv msg from client msgId:", request.GetMsgID(), " msgData ", string(request.GetDate()))

	err := request.GetConnection().SendMsg(201, []byte("welcome to zinxV0.6"))
	if err != nil {
		fmt.Println("sendMsg error", err)
	}

}

func main() {
	// 1. create a server with zinx
	s := znet.NewServer("[Zinx V0.6]")
	// 2. 添加 router

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	// 3. start server
	s.Server()
}
