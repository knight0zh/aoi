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
	player1 := &Entity{
		X:    60.9,
		Y:    24.9,
		Name: "player1",
	}
	tree.Add(player1)

	assert.Equal(t, false, tree.needCut())
	player2 := &Entity{
		X:    25,
		Y:    25,
		Name: "player2",
	}
	tree.Add(player2)

	assert.Equal(t, true, tree.needCut())
}

func TestNode_Search(t *testing.T) {
	aoi := NewQuadTree(0, 0, 100)
	tree := aoi.(*QuadTree)
	tree.maxCap = 2 // 超过两人节点分裂
	tree.radius = 5
	player1 := &Entity{
		X:    60.9,
		Y:    24.9,
		Name: "player1",
	}
	tree.Add(player1)

	player2 := &Entity{
		X:    25,
		Y:    25,
		Name: "player2",
	}
	tree.Add(player2)
	entities := tree.Search(player1)
	assert.Equal(t, 2, len(entities), "player1 player2")

	// 当出现第三个玩家超过节点最大容量产生分裂
	player3 := &Entity{
		X:    99,
		Y:    24,
		Name: "player3",
	}
	tree.Add(player3)
	entities = tree.Search(player1)
	assert.Equal(t, 2, len(entities), "player1 player3")

	// 添加第四个玩家
	player4 := &Entity{
		X:    72,
		Y:    23,
		Name: "player4",
	}
	tree.Add(player4)
	entities = tree.Search(player1)
	assert.Equal(t, 2, len(entities), "player1 player4")

	entities = tree.Search(player2)
	assert.Equal(t, 1, len(entities), "player2")

	// 添加第五个玩家
	player5 := &Entity{
		X:    49.9,
		Y:    49.9,
		Name: "player5",
	}
	tree.Add(player5)
	entities = tree.Search(player2)
	assert.Equal(t, 2, len(entities), "player2 player5")

	tree.Delete(player5)
	entities = tree.Search(player2)
	assert.Equal(t, 1, len(entities), "player2")
}

func BenchmarkQuadtree(b *testing.B) {
	var wg sync.WaitGroup
	aoi := NewQuadTree(0, 0, 100)
	tree := aoi.(*QuadTree)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < b.N; i++ {
		wg.Add(3000)
		for j := 0; j < 1000; j++ {
			go func() {
				tree.Add(&Entity{
					X:    float64(rand.Intn(5) * 10),
					Y:    float64(rand.Intn(5) * 10),
					Name: fmt.Sprintf("player%d", rand.Intn(50)),
				})
				wg.Done()
			}()

			go func() {
				tree.Delete(&Entity{
					X:    float64(rand.Intn(5) * 10),
					Y:    float64(rand.Intn(5) * 10),
					Name: fmt.Sprintf("player%d", rand.Intn(50)),
				})
				wg.Done()
			}()

			go func() {
				tree.Search(&Entity{
					X: float64(rand.Intn(5) * 10),
					Y: float64(rand.Intn(5) * 10),
				})
				wg.Done()
			}()
		}
	}
}