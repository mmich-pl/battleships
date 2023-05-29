package board_utils

import (
	"battleships/internal/utils"
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

func countOccurrences(matrix [][]int) []int {
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

func ValidateShipPlacement(fleet []string) (bool, string) {
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
		x, y, err := utils.MapCoords(coords)
		if err != nil {
			return false, fmt.Sprintf("invalid coordinate: %s", err)
		}
		board[x][y] = 1
	}

	blobs := ConnectedComponentLabeling(board)
	occurrences := countOccurrences(blobs)
	if reflect.DeepEqual(shipList, occurrences) {
		return true, "board_utils valid"
	}

	return false, "board_utils invalid"
}