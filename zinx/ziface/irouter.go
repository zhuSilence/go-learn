package ziface

// IRouter 路由抽象接口，路由里面的数据都是 IRequest

type IRouter interface {
	// PreHandle 处理 conn 业务之前的钩子方法
	PreHandle(request IRequest)
	// Handle 处理 conn 业务方法
	Handle(request IRequest)
	// PostHandle 处理 conn 业务之后的钩子方法
	PostHandle(request IRequest)
}
