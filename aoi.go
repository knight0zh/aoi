package aoi

import "sync"

// AOI (Area of Interest) represents an interface for managing entities within a specific area.
type AOI interface {
	Add(x, y float64, name string)         // Add an entity to the AOI
	Delete(x, y float64, name string)      // Delete an entity from the AOI
	Search(x, y float64) (result []string) // Search for entities within a specified range
}

// Entity represents an object with coordinates and a key.
type Entity struct {
	X, Y float64
	Key  string
}

var (
	resultPool sync.Pool // Pool for recycling result slices
	entityPool sync.Pool // Pool for recycling Entity objects
)

func init() {
	// Initialize the resultPool to recycle result slices
	resultPool.New = func() interface{} {
		return make([]string, 0, 500)
	}

	// Initialize the entityPool to recycle Entity objects
	entityPool.New = func() interface{} {
		return &Entity{}
	}
}
