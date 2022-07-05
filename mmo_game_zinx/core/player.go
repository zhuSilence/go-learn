package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/zhuSilence/go-learn/mmo_game_zinx/pb/pb"
	"github.com/zhuSilence/go-learn/zinx/ziface"
	"math/rand"
	"sync"
)

type Player struct {
	Pid  int32
	Conn ziface.IConnection // 当前玩家的连接
	X    float32            // 平面的 x 坐标
	Y    float32            // 高度
	Z    float32            // 平面的 y 坐标
	V    float32            // 旋转的角度 0 - 360 度
}

var PidGen int32 = 32 //用来生成的玩家 id 计数器
var IdLock sync.Mutex //保护锁

// NewPlayer 创建玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {

	// 生成一个玩家 id
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()
	// 创建玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), // 随机偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), // 随机偏移
		V:    0,                            // 角度默认为 0
	}
	return p
}

// SendMsg 提供一个发送消息的方法，将 protobuf 数据序列化
func (p *Player) SendMsg(msgId uint32, data proto.Message) {

	//将 Proto Message 结构体序列化，转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return
	}

	// 将二进制文件，通过 zinx 发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil.")
		return
	}

	// 发送消息
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("player send msg err", err)
		return
	}
	return
}

// SyncPid 同步 PlayerId 给客户端
func (p *Player) SyncPid() {
	// 组装 msgId ：1 的 proto 数据
	data := &pb.SyncPID{Pid: p.Pid}
	// 发送数据
	p.SendMsg(1, data)
}

// BroadCastStartPosition 同步 Payer 位置给客户端
func (p *Player) BroadCastStartPosition() {
	// 组装 msgId ：200 的 proto 数据
	data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{P: &pb.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
			V: p.V,
		}},
	}
	// 发送数据
	p.SendMsg(200, data)
}

// Talk 玩家发送talk
func (p *Player) Talk(content string) {
	// 组装 msgId ：200 的 proto 数据
	data := &pb.BroadCast{
		Pid:  p.Pid,
		Tp:   1, // 代表聊天广播
		Data: &pb.BroadCast_Content{Content: content},
	}
	// 得到所有的玩家
	players := WorldMgrObj.GetAllPlayers()
	// 向所有的玩家发送数据
	for _, player := range players {
		player.SendMsg(200, data)
	}
}

// SyncSurrounding 同步玩家位置
func (p *Player) SyncSurrounding() {
	// 1. 获取周边玩家
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	fmt.Println("=========================size is ", len(pids))
	if len(pids) <= 0 {
		return
	}
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		player := WorldMgrObj.GetPlayerByPid(int32(pid))
		players = append(players, player)
	}

	// 2. 将当前玩家的位置信息发送给其他玩家
	data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{P: &pb.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
			V: p.V,
		}},
	}
	// 发送数据
	for _, player := range players {
		player.SendMsg(200, data)
	}

	// 3. 向所有玩家发送位置 202 消息
	positionList := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		position := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		positionList = append(positionList, position)
	}

	syncProtoMsg := &pb.SyncPlayers{
		Ps: positionList[:],
	}
	p.SendMsg(202, syncProtoMsg)
}

// UpdatePos 更新当前玩家的位置，并广播给相邻玩家
func (p *Player) UpdatePos(x, y, z, v float32) {
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	// 组装广播消息
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4, // 4 移动后的位置信息
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 获取周边相邻玩家
	players := p.GetSurroundPlayers()

	for _, p := range players {
		p.SendMsg(200, protoMsg)
	}
}

// GetSurroundPlayers 获取周边相邻玩家
func (p *Player) GetSurroundPlayers() []*Player {
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	// 将所有的 pid 对应的 player 存放到切片中
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		player := WorldMgrObj.GetPlayerByPid(int32(pid))
		players = append(players, player)
	}
	return players
}

func (p *Player) Offline() {

	protoMsg := &pb.SyncPID{
		Pid: p.Pid,
	}

	players := p.GetSurroundPlayers()
	fmt.Println("===============current players size is ", len(players))

	for _, player := range players {
		fmt.Println("current player is ", player.Pid)
		if player.Pid != p.Pid {
			player.SendMsg(201, protoMsg)
		}
	}

	// 移除玩家信息
	WorldMgrObj.RemovePlayerByPid(p.Pid)
}
