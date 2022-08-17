package aoi

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func Test_FindQuadrant(t *testing.T) {
	aoi := NewQuadTree(0, 0, 100)
	tree := aoi.(*QuadTree)
	tree.cutNode()

	tests := []struct {
		x, y float64
		want int
	}{
		{
			x: 49.9, y: 49.9, want: leftUp,
		},
		{
			x: 50, y: 50, want: rightDown,
		},
		{
			x: 49.9, y: 50, want: leftDown,
		},
		{
			x: 50, y: 49.9, want: rightUp,
		},
	}

	for _, tt := range tests {
		d := tree.findSonQuadrant(tt.x, tt.y)
		assert.Equal(t, tt.want, d)
	}

	// 再次分割
	tree.Child[rightUp].cutNode()
	tests2 := []struct {
		x, y float64
		want int
	}{
		{
			x: 74.9, y: 24.9, want: leftUp,
		},
		{
			x: 75, y: 25, want: rightDown,
		},
		{
			x: 74.9, y: 25, want: leftDown,
		},
		{
			x: 75, y: 24.9, want: rightUp,
		},
	}
	for _, tt := range tests2 {
		d := tree.Child[rightUp].findSonQuadrant(tt.x, tt.y)
		assert.Equal(t, tt.want, d)
	}
}

func Test_NeedCut(t *testing.T) {
	aoi := NewQuadTree(0, 0, 100)
	tree := aoi.(*QuadTree)

	tree.maxCap = 2 // 超过两人节点分裂
	assert.Equal(t, false, tree.needCut())
	tree.Add(60.9, 24.9, "player1")

	assert.Equal(t, false, tree.needCut())
	tree.Add(25, 25, "player2")

	assert.Equal(t, true, tree.needCut())
}

func TestNode_Search(t *testing.T) {
	aoi := NewQuadTree(0, 0, 100)
	tree := aoi.(*QuadTree)
	tree.maxCap = 2 // 超过两人节点分裂
	tree.radius = 5

	tree.Add(60.9, 24.9, "player1")
	tree.Add(25, 25, "player2")

	// 查询player1附近
	entities := tree.Search(60.9, 24.9)
	assert.Equal(t, 2, len(entities), "player1 player2")

	// 当出现第三个玩家超过节点最大容量产生分裂
	tree.Add(99, 24, "player3")

	// 查询player1附近
	entities = tree.Search(60.9, 24.9)
	assert.Equal(t, 2, len(entities), "player1 player3")

	// 添加第四个玩家
	tree.Add(72, 23, "player4")

	// 查询player1附近
	entities = tree.Search(60.9, 24.9)
	assert.Equal(t, 2, len(entities), "player1 player4")

	// 查询player2附近
	entities = tree.Search(25, 25)
	assert.Equal(t, 1, len(entities), "player2")

	// 添加第五个玩家
	tree.Add(49.9, 49.9, "player5")

	// 查询player2附近
	entities = tree.Search(25, 25)
	assert.Equal(t, 2, len(entities), "player2 player5")

	// 移除player5
	tree.Delete(49.9, 49.9, "player5")

	// 查询player2附近
	entities = tree.Search(25, 25)
	assert.Equal(t, 1, len(entities), "player2")
}

func BenchmarkQuadtree(b *testing.B) {
	var wg sync.WaitGroup
	aoi := NewQuadTree(0, 0, 1024)
	tree := aoi.(*QuadTree)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < b.N; i++ {
		wg.Add(30000)
		for j := 0; j < 10000; j++ {
			go func() {
				tree.Add(
					float64(rand.Intn(10)*10+rand.Intn(10)),
					float64(rand.Intn(10)*10+rand.Intn(10)),
					fmt.Sprintf("player%d", rand.Intn(100)),
				)
				wg.Done()
			}()

			go func() {
				tree.Delete(
					float64(rand.Intn(10)*10+rand.Intn(10)),
					float64(rand.Intn(10)*10+rand.Intn(10)),
					fmt.Sprintf("player%d", rand.Intn(100)),
				)
				wg.Done()
			}()

			go func() {
				tree.Search(
					float64(rand.Intn(10)*10+rand.Intn(10)),
					float64(rand.Intn(10)*10+rand.Intn(10)),
				)
				wg.Done()
			}()
		}
	}
}
