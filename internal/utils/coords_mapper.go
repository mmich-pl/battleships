package utils

import (
	"fmt"
	"strconv"
)

// Parses coordinates to two integers that represents board square in matrix
func MapCoords(coordinate string) (int, int, error) {
	if len(coordinate) == 0 {
		return -1, -1, fmt.Errorf("coordinate is empty")
	}
	column := coordinate[0]
	if 'A' > column || column > 'J' {
		return -1, -1, fmt.Errorf("coordinate out of bound: expected column in bounds [A-J]")
	}

	x := int(column - 'A')
	y, err := strconv.Atoi(coordinate[1:])
	if err != nil {
		return -1, -1, fmt.Errorf("wrong coordinates format: %w", err)
	} else if y < 1 || y > 10 {
		return -1, -1, fmt.Errorf("coordinate [%s] out of bound: expected row in bounds [1-10]", coordinate)
	}
	return x, y - 1, nil
}
