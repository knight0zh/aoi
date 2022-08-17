package aoi

import "sync"

const (
	leftUp int = iota
	rightUp
	leftDown
	rightDown

	maxCap  = 500 // 节点最大容量
	maxDeep = 4   // 节点最大深度
	radius  = 16  // 视野半径
)

type Node struct {
	Leaf      bool      // 是否为叶子节点
	Deep      int       // 深度
	AreaWidth float64   // 格子宽度(长=宽)
	XStart    float64   // 起始范围
	YStart    float64   // 起始范围
	Tree      *QuadTree // 树指针
	Child     [4]*Node  // 子节点
	Entities  *sync.Map // 实体
}

type QuadTree struct {
	maxCap, maxDeep int
	radius          float64
	mPool           sync.Pool
	*Node
}

func NewSonNode(xStart, yStart float64, parent *Node) *Node {
	son := &Node{
		Leaf:      true,
		Deep:      parent.Deep + 1,
		AreaWidth: parent.AreaWidth / 2,
		XStart:    xStart,
		YStart:    yStart,
		Tree:      parent.Tree,
		Entities:  parent.Tree.mPool.Get().(*sync.Map),
	}

	return son
}

// canCut 检查节点是否可以分割
func (n *Node) canCut() bool {
	if n.XStart+n.AreaWidth/2 > 0 && n.YStart+n.AreaWidth/2 > 0 {
		return true
	}
	return false
}

// needCut 检查节点是否需要分割
func (n *Node) needCut() bool {
	lens := 0
	n.Entities.Range(func(key, value interface{}) bool {
		lens++
		return true
	})
	return lens+1 > n.Tree.maxCap && n.Deep+1 <= n.Tree.maxDeep && n.canCut()
}

// intersects 检查坐标是否在节点范围内
func (n *Node) intersects(x, y float64) bool {
	if n.XStart <= x && x < n.XStart+n.AreaWidth && n.YStart <= y && y < n.YStart+n.AreaWidth {
		return true
	}
	return false
}

// findSonQuadrant 根据坐标寻找子节点的方位
func (n *Node) findSonQuadrant(x, y float64) int {
	if x < n.Child[rightDown].XStart {
		if y < n.Child[rightDown].YStart {
			return leftUp
		}
		return leftDown
	}
	if y < n.Child[rightDown].YStart {
		return rightUp
	}
	return rightDown
}

// cutNode 分割节点
func (n *Node) cutNode() {
	n.Leaf = false
	half := n.AreaWidth / 2

	n.Child[leftUp] = NewSonNode(n.XStart, n.YStart, n)
	n.Child[rightUp] = NewSonNode(n.XStart+half, n.YStart, n)
	n.Child[leftDown] = NewSonNode(n.XStart, n.YStart+half, n)
	n.Child[rightDown] = NewSonNode(n.XStart+half, n.YStart+half, n)

	// 将实体迁移到对应子节点
	n.Entities.Range(func(k, v interface{}) bool {
		entity := v.(*Entity)
		for _, node := range n.Child {
			if node.intersects(entity.X, entity.Y) {
				node.Entities.Store(entity.Key, entity)
			}
		}
		n.Entities.Delete(k)
		return true
	})

	n.Tree.mPool.Put(n.Entities)
	n.Entities = nil
}

func NewQuadTree(xStart, yStart, width float64) AOI {
	basicNode := &Node{
		Leaf:      true,
		Deep:      1,
		AreaWidth: width,
		XStart:    xStart,
		YStart:    yStart,
		Child:     [4]*Node{},
	}
	tree := &QuadTree{
		maxDeep: maxDeep,
		maxCap:  maxCap,
		radius:  radius,
		Node:    basicNode,
	}
	tree.mPool.New = func() interface{} {
		return &sync.Map{}
	}
	basicNode.Tree = tree
	basicNode.Entities = tree.mPool.Get().(*sync.Map)
	return tree
}

func (n *Node) Add(x, y float64, name string) {
	// 判断是否需要分割
	if n.Leaf && n.needCut() {
		n.cutNode()
	}

	// 非叶子节点往下递归
	if !n.Leaf {
		n.Child[n.findSonQuadrant(x, y)].Add(x, y, name)
		return
	}

	entity := entityPool.Get().(*Entity)
	entity.X = x
	entity.Y = y
	entity.Key = name

	// 叶子节点进行存储
	n.Entities.Store(entity.Key, entity)
}

func (n *Node) Delete(x, y float64, name string) {
	if !n.Leaf {
		n.Child[n.findSonQuadrant(x, y)].Delete(x, y, name)
		return
	}

	if entity, ok := n.Entities.Load(name); ok {
		n.Entities.Delete(name)
		entityPool.Put(entity)
	}
}

func (n *Node) Search(x, y float64) []string {
	result := resultPool.Get().([]string)
	defer func() {
		result = result[:0]
		resultPool.Put(result)
	}()
	n.search(x, y, &result)
	return result
}

func (n *Node) search(x, y float64, result *[]string) {
	if !n.Leaf {
		minX, maxX := x-n.Tree.radius, x+n.Tree.radius
		minY, maxY := y-n.Tree.radius, y+n.Tree.radius

		for _, son := range n.Child {
			if son.intersects(minX, minY) || son.intersects(maxX, minY) ||
				son.intersects(minX, maxY) || son.intersects(maxX, maxY) {
				son.search(x, y, result)
			}
		}
		return
	}

	n.Entities.Range(func(key, value interface{}) bool {
		*result = append(*result, value.(*Entity).Key)
		return true
	})
	return
}
