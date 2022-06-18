package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	//1. 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	//2. 建立连接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net Dial is error: ", err)
		return nil
	}
	//3. 返回对象
	client.conn = conn
	return client
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag = -1
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>>>>请输入合法范围内的数字...")
		return false
	}
}

func (client *Client) PublicChat() {
	var chatMsg string
	//提示用户输入消息
	fmt.Println(">>>>>>请输入聊天内容，exit 退出")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		//发送给服务器
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write error : ", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>>请输入聊天内容，exit 退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write error: ", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>>>>请输入需要聊天的用户名，exit 退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>>>请输入消息内容，exit 退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			//发送给服务器
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write error : ", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>>>请输入消息内容，exit 退出")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println(">>>>>>>请输入需要聊天的用户名，exit 退出")
		fmt.Scanln(&remoteName)

	}

}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>>>>请输入用户名：")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"

	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for !client.menu() {

		}
		//根据 flag 处理不同的业务
		switch client.flag {
		case 1:
			//公聊模式
			client.PublicChat()
		case 2:
			// 私聊模式
			client.PrivateChat()
		case 3:
			//更新用户名
			client.UpdateName()
		}
	}
}

var serverIp string
var serverPort int

//./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器端ip地址（默认是 127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端ip端口（默认是8888）")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>> 连接服务器失败...")
		return
	}

	//单独开启一个 go 程处理回执消息
	go client.DealResponse()

	fmt.Println(">>>>>>>>>>>> 连接服务器成功....")
	//启动客户端

	client.Run()
}
