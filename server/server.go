package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"

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
	server.On("connection", func(so socketio.Socket) {
		fmt.Println("A user has connected")
		cookie, _ := so.Request().Cookie("connect.sid")
		uid := user.GetIDByCookie(cookie.Value, db)
		sh.Add(uid, &so)
		so.On("disconnection", func() {
			fmt.Println("A user has disconnected")
			sh.Remove(&so)
		})
	})
	http.Handle("/socket.io/", c.Handler(server))
	fmt.Println("Starting HTTP/Socket server...")
	http.ListenAndServe(":8000", nil)
}
