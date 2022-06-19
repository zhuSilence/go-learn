package ziface

type IServer interface {
	// Start server
	Start()

	// Stop server
	Stop()

	// Server run server
	Server()

	// AddRouter 路由功能，给当前的服务注册一个路由方法，供客户端的链接使用
	AddRouter(msgId uint32, router IRouter)
	// GetConnMgr 获取当前 server 的连接管理器
	GetConnMgr() IConnManager
}
