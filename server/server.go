package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"
)

// StartHTTP begins the socket server
func StartHTTP() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		fmt.Println("A user has connected")
		so.On("disconnection", func() {
			fmt.Println("A user has disconnected")
		})
	})
	http.Handle("/socket.io/", c.Handler(server))
	fmt.Println("Starting HTTP/Socket server...")
	http.ListenAndServe(":8000", nil)
}
