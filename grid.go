package aoi

import "sync"

var (
	// 分别将这8个方向的方向向量按顺序写入x, y的分量数组
	dx = []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy = []int{-1, 0, 1, -1, 1, -1, 0, 1}
)

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
	pool      sync.Pool
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
	manager.pool.New = func() interface{} {
		return make([]*Grid, 0, 9)
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

// getGIDByPos 通过横纵坐标获取对应的格子ID
func (g *GridManger) getGIDByPos(x, y float64) int {
	gx := (int(x) - g.StartX) / g.gridWidth()
	gy := (int(y) - g.StartY) / g.gridWidth()

	return gy*g.GridCount + gx
}

// getSurroundGrids 根据格子的gID得到当前周边的九宫格信息
func (g *GridManger) getSurroundGrids(gID int) []*Grid {
	grids := g.pool.Get().([]*Grid)
	defer func() {
		grids = grids[:0]
		g.pool.Put(grids)
	}()

	if _, ok := g.grids[gID]; !ok {
		return grids
	}
	grids = append(grids, g.grids[gID])
	// 根据gID, 得到格子所在的坐标
	x, y := gID%g.GridCount, gID/g.GridCount

	for i := 0; i < 8; i++ {
		newX := x + dx[i]
		newY := y + dy[i]

		if newX >= 0 && newX < g.GridCount && newY >= 0 && newY < g.GridCount {
			grids = append(grids, g.grids[newY*g.GridCount+newX])
		}
	}

	return grids
}

func (g *GridManger) Add(x, y float64, key string) {
	entity := entityPool.Get().(*Entity)
	entity.X = x
	entity.Y = y
	entity.Key = key

	ID := g.getGIDByPos(x, y)
	grid := g.grids[ID]
	grid.Entities.Store(key, entity)
}

func (g *GridManger) Delete(x, y float64, key string) {
	ID := g.getGIDByPos(x, y)
	grid := g.grids[ID]

	if entity, ok := grid.Entities.Load(key); ok {
		grid.Entities.Delete(key)
		entityPool.Put(entity)
	}
}

func (g *GridManger) Search(x, y float64) []string {
	result := resultPool.Get().([]string)
	defer func() {
		result = result[:0]
		resultPool.Put(result)
	}()
	ID := g.getGIDByPos(x, y)
	grids := g.getSurroundGrids(ID)
	for _, grid := range grids {
		grid.Entities.Range(func(_, value interface{}) bool {
			result = append(result, value.(*Entity).Key)
			return true
		})
	}

	return result
}
