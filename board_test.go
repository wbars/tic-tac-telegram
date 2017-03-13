package tic_tac_bot

import "testing"

func TestBoardCreating(t *testing.T) {
	game := createGame(4)
	if game.Size != 4 {
		t.Error("Expected 4, got", game.Size)
	}
	for i := 0; i < game.Size; i++ {
		for j := 0; j < game.Size; j++ {
			if game.Board[i][j] != NONE {
				t.Error("Expected default board with all NONE's (0) states, got", game.Board[i][j], "on", i, j)
			}
		}
	}

	winner := game.getWinner()
	if winner != NONE {
		t.Error("Expected no winner in default game, got", winner)
	}

	repr := "....\n" +
		"....\n" +
		"....\n" +
		"...."
	expectedRepr(game, repr, t)
}

func TestBoardTurns(t *testing.T) {
	game := createGame(4)
	turn := game.makeNextTurn(2, 3)

	if turn != nil {
		t.Error("Expected turn without errors on empty cells, got", turn.Error())
	}
	if game.Board[2][3] != FIRST {
		t.Error("Expected first being set after turn, got", game.Board[2][3])
	}

	if game.makeNextTurn(2, 3) == nil {
		t.Error("Expected error when making turn on taken cell, got clean turn")
	}

	game.makeNextTurn(1, 1)
	game.makeNextTurn(3, 3)
	game.makeNextTurn(2, 2)

	repr := "....\n" +
		".O..\n" +
		"..OX\n" +
		"...X"
	expectedRepr(game, repr, t)
}

func TestGameRepr(t *testing.T) {
	repr := "....\n" +
		".O..\n" +
		"..OX\n" +
		"...X"
	game := fromRepr(repr)
	if repr != game.repr() {
		t.Error("Repr was muted after unpacking, expected", repr, "got", game.repr())
	}
	for i := 0; i < game.Size; i++ {
		for j := 0; j < game.Size; j++ {
			cell := game.Board[i][j]
			if i == 1 && j == 1 || i == 2 && j == 2 {
				if cell != SECOND {
					t.Error("Expected second player on this position:", i, j, "got", cell)
				}
			} else if i == 2 && j == 3 || i == 3 && j == 3 {
				if cell != FIRST {
					t.Error("Expected first player on this position:", i, j, "got", cell)
				}
			} else if cell != NONE {
				t.Error("Expected None on this position:", i, j, "got", cell)
			}
		}
	}
}

func testWin(repr string, winner Player, t *testing.T) {
	game := fromRepr(repr)
	if game.getWinner() != winner {
		t.Error(winner, " player win board", repr, "\nHas winner: ", game)
	}
}

func TestGameWinFirst(t *testing.T) {
	repr := "X...\n" +
		"X.O.\n" +
		"X...\n" +
		"X..O"
	testWin(repr, FIRST, t)
}

func TestGameWinSecond(t *testing.T) {
	repr := "O.X.\n" +
		"O.O.\n" +
		"OX..\n" +
		"O..O"
	testWin(repr, SECOND, t)
}

func TestGameWinFalsePositive(t *testing.T) {
	repr := "O...\n" +
		".O..\n" +
		"..OX\n" +
		"...X"
	testWin(repr, NONE, t)
}

func TestGameWinMainDiagonal(t *testing.T) {
	repr := "O.X.\n" +
		"XO..\n" +
		"OXO.\n" +
		"X..O"
	testWin(repr, SECOND, t)
}

func TestGameWinReverseDiagonal(t *testing.T) {
	repr := "X.XX\n" +
		"XOX.\n" +
		"OXO.\n" +
		"X..O"
	testWin(repr, FIRST, t)
}

func expectedRepr(game *Game, repr string, t *testing.T) {
	if game.repr() != repr {
		t.Error("Expected this repr:\n", repr, "got\n", game.repr())
	}
}
