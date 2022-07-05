package main

import (
	"fmt"
	"github.com/zhuSilence/go-learn/mmo_game_zinx/apis"
	"github.com/zhuSilence/go-learn/mmo_game_zinx/core"
	"github.com/zhuSilence/go-learn/zinx/ziface"
	"github.com/zhuSilence/go-learn/zinx/znet"
)

func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个 Player
	player := core.NewPlayer(conn)
	// 给客户端发送 msgId：1 的消息
	player.SyncPid()
	// 给客户端发送 msgId：200 的消息
	player.BroadCastStartPosition()
	// 将当前上线的玩家添加到 word 中
	core.WorldMgrObj.AddPlayer(player)
	// 当前连接绑定 pid
	conn.SetProperty("pid", player.Pid)
	// 通知其他玩家当前玩家上线，广播当前玩家位置信息
	player.SyncSurrounding()

	fmt.Println("player pid=", player.Pid, " is online")
}

func OnConnectionLost(conn ziface.IConnection) {
	// 得到当前玩家
	pid, _ := conn.GetProperty("pid")
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	// 玩家下线
	player.Offline()

	fmt.Println("=====> Player pid = ", pid, " offline...")
}

func main() {

	//创建 Zinx 服务句柄
	s := znet.NewServer("MMO Game")
	// 连接创建和销毁钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)
	// 注册接口
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})
	// 启动服务
	s.Server()
}
