package ziface

// IMsgHandler 消息管理抽象层
type IMsgHandler interface {

	// DoMsgHandler 调度执行对的 Router 消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
	// StartWorkerPool 启动一个 zinx 线程池，用于处理客户端的任务
	StartWorkerPool()
	// SendMsgToTaskQueue 将消息发送给消息队列
	SendMsgToTaskQueue(request IRequest)
}
