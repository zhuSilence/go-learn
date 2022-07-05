package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/zhuSilence/go-learn/mmo_game_zinx/core"
	"github.com/zhuSilence/go-learn/mmo_game_zinx/pb/pb"
	"github.com/zhuSilence/go-learn/zinx/ziface"
	"github.com/zhuSilence/go-learn/zinx/znet"
)

type WorldChatApi struct {
	znet.BaseRouter // 继承 base 路由
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	// 1. 解析客户端传递的 proto 协议
	protoMsg := &pb.Talk{}
	err := proto.Unmarshal(request.GetDate(), protoMsg)
	if err != nil {
		fmt.Println("talk unmarshal err:", err)
		return
	}
	// 2. 获取当前消息属于哪个玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("pid not exist")
		return
	}

	// 3. 根据 pid 获取 player 对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	// 4. 将消息广播到所有玩家
	player.Talk(protoMsg.Content)

}
