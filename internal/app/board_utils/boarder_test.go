package board_utils

import (
	. "github.com/grupawp/warships-gui/v2"
	"testing"
)

func TestMarkBorder(t *testing.T) {
	board := &[10][10]State{
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Hit, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Hit, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Hit, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Hit, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
	}

	expectedBoard := &[10][10]State{
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Miss, Miss, Miss, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Miss, Hit, Miss, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Miss, Hit, Miss, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Miss, Hit, Miss, Miss, Miss, Empty, Empty, Empty, Empty},
		{Empty, Miss, Miss, Miss, Hit, Miss, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Miss, Miss, Miss, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
	}

	MarkBorder(board, 2, 2)
	MarkBorder(board, 5, 4)

	// Check if the board matches the expected board after marking the border
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if board[i][j] != expectedBoard[i][j] {
				t.Errorf("Board mismatch at row %d, column %d (expected: %v, got:%v)", i, j, expectedBoard[i][j], board[i][j])
			}
		}
	}
}
