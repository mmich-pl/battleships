package board_validation

import (
	"battleships/internal/app"
	"fmt"
	"reflect"
	"sort"
)

const (
	battleship = 4
	cruiser    = 3
	destroyer  = 2
	submarine  = 1
	boardSize  = 10
)

func CountOccurrences(matrix [][]int) []int {
	counts := make(map[int]int)
	for _, row := range matrix {
		for _, num := range row {
			if num == 0 {
				continue
			}
			counts[num]++
		}
	}
	v := make([]int, 0, len(counts))

	for _, value := range counts {
		v = append(v, value)
	}
	sort.Ints(v)
	return v
}

func validateShipPlacement(fleet []string) (bool, string) {
	if len(fleet) != 20 {
		return false, "ship overlap or missing"
	}

	shipQuantities := map[int]int{
		battleship: 1,
		cruiser:    2,
		destroyer:  3,
		submarine:  4,
	}

	var shipList []int
	for ship, count := range shipQuantities {
		for i := 0; i < count; i++ {
			shipList = append(shipList, ship)
		}
	}
	sort.Ints(shipList)

	board := make([][]int, boardSize)
	for i := range board {
		board[i] = make([]int, boardSize)
	}

	for _, coords := range fleet {
		x, y, err := app.MapCoords(coords)
		if err != nil {
			return false, fmt.Sprintf("invalid coordinate: %s", err)
		}
		board[x][y] = 1
	}

	blobs := ConnectedComponentLabeling(board)
	occurrences := CountOccurrences(blobs)
	if reflect.DeepEqual(shipList, occurrences) {
		return true, "board valid"
	}

	return false, "board invalid"
}
