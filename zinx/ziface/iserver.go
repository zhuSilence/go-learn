package ziface

type IServer interface {
	// Start server
	Start()

	// Stop server
	Stop()

	// Server run server
	Server()

	// AddRouter 路由功能，给当前的服务注册一个路由方法，供客户端的链接使用
	AddRouter(router IRouter)
}
