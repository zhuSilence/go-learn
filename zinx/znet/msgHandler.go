package znet

import (
	"fmt"
	"github.com/zhuSilence/go-learn/zinx/utils"
	"github.com/zhuSilence/go-learn/zinx/ziface"
	"strconv"
)

// 消息处理模块的实现层

type MsgHandle struct {
	// 存放每个msgId 对应的 Router 处理方法
	Apis map[uint32]ziface.IRouter
	// 负责 Worker 任务的队列
	TaskQueue []chan ziface.IRequest
	// Worker 线程数
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	router, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		panic("msgId = " + strconv.Itoa(int(request.GetMsgID())) + " doesn't exist need register")
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断是否存在，不存在则添加
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeat api , msgId = " + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId, " success ")
}

func (mh *MsgHandle) StartWorkerPool() {
	// 根据 workerPoolSize 开启 worker 和分配空间
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动单个 worker
		go mh.startOneWorker(i)
	}
}

func (mh *MsgHandle) startOneWorker(wordId int) {
	fmt.Println("Worker ID = ", wordId, " started... ")

	for {
		select {
		// 如果有消息，获取的就是一个客户端的 request
		case request := <-mh.TaskQueue[wordId]:
			mh.DoMsgHandler(request)
		}
	}

}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 将消息平均分配给不同的 worker，根据客户端的链接 id 来分配

	workId := request.GetConnection().GetConnId() % mh.WorkerPoolSize
	fmt.Println("Add ConnId = ", request.GetConnection().GetConnId(),
		" request MsgId = ", request.GetMsgID(),
		" to WorkerId = ", workId)

	mh.TaskQueue[workId] <- request
}
