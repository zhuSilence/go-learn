package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/zhuSilence/go-learn/mmo_game_zinx/core"
	"github.com/zhuSilence/go-learn/mmo_game_zinx/pb/pb"
	"github.com/zhuSilence/go-learn/zinx/ziface"
	"github.com/zhuSilence/go-learn/zinx/znet"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	// 解析客户端传递的 proto 协议
	protoMsg := &pb.Position{}
	err := proto.Unmarshal(request.GetDate(), protoMsg)
	if err != nil {
		fmt.Println("proto unmarshal err:", err)
		return
	}
	// 获取当前移动的玩家位置
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty error pid not exist")
		return
	}
	fmt.Printf("player is moving pid = %d, move(%f, %f, %f, %f)", pid, protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32)) // 强转类型

	// 将移动的玩家位置进行更新并广播给其他玩家
	player.UpdatePos(protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)

}
