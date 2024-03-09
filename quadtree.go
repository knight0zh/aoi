package aoi

import "sync"

const (
	leftUp int = iota
	rightUp
	leftDown
	rightDown

	maxCap  = 500 // Maximum capacity of a node
	maxDeep = 4   // Maximum depth of the quadtree
	radius  = 16  // Field of view radius
)

// Node represents a node in the quadtree.
type Node struct {
	Leaf      bool      // Indicates whether the node is a leaf node
	Deep      int       // Depth of the node in the quadtree
	AreaWidth float64   // Width of the grid (assuming square grids)
	XStart    float64   // Starting X-coordinate of the node's area
	YStart    float64   // Starting Y-coordinate of the node's area
	Tree      *QuadTree // Pointer to the quadtree
	Child     [4]*Node  // Child nodes (quadrants)
	Entities  *sync.Map // Entities within the node
}

// QuadTree represents a quadtree data structure for spatial partitioning.
type QuadTree struct {
	maxCap, maxDeep int
	radius          float64
	mPool           sync.Pool
	*Node
}

// NewSonNode creates a new child node with the specified parameters.
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

// canCut checks whether the node can be split.
func (n *Node) canCut() bool {
	if n.XStart+n.AreaWidth/2 > 0 && n.YStart+n.AreaWidth/2 > 0 {
		return true
	}
	return false
}

// needCut checks whether the node needs to be split.
func (n *Node) needCut() bool {
	lens := 0
	n.Entities.Range(func(key, value interface{}) bool {
		lens++
		return true
	})
	return lens+1 > n.Tree.maxCap && n.Deep+1 <= n.Tree.maxDeep && n.canCut()
}

// intersects checks if the coordinates are within the node's range.
func (n *Node) intersects(x, y float64) bool {
	if n.XStart <= x && x < n.XStart+n.AreaWidth && n.YStart <= y && y < n.YStart+n.AreaWidth {
		return true
	}
	return false
}

// findSonQuadrant finds the quadrant of a child node based on coordinates.
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

// cutNode splits the node into four child nodes.
func (n *Node) cutNode() {
	n.Leaf = false
	half := n.AreaWidth / 2

	n.Child[leftUp] = NewSonNode(n.XStart, n.YStart, n)
	n.Child[rightUp] = NewSonNode(n.XStart+half, n.YStart, n)
	n.Child[leftDown] = NewSonNode(n.XStart, n.YStart+half, n)
	n.Child[rightDown] = NewSonNode(n.XStart+half, n.YStart+half, n)

	// Move entities to the corresponding child nodes
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

// NewQuadTree initializes a new QuadTree with the specified parameters.
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

// Add adds an entity to the quadtree based on its coordinates.
func (n *Node) Add(x, y float64, name string) {
	// Check if splitting is required
	if n.Leaf && n.needCut() {
		n.cutNode()
	}

	// Recursively add to non-leaf nodes
	if !n.Leaf {
		n.Child[n.findSonQuadrant(x, y)].Add(x, y, name)
		return
	}

	entity := entityPool.Get().(*Entity)
	entity.X = x
	entity.Y = y
	entity.Key = name

	// Store in leaf node
	n.Entities.Store(entity.Key, entity)
}

// Delete removes an entity from the quadtree based on its coordinates.
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

// Search retrieves a list of entity keys within the specified coordinates' range.
func (n *Node) Search(x, y float64) []string {
	result := resultPool.Get().([]string)
	defer func() {
		result = result[:0]
		resultPool.Put(result)
	}()
	n.search(x, y, &result)
	return result
}

// search recursively searches for entities within the specified coordinates' range.
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

	// Collect entity keys within the leaf node
	n.Entities.Range(func(key, value interface{}) bool {
		*result = append(*result, value.(*Entity).Key)
		return true
	})
	return
}
