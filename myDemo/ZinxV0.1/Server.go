package main

import "github.com/zhuSilence/go-learn/zinx/znet"

/**
基于 Zinx 框架开发的服务器用于程序
*/
func main() {
	// 1. create a server with zinx
	s := znet.NewServer("[Zinx V0.1]")
	// 2. start server
	s.Server()
}
