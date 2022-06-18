package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)
	//1. 直接链接远程服务器，的带一个 conn 链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err:", err)
		return
	}
	//2. 链接调用 Write 写数据
	for {
		_, err := conn.Write([]byte("hello zinx0.1"))
		if err != nil {
			fmt.Println("client Write err:", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("client Read err:", err)
			return
		}
		fmt.Printf("server call back, %s, cnt: %d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}

}
