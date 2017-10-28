package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"

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
	games := gamelist.CreateGameList()

	socketIOMux, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	socketIOMux.On("connection", func(s socketio.Socket) {
		go initSocket(&s, db, &sh, &games)
	})
	http.Handle("/socket.io/", c.Handler(socketIOMux))
	http.Handle("/game/", c.Handler(createGameMux("/game", db, &sh, &games)))
	http.Handle("/gamelist/", c.Handler(createGameListMux("/gamelist", db, &sh, &games)))
	fmt.Println("Starting HTTP/Socket server...")
	http.ListenAndServe(":8000", nil)
}

func initSocket(so *socketio.Socket, db *sql.DB, sh *socket.Handler, games *gamelist.GameList) {
	cookie, e := (*so).Request().Cookie("connect.sid")
	if e != nil {
		return
	}
	u, e := user.GetByCookie(cookie.Value, db)
	if e != nil {
		return
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
	// TODO - Allow function to accept other input types without crashing the server
	(*so).On("playcard", func(cID int) {
		games.PlayCard(u, cID)
	})
	(*so).On("vote", func() {
	})
}

func createGameMux(path string, db *sql.DB, sh *socket.Handler, gl *gamelist.GameList) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(path+"/state", func(w http.ResponseWriter, r *http.Request) {
		cookie, e := r.Cookie("connect.sid")
		if e != nil {
			http.Error(w, "Not logged in", 500)
			return
		}
		u, e := user.GetByCookie(cookie.Value, db)
		if e != nil {
			http.Error(w, "Not logged in", 500)
			return
		}
		json.NewEncoder(w).Encode(gl.GetStateForUser(u))
	})
	return mux
}

func createGameListMux(path string, db *sql.DB, sh *socket.Handler, gl *gamelist.GameList) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(path+"/testt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Method Type: " + r.Method)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	mux.HandleFunc(path+"/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Method Type: " + r.Method)
		fmt.Fprintf(w, "Hello there, %s", r.Method)
	})
	return mux
}
