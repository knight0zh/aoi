package aoi

type AOI interface {
	Add(entity *Entity)                       // 添加实体
	Delete(entity *Entity)                    // 移除实体
	Search(entity *Entity) (result []*Entity) // 范围查询
}

type Entity struct {
	X, Y float64
	Name string
}
