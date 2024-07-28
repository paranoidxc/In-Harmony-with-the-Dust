package game

import (
	"math"
)

// draw a line
func bresenham(start Pos, end Pos) []Pos {
	result := make([]Pos, 0)
	steep := math.Abs(float64(end.Y-start.Y)) > math.Abs(float64(end.X-start.X))

	if steep {
		start.X, start.Y = start.Y, start.X
		end.X, end.Y = end.Y, end.X
	}

	if start.X > end.X {
		start.X, end.X = end.X, start.X
		start.Y, end.Y = end.Y, start.Y
	}

	deltaX := end.X - start.X
	deltaY := int(math.Abs(float64(end.Y - start.Y)))

	err := 0
	y := start.Y
	ystep := 1
	if start.Y >= end.Y {
		ystep = -1
	}

	for x := start.X; x < end.X; x++ {
		if steep {
			result = append(result, Pos{y, x})
		} else {
			result = append(result, Pos{x, y})
		}
		err += deltaY
		if 2*err >= deltaY {
			y += ystep
			err -= deltaX
		}
	}

	return result
}
