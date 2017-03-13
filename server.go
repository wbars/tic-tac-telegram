package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"encoding/json"
	"time"
	"math/rand"
)

type Response struct {
	Success bool `json:"success"`
	Data    interface{} `json:"data"`
	Message string `json:"message"`
}

func success(data interface{}) Response {
	return Response{Success: true, Data: data, Message: ""}
}

func fail(message string) Response {
	return Response{Success: false, Data: nil, Message: message}
}

const DEFAULT_GAME_SIZE int = 5

func createRouter() *mux.Router {
	router := mux.NewRouter()
	games := map[int]*Game{}
	tokens := map[int]string{}
	userPlayer := map[int]Player{}

	router.HandleFunc("/game", createNewGameHandler(games, tokens, userPlayer)).Methods("POST")
	router.HandleFunc("/game/{id}", createGetGameHandler(games)).Methods("GET")
	router.HandleFunc("/player/{id}", createUserPlayerHandler(userPlayer)).Methods("GET")
	router.HandleFunc("/game/{id}/win", createGetGameWin(games)).Methods("GET")
	router.HandleFunc("/game/{id}/turn", NextTurn(games, tokens, userPlayer)).Methods("POST")

	return router
}
func createUserPlayerHandler(players map[int]Player) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(req)["id"])
		if player, ok := players[id]; ok {
			writeResponse(w, success(player))
			return
		}
		writeResponse(w, fail("Cant find game"))
	}
}

type EndGame struct {
	Game   Game `json:"game"`
	Winner Player `json:"winner"`
}

func NextTurn(games map[int]*Game, tokens map[int]string, usersPlayer map[int]Player) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(req)["id"])
		if game, ok := games[id]; ok {
			req.ParseForm()
			if game.isFinished() {
				writeResponse(w, success(CreateEndGame(game)))
				return
			}

			if !isAuthorized(req, tokens, id) {
				writeResponse(w, fail("Invallid token"))
				return
			}

			if game.LastPlayer == usersPlayer[id] {
				writeResponse(w, fail("Not your turn"))
				return
			}

			i, err1 := strconv.Atoi(req.Form.Get("i"))
			j, err2 := strconv.Atoi(req.Form.Get("j"))
			if err1 != nil || err2 != nil {
				writeResponse(w, fail(err1.Error() + err2.Error()))
			}
			err := game.makeNextTurn(i, j)
			if err != nil {
				writeResponse(w, fail(err.Error()))
				return
			}
			if game.isFinished() {
				writeResponse(w, success(CreateEndGame(game)))
				return
			}

			makeBotTurn(games, id, usersPlayer)
			if game.isFinished() {
				writeResponse(w, success(CreateEndGame(game)))
				return
			}

			writeResponse(w, success(games[id]))
			return
		}
		writeResponse(w, fail("Cant find game"))
	}
}
func writeResponse(w http.ResponseWriter, fail Response) error {
	return json.NewEncoder(w).Encode(fail)
}
func CreateEndGame(game *Game) EndGame {
	return EndGame{Game: *game, Winner: game.getWinner()}
}
func (game Game) isFinished() bool {
	if game.getWinner() != NONE {
		return true
	}
	for i := 0; i < len(game.Board); i++ {
		for j := 0; j < len(game.Board); j++ {
			if game.Board[i][j] == NONE {
				return false
			}
		}
	}
	return true
}
func (player Player) getAlias() string {
	if player == FIRST {
		return "FIRST"
	}
	if player == SECOND {
		return "SECOND"
	}
	return "NONE"
}
func createGetGameWin(games map[int]*Game) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(req)["id"])
		if game, ok := games[id]; ok {
			writeResponse(w, success(game.getWinner()))
			return
		}
		writeResponse(w, fail("Cant find game"))
	}
}
func isAuthorized(req *http.Request, tokens map[int]string, id int) bool {
	return req.Form.Get("token") == tokens[id]
}

func createGetGameHandler(games map[int]*Game) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(req)["id"])
		if game, ok := games[id]; ok {
			writeResponse(w, success(game))
			return
		}
		writeResponse(w, fail("Cant find game"))
	}
}
func createNewGameHandler(games map[int]*Game, tokens map[int]string, usersPlayer map[int]Player) func(w http.ResponseWriter, req *http.Request) {
	gameCounter := 0
	return func(w http.ResponseWriter, req *http.Request) {
		games[gameCounter] = createGame(DEFAULT_GAME_SIZE)
		tokens[gameCounter] = RandString(15)
		usersPlayer[gameCounter] = getRandomPlayer()

		if usersPlayer[gameCounter] != FIRST {
			makeBotTurn(games, gameCounter, usersPlayer)
		}

		type GameCreated struct {
			Id         int `json:"id"`
			Token      string `json:"token"`
			Game       Game `json:"game"`
			UserPlayer Player `json:"user_player"`
		}
		gameCreated := GameCreated{Id: gameCounter, Token: tokens[gameCounter], Game: *games[gameCounter], UserPlayer: usersPlayer[gameCounter]}
		gameCounter++
		writeResponse(w, success(gameCreated))
	}
}
func makeBotTurn(games map[int]*Game, gameId int, usersPlayer map[int]Player) {
	games[gameId].makeNextTurn(getNaiveMove(*games[gameId], getOppositePlayer(usersPlayer[gameId])))
}
func getOppositePlayer(player Player) Player {
	if player == FIRST {
		return SECOND
	}
	return FIRST
}

func getRandomPlayer() Player {
	if rand.Intn(2) == 0 {
		return FIRST
	}
	return SECOND
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
