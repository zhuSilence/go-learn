package ziface

import "net"

type IConnection interface {
	// Start 启动连接，让当前的连接准备开始工作
	Start()
	// Stop 停止连接，结束当前连接的工作
	Stop()
	// GetTCPConnection 获取当前连接的绑定 socket conn
	GetTCPConnection() *net.TCPConn
	// GetConnId 获取当前连接的连接 ID
	GetConnId() uint32
	// RemoteAddr 获取远程客户端的 TCP 状态 IP Port
	RemoteAddr() net.Addr
	// SendMsg 发送数据，将数据发送给远程的客户端
	SendMsg(msgId uint32, data []byte) error
	// SetProperty 设置连接属性
	SetProperty(key string, value interface{})
	// GetProperty 获取连接属性
	GetProperty(key string) (interface{}, error)
	// Remove 删除连接属性
	Remove(key string)
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
