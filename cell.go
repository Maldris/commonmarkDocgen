package docgen

import (
	"fmt"
	"reflect"
)

// Cell is used to provide nested data struct access in templates
type Cell struct{ v interface{} }

func NewCell(v ...interface{}) (*Cell, error) {
	switch len(v) {
	case 0:
		return new(Cell), nil
	case 1:
		return &Cell{v[0]}, nil
	default:
		return nil, fmt.Errorf("wrong number of args: want 0 or 1, got %v", len(v))
	}
}

func (c *Cell) Set(v interface{}) *Cell { c.v = v; return c }
func (c *Cell) Get() interface{}        { return c.v }

// eq must be overwritten to support cell
func eq(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case string, int, int64, byte, float32, float64:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}

	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}
