package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	games := gamelist.CreateGameList(sh)

	socketIOMux, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	socketIOMux.On("connection", func(s socketio.Socket) {
		go initSocket(&s, db, sh, &games)
	})
	http.Handle("/socket.io/", c.Handler(socketIOMux))
	http.Handle("/game/", c.Handler(createGameMux("/game", db, sh, &games)))
	http.Handle("/gamelist", c.Handler(createGameListMux("/gamelist", db, sh, &games)))
	fmt.Println("Starting HTTP/Socket server...")
	http.ListenAndServe(":8000", nil)
}

func initSocket(so *socketio.Socket, db *sql.DB, sh *socket.Handler, games *gamelist.GameList) {
	cookie, err := (*so).Request().Cookie("connect.sid")
	if err != nil {
		return
	}
	u, err := user.GetByCookie(cookie.Value, db)
	if err != nil {
		return
	}
	fmt.Println("A user has connected")
	sh.Add(u.ID, so)
	(*so).On("disconnection", func() {
		fmt.Println("A user has disconnected")
		sh.Remove(so)
	})
	sh.SendActionToUser(u.ID, socket.Action{Type: "game/SET_GAME_STATE", Payload: games.GetStateForUser(u)})
}

// GameCreateMessage JSON structure for HTTP requests to the game creation endpoint
type GameCreateMessage struct {
	Name        string `json:"name"`
	CardpackIDs []int  `json:"cardpackIDs"`
	MaxPlayers  int    `json:"maxPlayers"`
}

func createGameMux(path string, db *sql.DB, sh *socket.Handler, gl *gamelist.GameList) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(path+"/state", func(w http.ResponseWriter, r *http.Request) {
		u, err := user.GetByRequest(r, db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(gl.GetStateForUser(u))
	})
	mux.HandleFunc(path+"/create", func(w http.ResponseWriter, r *http.Request) {
		u, err := user.GetByRequest(r, db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var msg GameCreateMessage
		err = json.Unmarshal(b, &msg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		bc, wc := card.GetCards(msg.CardpackIDs, db)
		err = gl.CreateGame(u, msg.Name, msg.MaxPlayers, bc, wc)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(gl.GetStateForUser(u))
	})
	mux.HandleFunc(path+"/start", func(w http.ResponseWriter, r *http.Request) {
		u, err := user.GetByRequest(r, db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = gl.StartGame(u.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(gl.GetStateForUser(u))
	})
	mux.HandleFunc(path+"/stop", func(w http.ResponseWriter, r *http.Request) {
		u, err := user.GetByRequest(r, db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = gl.StopGame(u.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(gl.GetStateForUser(u))
	})
	mux.HandleFunc(path+"/join", func(w http.ResponseWriter, r *http.Request) {
		u, err := user.GetByRequest(r, db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		gl.JoinGame(u, string(b))
		json.NewEncoder(w).Encode(gl.GetStateForUser(u))
	})
	mux.HandleFunc(path+"/leave", func(w http.ResponseWriter, r *http.Request) {
		u, err := user.GetByRequest(r, db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		gl.LeaveGame(u)
		json.NewEncoder(w).Encode(gl.GetStateForUser(u))
	})
	mux.HandleFunc(path+"/card", func(w http.ResponseWriter, r *http.Request) {
		u, err := user.GetByRequest(r, db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var msg int
		err = json.Unmarshal(b, &msg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		gl.PlayCard(u, msg)
		json.NewEncoder(w).Encode(true)
	})
	mux.HandleFunc(path+"/kickplayer", func(w http.ResponseWriter, r *http.Request) {
		u, err := user.GetByRequest(r, db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var msg int
		err = json.Unmarshal(b, &msg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		gl.KickUser(u, msg)
		json.NewEncoder(w).Encode(true)
	})
	mux.HandleFunc(path+"/vote", func(w http.ResponseWriter, r *http.Request) {
		u, err := user.GetByRequest(r, db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var msg int
		err = json.Unmarshal(b, &msg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		gl.VoteCard(u, msg)
		json.NewEncoder(w).Encode(true)
	})
	return mux
}

func createGameListMux(path string, db *sql.DB, sh *socket.Handler, gl *gamelist.GameList) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(gl.GetList())
	})
	return mux
}
