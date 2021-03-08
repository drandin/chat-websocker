package app

import (
	"encoding/json"
	"fmt"
	"time"
)

type message struct {
	payload []byte
	roomId string
	userId int
}

type MessagePublic struct {
	Payload string `json:"payload"`
	RoomId string `json:"roomId"`
	UserId int `json:"userId"`
}

type subscription struct {
	conn *connection
	room string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {

	// Registered connections.
	rooms map[string]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var h = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

// Отправка сообщения в комнату "room-1"
func (h *hub) test()  {
	duration := time.Second * 7
	time.Sleep(duration)
	fmt.Println("pause")
	fmt.Println(h.rooms["room-1"])
	m := message{[]byte("Сообщение от сервера."), "room-1", 0}
	h.broadcast <- m
}

func (h *hub) run() {

	go h.test()

	for {

		select {
		case s := <-h.register:

			connections := h.rooms[s.room]

			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
			}

			h.rooms[s.room][s.conn] = true

		case s := <-h.unregister:

			connections := h.rooms[s.room]

			if connections != nil {

				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}

		case m := <-h.broadcast:

			connections := h.rooms[m.roomId]

			for c := range connections {

				message := MessagePublic{
					Payload: string(m.payload),
					RoomId: m.roomId,
					UserId: m.userId,
				}

				messageJson, _ := json.Marshal(message)

				fmt.Println(string(messageJson))

				select {

				case c.send <- messageJson:

				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.roomId)
					}
				}

			}
		}
	}
}