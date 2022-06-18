package znet

import (
	"fmt"
	"github.com/zhuSilence/go-learn/zinx/ziface"
	"net"
)

type Connection struct {
	// 当前连接的 socket TCP 套接字
	Conn *net.TCPConn
	// 连接的 ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool
	// 当前连接所绑定的处理业务的方法
	handleApi ziface.HandleFunc
	// 告知岗前连接已经退出的 channel
	ExitChan chan bool
	// 当前的链接处理方法 router
	Router ziface.IRouter
}

// NewConnection 初始化连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandleFunc) *Connection {

	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		handleApi: callback_api,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connId=", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		// 得到当前链接的 Request 数据
		req := Request{
			conn: c,
			data: buf,
		}
		// 从路由中找到注册绑定的 Conn 对应的 router，进行执行
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

		//// 调用当前连接所绑定的 handleAPI
		//if err := c.handleApi(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("connId=", c.ConnID, "handleApi err", err)
		//	break
		//}
	}

}
func (c *Connection) Start() {
	fmt.Println("Connection start , connId ", c.ConnID)
	// 启动从当前的连接读数据
	go c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Println("Connection stop , connId ", c.ConnID)
	if c.isClosed {
		return
	}
	// 回收资源
	c.isClosed = true
	c.Conn.Close()
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}
