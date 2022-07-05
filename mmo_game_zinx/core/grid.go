package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	//格子 id
	GID  int
	MinX int
	MaxX int
	MinY int
	MaxY int
	// 格子内玩家成员集合
	playerIds map[int]bool
	pIdLock   sync.RWMutex
}

// NewGrid 初始化格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIds: make(map[int]bool),
	}
}

// Add 给格子添加一个玩家
func (g *Grid) Add(playerId int) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()
	g.playerIds[playerId] = true
}

// Remove 给格子删除一个玩家
func (g *Grid) Remove(playerId int) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()
	delete(g.playerIds, playerId)
}

// GetPlayerIds 获取当前格子中所有的玩家 ID
func (g *Grid) GetPlayerIds() (playerIds []int) {
	g.pIdLock.RLock()
	defer g.pIdLock.RUnlock()

	for k, _ := range g.playerIds {
		playerIds = append(playerIds, k)
	}
	return
}

// 调试使用：打印格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX:%d, minY: %d, maxY: %d, playerIds: %v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIds)
}
