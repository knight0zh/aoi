package aoi

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"
)

/**
         0   50   100   150    200
     -----------------------------
0    |   0    1     2     3      4
50   |   5    6     7     8      9
100  |  10   11    12    13     14
150  |  15   16    17    18     19
200  |  20   21    22    23     24

*/
func TestGridManger_GetSurroundGrids(t *testing.T) {
	aol := NewGridManger(0, 0, 250, 5)
	manger := aol.(*GridManger)
	tests := []struct {
		x, y float64
		want []int
	}{
		{
			x: 0, y: 0,
			want: []int{0, 1, 5, 6},
		},
		{
			x: 150, y: 0,
			want: []int{2, 3, 4, 7, 8, 9},
		},
		{
			x: 50, y: 50,
			want: []int{0, 1, 2, 5, 6, 7, 10, 11, 12},
		},
		{
			x: 200, y: 100,
			want: []int{8, 9, 13, 14, 18, 19},
		},
		{
			x: 200, y: 200,
			want: []int{18, 19, 23, 24},
		},
	}

	for _, tt := range tests {
		ID := manger.getGIDByPos(tt.x, tt.y)
		grids := manger.getSurroundGrids(ID)
		gID := make([]int, 0)
		for _, grid := range grids {
			gID = append(gID, grid.GID)
		}
		sort.Ints(gID)
		assert.Equal(t, tt.want, gID)
	}
}

func TestNewGridManger(t *testing.T) {
	aol := NewGridManger(0, 0, 250, 5)
	manger := aol.(*GridManger)
	entities := []*Entity{
		{
			X: 0, Y: 0, Key: "a",
		},
		{
			X: 50, Y: 0, Key: "b",
		},
		{
			X: 100, Y: 0, Key: "c",
		},
		{
			X: 50, Y: 0, Key: "d",
		},
		{
			X: 50, Y: 50, Key: "e",
		},
		{
			X: 50, Y: 100, Key: "f",
		},
		{
			X: 100, Y: 0, Key: "g",
		},
		{
			X: 100, Y: 50, Key: "h",
		},
		{
			X: 100, Y: 100, Key: "i",
		},
	}

	for _, entity := range entities {
		manger.Add(entity.X, entity.Y, entity.Key)
	}

	search := manger.Search(50, 50)
	result := make([]string, 0)
	for _, entity := range search {
		result = append(result, entity)
	}
	sort.Strings(result)
	assert.Equal(t, []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}, result)

	manger.Delete(100, 100, "i")
	search2 := manger.Search(50, 50)
	result2 := make([]string, 0)
	for _, entity := range search2 {
		result2 = append(result2, entity)
	}
	sort.Strings(result2)
	assert.Equal(t, []string{"a", "b", "c", "d", "e", "f", "g", "h"}, result2)
}

func BenchmarkGridManger(b *testing.B) {
	var wg sync.WaitGroup
	aol := NewGridManger(0, 0, 1024, 16)
	manger := aol.(*GridManger)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		wg.Add(30000)
		for j := 0; j < 10000; j++ {
			go func() {
				manger.Add(
					float64(rand.Intn(10)*10+rand.Intn(10)),
					float64(rand.Intn(10)*10+rand.Intn(10)),
					fmt.Sprintf("player%d", rand.Intn(100)),
				)
				wg.Done()
			}()

			go func() {
				manger.Delete(
					float64(rand.Intn(10)*10+rand.Intn(10)),
					float64(rand.Intn(10)*10+rand.Intn(10)),
					fmt.Sprintf("player%d", rand.Intn(100)),
				)
				wg.Done()
			}()

			go func() {
				manger.Search(
					float64(rand.Intn(10)*10+rand.Intn(10)),
					float64(rand.Intn(10)*10+rand.Intn(10)),
				)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
