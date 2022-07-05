package core

import "fmt"

const (
	AOI_MIN_X = 85
	AOI_MAX_X = 410
	AOI_CNT_X = 10
	AOI_MIN_Y = 75
	AOI_MAX_Y = 510
	AOI_CNT_Y = 20
)

type AOIManager struct {
	// 区域左边界
	MinX int
	// 区域右边界
	MaxX int
	// X 方向的各自数量
	CntsX int
	// 区域上边界
	MinY int
	// 区域下边界
	MaxY int
	// Y 方向各自数量
	CntsY int
	// 当前区域中有哪些格子
	grids map[int]*Grid
}

// NewAOIManager 初始化方法
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiManager := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}
	// 给 AOI 初始化区域的格子
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//计算格子 id
			gid := y*cntsX + x
			// 初始化 gid 格子
			aoiManager.grids[gid] = NewGrid(gid,
				aoiManager.MinX+x*aoiManager.gridWidth(),
				aoiManager.MinX+(x+1)*aoiManager.gridWidth(),
				aoiManager.MinY+y*aoiManager.gridLength(),
				aoiManager.MinY+(y+1)*aoiManager.gridLength())
		}
	}
	return aoiManager
}

// 得到每个格子在 X 轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子在 Y 轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印区域的信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX:%d, MaxX:%d, cntsX:%d, minY: %d, maxY: %d, cntsY:%d\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// GetSurroundGridsByGid 根据格子 GID 得到周边九宫格格子集合
func (m *AOIManager) GetSurroundGridsByGid(gId int) (grids []*Grid) {
	// 判断传入的 gId 是否在 AOIManager 中
	if _, ok := m.grids[gId]; !ok {
		return
	}
	// 初始化返回值切片, 将当前 gid 加入到要返回的格子中
	grids = append(grids, m.grids[gId])

	idx := gId % m.CntsX
	// 根据 gId 判断左右是否有格子
	if idx > 0 {
		grids = append(grids, m.grids[gId-1])
	}
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gId+1])
	}

	// 根据 gId 得到当前格子的编号 --idx = id % nx
	// 判断 idx 编号左右是否有格子，分别存放在 gidX 和 gidY 中
	gridsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gridsX = append(gridsX, v.GID)
	}

	for _, v := range gridsX {
		// 得到当前格子 id 的 y 轴的编号 idy = id / ny
		idy := v / m.CntsY
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])
		}
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}
	}
	return grids
}

// GetPidsByPos 通过坐标得到周边九宫格全部的 playerIds
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {

	// 得到当前位置的 gid
	gid := m.GetGidByPos(x, y)
	// 根据 gid 得到周边九宫格
	grids := m.GetSurroundGridsByGid(gid)
	// 获取九宫格中的 gid
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIds()...)
		//fmt.Printf("====> grid id: %d, pids: %v=====\n", grid.GID, grid.GetPlayerIds())
	}
	return playerIDs
}

// GetGidByPos 通过坐标得到周边九宫格全部的 playerIds
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CntsX + idx
}

// AddPidToGrid 添加一个 playerId 到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// RemovePidFromGrid  移除一个格子中的 playerId
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

// GetPidsByGid 通过 GID 获取全部的 PlayerId
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIds()
	return playerIDs
}

// AddToGridByPos 通过坐标将 Player 添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.AddPidToGrid(pID, gID)
}

// RemoveFromGridByPos 通过坐标把一个 Player 从一个格子删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.RemovePidFromGrid(pID, gID)
}
