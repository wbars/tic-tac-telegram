package main

func getNaiveMove(game Game, player Player) (int, int) {
	board := game.Board
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			if board[i][j] == NONE {
				return i, j
			}
		}
	}
	return -1, -1
}
