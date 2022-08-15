package aoi

import "sync"

// Grid 格子
type Grid struct {
	GID      int      //格子ID
	Entities sync.Map //当前格子内的实体
}

// GridManger AOI九宫格实现矩形
type GridManger struct {
	StartX    int // X区域左边界坐标
	StartY    int // Y区域上边界坐标
	AreaWidth int // 格子宽度(长=宽)
	GridCount int // 格子数量
	grids     map[int]*Grid
}

func NewGrid(gid int) *Grid {
	return &Grid{
		GID: gid,
	}
}

func NewGridManger(startX, startY, areaWidth, gridCount int) AOI {
	manager := &GridManger{
		StartX:    startX,
		StartY:    startY,
		AreaWidth: areaWidth,
		GridCount: gridCount,
		grids:     make(map[int]*Grid),
	}

	for y := 0; y < gridCount; y++ {
		for x := 0; x < gridCount; x++ {
			//格子编号：ID = IDy *nx + IDx  (利用格子坐标得到格子编号)
			gID := y*gridCount + x
			manager.grids[gID] = NewGrid(gID)
		}
	}

	return manager
}

func (g *GridManger) gridWidth() int {
	return g.AreaWidth / g.GridCount
}

// GetGIDByPos 通过横纵坐标获取对应的格子ID
func (g *GridManger) GetGIDByPos(entity *Entity) int {
	gx := (int(entity.X) - g.StartX) / g.gridWidth()
	gy := (int(entity.Y) - g.StartY) / g.gridWidth()

	return gy*g.GridCount + gx
}

// GetSurroundGrids 根据格子的gID得到当前周边的九宫格信息
func (g *GridManger) GetSurroundGrids(gID int) (grids []*Grid) {
	if _, ok := g.grids[gID]; !ok {
		return
	}
	grids = append(grids, g.grids[gID])

	// 根据gID, 得到格子所在的坐标
	x, y := gID%g.GridCount, gID/g.GridCount

	// 分别将这8个方向的方向向量按顺序写入x, y的分量数组
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	surroundGID := make([]int, 0)
	for i := 0; i < 8; i++ {
		newX := x + dx[i]
		newY := y + dy[i]

		if newX >= 0 && newX < g.GridCount && newY >= 0 && newY < g.GridCount {
			surroundGID = append(surroundGID, newY*g.GridCount+newX)
		}
	}

	for _, gID := range surroundGID {
		grids = append(grids, g.grids[gID])
	}

	return
}

func (g *GridManger) Add(entity *Entity) {
	ID := g.GetGIDByPos(entity)
	grid := g.grids[ID]
	grid.Entities.Store(entity.Name, entity)
}

func (g *GridManger) Delete(entity *Entity) {
	ID := g.GetGIDByPos(entity)
	grid := g.grids[ID]
	grid.Entities.Delete(entity.Name)
}

func (g *GridManger) Search(entity *Entity) (result []*Entity) {
	ID := g.GetGIDByPos(entity)
	grids := g.GetSurroundGrids(ID)
	for _, grid := range grids {
		grid.Entities.Range(func(_, value interface{}) bool {
			result = append(result, value.(*Entity))
			return true
		})
	}

	return
}
