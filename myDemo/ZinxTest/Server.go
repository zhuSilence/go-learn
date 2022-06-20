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
	s := znet.NewServer("[Zinx V0.8]")

	// 2. 注册钩子函数
	s.SetOnConnStart(func(connection ziface.IConnection) {
		fmt.Println("====> SetOnConnStart is called....")
		if err := connection.SendMsg(202, []byte("====> SetOnConnStart is called....")); err != nil {
			fmt.Println("====> SetOnConnStart is called err....", err)
		}

		connection.SetProperty("name", "ziyou")
		connection.SetProperty("name1", "ziyou1")
		connection.SetProperty("name2", "ziyou2")
		connection.SetProperty("name3", "ziyou3")

	})

	s.SetOnConnStop(func(connection ziface.IConnection) {
		fmt.Println("====> SetOnConnStop is called....")
		fmt.Println(connection.GetProperty("name"))
		fmt.Println(connection.GetProperty("name1"))
		fmt.Println(connection.GetProperty("name2"))
		fmt.Println(connection.GetProperty("name3"))
	})

	// 3. 添加 router

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	// 4. start server
	s.Server()
}
