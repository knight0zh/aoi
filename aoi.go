package aoi

import "sync"

type AOI interface {
	Add(x, y float64, name string)         // 添加实体
	Delete(x, y float64, name string)      // 移除实体
	Search(x, y float64) (result []string) // 范围查询
}

type Entity struct {
	X, Y float64
	Key  string
}

var (
	resultPool sync.Pool
	entityPool sync.Pool
)

func init() {
	resultPool.New = func() interface{} {
		return make([]string, 0, 500)
	}
	entityPool.New = func() interface{} {
		return &Entity{}
	}
}
