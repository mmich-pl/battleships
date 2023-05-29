package board_utils

import (
	"battleships/internal/utils"
	"fmt"
	"reflect"
	"sort"
)

var (
	ShipQuantities = map[int]int{
		battleship: 1,
		cruiser:    2,
		destroyer:  3,
		submarine:  4,
	}
)

const (
	battleship = 4
	cruiser    = 3
	destroyer  = 2
	submarine  = 1
	BoardSize  = 10
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

func ValidateShipPlacement(fleet []string) (bool, string) {
	if len(fleet) != 20 {
		return false, "ship overlap or missing"
	}

	var shipList []int
	for ship, count := range ShipQuantities {
		for i := 0; i < count; i++ {
			shipList = append(shipList, ship)
		}
	}
	sort.Ints(shipList)

	board, msg := MapCoordToBoard(fleet)
	if msg != "" {
		return false, msg
	}

	blobs := ConnectedComponentLabeling(board)
	occurrences := CountOccurrences(blobs)
	if reflect.DeepEqual(shipList, occurrences) {
		return true, "board_utils valid"
	}

	return false, "board_utils invalid"
}

func MapCoordToBoard(coordinates []string) ([][]int, string) {
	board := make([][]int, BoardSize)
	for i := range board {
		board[i] = make([]int, BoardSize)
	}

	for _, coords := range coordinates {
		x, y, err := utils.MapCoords(coords)
		if err != nil {
			return nil, fmt.Sprintf("invalid coordinate: %s", err)
		}
		board[x][y] = 1
	}
	return board, ""
}
