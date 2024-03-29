package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Game struct {
	ID          int    `json:"id"`
	Fen         string `json:"fen"`
	Uci         string `json:"uci"`
	TimedUci    string `json:"timed_uci"`
	TimeControl string `json:"time_control"`
	White       string `json:"whitePlayer"`
	Black       string `json:"blackPlayer"`
	Winner      string `json:"winner"`
	Reason      string `json:"reason"`
	GameCreated string `json:"game_created"`
	GameEnded   string `json:"game_ended"`
}

// fetchLiveGames only needs id because node websocket will fetch the rest

func fetchLiveGames(amount int) []int {
	db, err := sql.Open("sqlite3", dbLoc)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	games, err := db.Query("SELECT id FROM game WHERE game_ended IS NULL ORDER BY id DESC LIMIT ?", amount)
	if err != nil {
		panic(err)
	}
	var ids []int
	for games.Next() {
		var id int
		err = games.Scan(&id)
		if err != nil {
			panic(err)
		}
		ids = append(ids, id)
	}
	return ids
}

func fetchFinishedGames(amount int, user string) []Game {
	db, err := sql.Open("sqlite3", dbLoc)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var games *sql.Rows
	if user == "" {
		games, err = db.Query("SELECT id, fen, IFNULL(uci, ''), IFNULL(timed_uci, ''), IFNULL(whitePlayerID, '0'), IFNULL(blackPlayerID, '0'), game_created, time_control FROM game WHERE game_ended IS NOT NULL AND id IN (SELECT id FROM game_ended WHERE reason != 'Game timed out') ORDER BY id DESC LIMIT ?", amount)
	} else {
		games, err = db.Query("SELECT id, fen, IFNULL(uci, ''), IFNULL(timed_uci, ''), IFNULL(whitePlayerID, '0'), IFNULL(blackPlayerID, '0'), game_created, time_control FROM game WHERE game_ended IS NOT NULL AND id IN (SELECT id FROM game_ended WHERE reason != 'Game timed out') AND id IN (SELECT id FROM game WHERE blackPlayerID IN (SELECT id FROM user WHERE username = ?) OR whitePlayerID IN (SELECT id FROM user WHERE username = ?)) ORDER BY id DESC LIMIT ?", user, user, amount)
	}
	if err != nil {
		panic(err)
	}
	var Games []Game
	for games.Next() {
		var game Game
		var whitePlayerID, blackPlayerID int
		err = games.Scan(&game.ID, &game.Fen, &game.Uci, &game.TimedUci, &whitePlayerID, &blackPlayerID, &game.GameCreated, &game.TimeControl)
		game.Winner, game.Reason, game.GameEnded = fetchResult(game.ID)
		game.White = idToUsername(whitePlayerID)
		game.Black = idToUsername(blackPlayerID)
		if err != nil {
			panic(err)
		}
		Games = append(Games, game)
	}
	return Games
}

func idToUsername(id int) string {
	if id == 0 {
		return ""
	}
	db, err := sql.Open("sqlite3", dbLoc)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var username string
	err = db.QueryRow("SELECT username FROM user WHERE id=?", id).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "Deleted User"
		}
		panic(err)
	}
	return username
}

func fetchResult(id int) (string, string, string) {
	db, err := sql.Open("sqlite3", dbLoc)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var winner, reason, date string
	err = db.QueryRow("SELECT winner, reason, game_ended FROM game_ended WHERE id=?", id).Scan(&winner, &reason, &date)
	if err != nil {
		panic(err)
	}
	return winner, reason, date
}

func fetchSingleGame(id int) (bool, Game, error) {
	db, err := sql.Open("sqlite3", dbLoc)
	var game Game
	if err != nil {
		return false, game, err
	}
	defer db.Close()

	var whitePlayerID, blackPlayerID, game_ended int
	var live = true
	err = db.QueryRow("SELECT id, fen, IFNULL(uci, ''), IFNULL(timed_uci, ''), IFNULL(whitePlayerID, '0'), IFNULL(blackPlayerID, '0'), IFNULL(game_ended, '0'), IFNULL(time_control, '0') FROM game WHERE id=?", id).Scan(&game.ID, &game.Fen, &game.Uci, &game.TimedUci, &whitePlayerID, &blackPlayerID, &game_ended, &game.TimeControl)
	if err != nil {
		return live, game, err
	}
	if game_ended != 0 {
		game.Winner, game.Reason, game.GameEnded = fetchResult(game.ID)
		live = false
	}
	game.White = idToUsername(whitePlayerID)
	game.Black = idToUsername(blackPlayerID)
	return live, game, nil
}
