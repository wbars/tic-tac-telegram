package tic_tac_bot

import (
	"strings"
	"errors"
	"fmt"
)

type Game struct {
	Board      [][]Player `json:"board"`
	Size       int `json:"size"`
	LastPlayer Player `json:"last_player"`
}

type Player int

const (
	NONE   Player = iota
	FIRST  Player = iota
	SECOND Player = iota
)

func createGame(size int) *Game {
	table := [][]Player{}
	for i := 0; i < size; i++ {
		table = append(table, make([]Player, size))
	}
	return &Game{Board: table, Size: size}
}

func (game Game) repr() string {
	rows := []string{}
	for i := 0; i < game.Size; i++ {
		row := ""
		for j := 0; j < game.Size; j++ {
			row += getMark(game.Board[i][j])
		}
		rows = append(rows, row)

	}
	return strings.Join(rows, "\n")
}

func fromRepr(repr string) *Game {
	table := [][]Player{}
	reprTable := strings.Split(repr, "\n")
	size := len(reprTable)
	for i := 0; i < size; i++ {
		row := []Player{}
		for j := 0; j < size; j++ {
			row = append(row, fromMark(string(reprTable[i][j])))
		}
		table = append(table, row)
	}
	return &Game{Board: table, Size: size}
}

func (game *Game) makeNextTurn(row int, col int) (error) {
	if game.Board[row][col] != NONE {
		return errors.New("Cell is taken")
	}

	game.Board[row][col] = game.nextPlayer()
	return nil
}
func (game *Game) nextPlayer() Player {
	fmt.Print(game.LastPlayer)
	if game.LastPlayer != FIRST {
		game.LastPlayer = FIRST
		return FIRST
	}
	game.LastPlayer = SECOND
	return SECOND

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

func (game Game) getWinner() Player {
	for row := 0; row < game.Size; row++ {
		winner := getSliceWinner(game.Board[row])
		if winner != NONE {
			return winner
		}
	}

	for col := 0; col < game.Size; col++ {
		winner := getSliceWinner(getVerticalSlice(game.Board, col))
		if winner != NONE {
			return winner
		}
	}

	winner := getSliceWinner(mainDiagonal(game.Board))
	if winner != NONE {
		return winner
	}

	winner = getSliceWinner(reverseDiagonal(game.Board))
	if winner != NONE {
		return winner
	}
	return NONE

}
func mainDiagonal(game [][]Player) (result []Player) {
	result = []Player{}
	for i := 0; i < len(game); i++ {
		result = append(result, game[i][i])
	}
	return
}
func getVerticalSlice(game [][]Player, col int) (result []Player) {
	result = []Player{}
	for row := 0; row < len(game); row++ {
		result = append(result, game[row][col])
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

func reverseDiagonal(game [][]Player) (result []Player) {
	result = []Player{}
	for i := 0; i < len(game); i++ {
		result = append(result, game[i][len(game)-i-1])
	}
	return
}
