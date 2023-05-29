package board_utils

import (
	gui "github.com/grupawp/warships-gui/v2"
)

func MarkBorder(board *[10][10]gui.State, x, y int) {
	var shipCoords [][2]int

	checkNeighboringShip(board, x, y, &shipCoords)
	for _, i := range shipCoords {
		drawBoarder(board, i)
	}
}

func checkNeighboringShip(board *[10][10]gui.State, row, col int, ship *[][2]int) {
	for _, i := range *ship {
		if i[0] == row && i[1] == col {
			return
		}
	}

	var restOfShip [][2]int
	*ship = append(*ship, [2]int{row, col})
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {

			// Skip the current tile itself
			if i == 0 && j == 0 {
				continue
			}

			neighborRow := row + i
			neighborCol := col + j

			if neighborRow >= 0 && neighborRow < len(board) && neighborCol >= 0 && neighborCol < len(board[0]) {
				if board[neighborRow][neighborCol] == gui.Hit {
					restOfShip = append(restOfShip, [2]int{neighborRow, neighborCol})
				}
			}
		}
	}

	for _, p := range restOfShip {
		checkNeighboringShip(board, p[0], p[1], ship)
	}
}

func max(a, b int) int {
	return a - (a-b)*(((a-b)>>31)&1)
}

func min(a, b int) int {
	return b + (a-b)*(((a-b)>>31)&1)
}

func drawBoarder(board *[10][10]gui.State, coordinate [2]int) {
	startRow := max(coordinate[0]-1, 0)
	endRow := min(coordinate[0]+1, len(board)-1)
	startCol := max(coordinate[1]-1, 0)
	endCol := min(coordinate[1]+1, len(board[0])-1)

	for i := startRow; i <= endRow; i++ {
		for j := startCol; j <= endCol; j++ {
			if board[i][j] != gui.Hit {
				board[i][j] = gui.Miss
			}
		}
	}
}
