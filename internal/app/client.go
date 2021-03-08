package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader {
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {

	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	userId int

}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {

	c := s.conn

	defer func() {
		h.unregister <- s
		_ = c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	_ = c.ws.SetReadDeadline(time.Now().Add(pongWait))

	c.ws.SetPongHandler(func(string) error {
		_ = c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.ws.ReadMessage()

		if err != nil {

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}

			break
		}

		m := message{msg, s.room, c.userId}
		h.broadcast <- m
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	_ = c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {

	c := s.conn

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		_ = c.ws.Close()
	}()

	for {
		select {

		case message, ok := <-c.send:

			if !ok {
				_ = c.write(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:

			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, settings *settings) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &connection{
		ws: ws,
		send: make(chan []byte, 256),
		userId: settings.userId,
	}

	s := subscription{c, settings.roomId}
	h.register <- s
	fmt.Println(s.room, s.conn.userId)
	go s.writePump()
	go s.readPump()
}