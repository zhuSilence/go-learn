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
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen),
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
