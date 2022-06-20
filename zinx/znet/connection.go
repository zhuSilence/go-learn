package znet

import (
	"errors"
	"fmt"
	"github.com/zhuSilence/go-learn/zinx/utils"
	"github.com/zhuSilence/go-learn/zinx/ziface"
	"io"
	"net"
	"sync"
)

type Connection struct {
	// 当前 conn 属于的 server
	TcpServer ziface.IServer
	// 当前连接的 socket TCP 套接字
	Conn *net.TCPConn
	// 连接的 ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool
	// 告知岗前连接已经退出的 channel
	ExitChan chan bool
	// 增加一个读写分离的 channel
	msgChan chan []byte
	// 当前的链接处理方法 msgHandler
	MsgHandler ziface.IMsgHandler
	// 连接属性集合
	Property map[string]interface{}
	// 读写锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化连接的方法
func NewConnection(TcpServer ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {

	c := &Connection{
		TcpServer:  TcpServer,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		Property:   make(map[string]interface{}),
	}

	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connId=", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 读取客户端的数据
		// 创建拆包解包对象
		dp := NewDataPack()
		// 读取 msg head 二进制流 8 个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}
		// 二进制流拆包，到的 msgId 和 msgDataLen
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		// 根据 msg dataLen 再次读取 data
		if msg.GetDataLen() > 0 {
			data := make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
			msg.SetData(data)
		}

		// 得到当前链接的 Request 数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 开启的工作池，则使用工作池处理消息
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中找到注册绑定的 Conn 对应的 router，进行执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}

}

// StartWrite 监听管道进行消息的回写
func (c *Connection) StartWrite() {
	fmt.Println("Write Goroutine is running...")
	defer fmt.Println("connId=", c.ConnID, " Writer is exit, remote addr is ", c.RemoteAddr().String())

	// 不断的阻塞等待 channel 的消息
	for {
		select {
		case data := <-c.msgChan:
			// 有数据写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error", err)
				return
			}
		case <-c.ExitChan:
			// 代表 Reader 已经退出，此时 Write 需要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Connection start , connId ", c.ConnID)
	// 启动从当前的连接读数据
	go c.StartReader()
	// 启动从当前的连接写数据
	go c.StartWrite()

	// 创建连接之后执行钩子函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Connection stop , connId ", c.ConnID)
	if c.isClosed {
		return
	}
	// 回收资源
	c.isClosed = true
	// 告知 write 关闭
	c.ExitChan <- true

	// 销毁连接之前执行钩子函数
	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack msg error", err)
		return errors.New("pack msg error")
	}

	// 将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.Property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.Property[key]; ok {
		return value, nil
	}
	return nil, errors.New("not found")
}

func (c *Connection) Remove(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.Property, key)
}
