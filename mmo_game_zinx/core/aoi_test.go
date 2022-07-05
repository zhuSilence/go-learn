package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	// 初始化 AOIManager
	aoiMgr := NewAOIManager(100, 300, 4, 200, 450, 5)

	fmt.Println(aoiMgr)

}

func TestAOIManager_GetSurroundGridsByGid(t *testing.T) {
	// 初始化 AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	for gid, _ := range aoiMgr.grids {
		grids := aoiMgr.GetSurroundGridsByGid(gid)
		fmt.Println("gid :", gid, "grids len :", len(grids))
		gIds := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIds = append(gIds, grid.GID)
		}
		fmt.Println("surround grid ids are", gIds)
	}
}
