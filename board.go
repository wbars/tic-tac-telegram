package tic_tac_bot

import "strings"

type Board struct {
	board [][]Player
	size  int
}

type Player int

const (
	NONE   Player = iota
	FIRST  Player = iota
	SECOND Player = iota
)

func createBoard(size int) *Board {
	board := [][]Player{}
	for i := 0; i < size; i++ {
		board = append(board, make([]Player, size))
	}
	return &Board{board: board, size: size}
}

func (board Board) repr() (result string) {
	result = ""
	for row := 0; row < board.size; row++ {
		for col := 0; col < board.size; col++ {
			result += getMark(board.board[row][col])
		}
		result += "\n"
	}
	return
}

func fromRepr(repr string) Board {
	size := len(repr)
	table := [][]Player{}
	reprTable := strings.Split(repr, "\n")
	for i := 0; i < size; i++ {
		row := []Player{}
		for j := 0; j < size; j++ {
			row = append(row, fromMark(string(reprTable[row][j])))
		}
		table = append(table, row)
	}
	return Board{table, size}
}

func fromMark(mark string) Player {
	if mark == "X" {
		return FIRST
	}
	if mark == "O" {
		return SECOND
	}
	return NONE
}

func getMark(player Player) string {
	if player == FIRST {
		return "X"
	}
	if player == SECOND {
		return "O"
	}
	return "."
}

func (board Board) getWinner() Player {
	for row := 0; row < board.size; row++ {
		winner := getSliceWinner(board.board[row])
		if winner != NONE {
			return winner
		}
	}

	for col := 0; col < board.size; col++ {
		winner := getSliceWinner(getVerticalSlice(board.board, col))
		if winner != NONE {
			return winner
		}
	}

	winner := getSliceWinner(mainDiagonal(board.board))
	if winner != NONE {
		return winner
	}

	winner = getSliceWinner(reverseDiagonal(board.board))
	if winner != NONE {
		return winner
	}
	return NONE

}
func mainDiagonal(board [][]Player) (result []Player) {
	result = []Player{}
	for i := 0; i < len(board); i++ {
		result = append(result, board[i][i])
	}
	return
}
func getVerticalSlice(board [][]Player, col int) (result []Player) {
	result = []Player{}
	for row := 0; row < len(board); row++ {
		result = append(result, board[row][col])
	}
	return
}
func getSliceWinner(slice []Player) Player {
	firstCount := 0
	secondCount := 0
	for _, mark := range slice {
		if mark == FIRST {
			firstCount++
		}
		if mark == SECOND {
			secondCount++
		}
	}
	if firstCount == len(slice) {
		return FIRST
	}
	if secondCount == len(slice) {
		return SECOND
	}
	return NONE
}

func reverseDiagonal(board [][]Player) (result []Player) {
	result = []Player{}
	for i := 0; i < len(board); i++ {
		result = append(result, board[i][len(board)-i-1])
	}
	return
}
