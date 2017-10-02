package socket

import "github.com/googollee/go-socket.io"

import "fmt"

// Handler manages user sockets
type Handler struct {
	uToS map[int][]socketio.Socket
	sToU map[socketio.Socket]int
}

// CreateHandler generates a socket handler
func CreateHandler() Handler {
	return Handler{uToS: make(map[int][]socketio.Socket), sToU: make(map[socketio.Socket]int)}
}

// Add registers reference to a socket
func (h Handler) Add(userID int, s *socketio.Socket) {
	if _, ok := h.uToS[userID]; ok {
		h.uToS[userID] = append(h.uToS[userID], *s)
	} else {
		h.uToS[userID] = []socketio.Socket{*s}
	}
	h.sToU[*s] = userID
}

// Remove deletes reference to a socket
func (h Handler) Remove(s *socketio.Socket) {
	if userID, ok := h.sToU[*s]; ok {
		for i, soc := range h.uToS[userID] {
			if *s == soc {
				h.uToS[userID][i] = h.uToS[userID][len(h.uToS)-1]
				h.uToS[userID] = h.uToS[userID][:len(h.uToS)-1]
				break
			}
		}
		if len(h.uToS[userID]) == 0 {
			delete(h.uToS, userID)
		}
		delete(h.sToU, *s)
	} else {
		fmt.Println("Attempted to delete a socket that does not exist")
	}
}

// SendDataToUser sends data to all sockets belonging to a particular user
func (h Handler) SendDataToUser(userID int, dataType string, data string) {
	if _, ok := h.uToS[userID]; ok {
		for _, s := range h.uToS[userID] {
			s.Emit(dataType, data)
		}
	}
}

// SendDataToUsers sends data to all sockets belonging to a list of users
func (h Handler) SendDataToUsers(userIDs []int, dataType string, data string) {
	for _, id := range userIDs {
		h.SendDataToUser(id, dataType, data)
	}
}

// SendDataToAllUsers sends data to all sockets
func (h Handler) SendDataToAllUsers(dataType string, data string) {
	for s := range h.sToU {
		s.Emit(dataType, data)
	}
}
