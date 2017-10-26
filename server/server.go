package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"

	"../card"
	"../gamelist"
	"../user"
	"./socket"
)

// StartHTTP begins the socket server
func StartHTTP(db *sql.DB) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	sh := socket.CreateHandler()
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	games := gamelist.CreateGameList()

	server.On("connection", func(s socketio.Socket) {
		go initSocket(&s, db, &sh, &games)
	})
	http.Handle("/socket.io/", c.Handler(server))
	fmt.Println("Starting HTTP/Socket server...")
	http.ListenAndServe(":8000", nil)
}

func initSocket(so *socketio.Socket, db *sql.DB, sh *socket.Handler, games *gamelist.GameList) {
	cookie, e := (*so).Request().Cookie("connect.sid")
	if e != nil {
		return;
	}
	u, e := user.GetByCookie(cookie.Value, db)
	if e != nil {
		return;
	}
	fmt.Println("A user has connected")
	sh.Add(u.ID, so)
	(*so).On("disconnection", func() {
		fmt.Println("A user has disconnected")
		sh.Remove(so)
	})
	// Game Logic Events
	(*so).On("refreshgame", func() {
		sh.SendActionToUser(u.ID, socket.Action{Type: "game/SET_GAME_STATE", Payload: games.GetStateForUser(u)})
	})
	(*so).On("refreshlist", func() {
		// TODO - Change nil to actual game list
		sh.SendActionToUser(u.ID, socket.Action{Type: "games/SET_GAMES", Payload: nil})
	})
	(*so).On("join", func(gn string) {
		if games.JoinGame(u, gn) == nil {
		}
		sh.SendActionToUser(u.ID, socket.Action{Type: "game/SET_GAME_STATE", Payload: games.GetStateForUser(u)})
	})
	(*so).On("leave", func() {
		games.LeaveGame(u)
		sh.SendActionToUser(u.ID, socket.Action{Type: "game/RESET_GAME_STATE", Payload: nil})
		sh.SendActionToUser(u.ID, socket.Action{Type: "game/SET_GAME_STATE", Payload: games.GetStateForUser(u)})
	})
	(*so).On("kickplayer", func() {
		games.LeaveGame(u)
		sh.SendActionToUser(u.ID, socket.Action{Type: "game/RESET_GAME_STATE", Payload: nil})
		sh.SendActionToUser(u.ID, socket.Action{Type: "game/SET_GAME_STATE", Payload: games.GetStateForUser(u)})
	})
	(*so).On("playcard", func(c card.WhiteCard) {
		games.PlayCard(u, c)
	})
	(*so).On("voteplayer", func() {
	})
}
