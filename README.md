# AOI (Area of Interest) Library

This library provides implementations of Area of Interest algorithms for spatial partitioning. Currently, it includes implementations for the following algorithms:

1. **九宫格 (Grid Manager)**
   - A simple grid-based AOI algorithm dividing the area into a grid of cells and associating entities with the corresponding grid cells.

2. **四叉树 (Quadtree)**
   - A hierarchical spatial partitioning algorithm dividing the area into four quadrants recursively, optimizing the search for entities within a specified range.

## Usage:

### 九宫格 (Grid Manager)

```go
// Example Usage:
aoiManager := NewGridManager(startX, startY, areaWidth, gridCount)
aoiManager.Add(x, y, "Entity1")
aoiManager.Delete(x, y, "Entity1")
result := aoiManager.Search(x, y)

// Example Usage:
quadTree := NewQuadTree(startX, startY, areaWidth)
quadTree.Add(x, y, "Entity1")
quadTree.Delete(x, y, "Entity1")
result := quadTree.Search(x, y)
```

## Features:
- Both implementations support adding, deleting, and searching for entities within a specified area of interest.
- The Grid Manager uses a simple grid-based approach, while the Quadtree provides a hierarchical and optimized solution for larger and dynamic environments.

## TODO:
Implement additional commonly used AOI algorithms:
- R-树 (R-tree)
- 六边形网格 (Hexagonal Grid)
- 基于事件的算法 (Event-driven Approaches)