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
		entity *Entity
		want   []int
	}{
		{
			entity: &Entity{
				X: 0, Y: 0,
			},
			want: []int{0, 1, 5, 6},
		},
		{
			entity: &Entity{
				X: 150, Y: 0,
			},
			want: []int{2, 3, 4, 7, 8, 9},
		},
		{
			entity: &Entity{
				X: 50, Y: 50,
			},
			want: []int{0, 1, 2, 5, 6, 7, 10, 11, 12},
		},
		{
			entity: &Entity{
				X: 200, Y: 100,
			},
			want: []int{8, 9, 13, 14, 18, 19},
		},
		{
			entity: &Entity{
				X: 200, Y: 200,
			},
			want: []int{18, 19, 23, 24},
		},
	}

	for _, tt := range tests {
		ID := manger.GetGIDByPos(tt.entity)
		grids := manger.GetSurroundGrids(ID)
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
			X: 0, Y: 0, Name: "a",
		},
		{
			X: 50, Y: 0, Name: "b",
		},
		{
			X: 100, Y: 0, Name: "c",
		},
		{
			X: 50, Y: 0, Name: "d",
		},
		{
			X: 50, Y: 50, Name: "e",
		},
		{
			X: 50, Y: 100, Name: "f",
		},
		{
			X: 100, Y: 0, Name: "g",
		},
		{
			X: 100, Y: 50, Name: "h",
		},
		{
			X: 100, Y: 100, Name: "i",
		},
	}

	for _, entity := range entities {
		manger.Add(entity)
	}

	search := manger.Search(&Entity{X: 50, Y: 50})
	result := make([]string, 0)
	for _, entity := range search {
		result = append(result, entity.Name)
	}
	sort.Strings(result)
	assert.Equal(t, []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}, result)

	manger.Delete(&Entity{X: 100, Y: 100, Name: "i"})
	search2 := manger.Search(&Entity{X: 50, Y: 50})
	result2 := make([]string, 0)
	for _, entity := range search2 {
		result2 = append(result2, entity.Name)
	}
	sort.Strings(result2)
	assert.Equal(t, []string{"a", "b", "c", "d", "e", "f", "g", "h"}, result2)
}

func BenchmarkGridManger(b *testing.B) {
	var wg sync.WaitGroup
	aol := NewGridManger(0, 0, 256, 16)
	manger := aol.(*GridManger)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		wg.Add(30000)
		for j := 0; j < 10000; j++ {
			go func() {
				manger.Add(&Entity{
					X:    float64(rand.Intn(5) * 10),
					Y:    float64(rand.Intn(5) * 10),
					Name: fmt.Sprintf("player%d", rand.Intn(50)),
				})
				wg.Done()
			}()

			go func() {
				manger.Delete(&Entity{
					X:    float64(rand.Intn(5) * 10),
					Y:    float64(rand.Intn(5) * 10),
					Name: fmt.Sprintf("player%d", rand.Intn(50)),
				})
				wg.Done()
			}()

			go func() {
				manger.Search(&Entity{
					X: float64(rand.Intn(5) * 10),
					Y: float64(rand.Intn(5) * 10),
				})
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
