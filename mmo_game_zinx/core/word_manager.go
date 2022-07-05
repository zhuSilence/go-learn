package core

import "sync"

type WorldManager struct {
	// AOIManager 当前世界地图 AOI 的管理模块
	AoiMgr *AOIManager
	//当前全部在线的 Player 集合
	Players map[int32]*Player
	// 读写锁
	pLock sync.RWMutex
}

// WorldMgrObj 提供一个全局的句柄
var WorldMgrObj *WorldManager

// 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNT_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNT_Y),
		Players: make(map[int32]*Player),
	}
}

// AddPlayer 添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {

	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

// RemovePlayerByPid 添加一个玩家
func (wm *WorldManager) RemovePlayerByPid(pID int32) {

	player := wm.Players[pID]
	wm.AoiMgr.RemoveFromGridByPos(int(pID), player.X, player.Z)

	wm.pLock.Lock()
	delete(wm.Players, pID)
	wm.pLock.Unlock()
}

// GetPlayerByPid 根据 pid 火气 Player 对象
func (wm *WorldManager) GetPlayerByPid(pID int32) *Player {

	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.Players[pID]
}

// GetAllPlayers 获取所有的 player
func (wm *WorldManager) GetAllPlayers() []*Player {

	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	players := make([]*Player, 0)

	for _, player := range wm.Players {
		players = append(players, player)
	}
	return players

}
