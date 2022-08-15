package aoi

import "sync"

const (
	leftUp int = iota
	rightUp
	leftDown
	rightDown
)

type Node struct {
	AreaWidth int      // 格子宽度(长=宽)
	XStart    int      // 起始范围
	YStart    int      // 起始范围
	Deep      int      // 深度
	Leaf      bool     // 是否为叶子节点
	Parent    *Node    // 父节点
	Child     [4]*Node // 子节点
	Entities  sync.Map // 实体
}

type QuadTree struct {
	Root *Node
}

func (q QuadTree) Add(entity *Entity) {
	panic("implement me")
}

func (q QuadTree) Delete(entity *Entity) {
	panic("implement me")
}

func (q QuadTree) Search(entity *Entity) (result []*Entity) {
	panic("implement me")
}
