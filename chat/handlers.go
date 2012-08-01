package chat

import (
	"code.google.com/p/go.net/websocket"
	. "github.com/paulbellamy/mango"
	"github.com/sunfmin/mangotemplate"
)

type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan string

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan string),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan string
}

func (c *connection) reader() {
	for {
		var message string
		err := websocket.Message.Receive(c.ws, &message)
		if err != nil {
			break
		}
		h.broadcast <- message
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := websocket.Message.Send(c.ws, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func ReceiveChatMessage(ws *websocket.Conn) {
	go h.run()

	c := &connection{send: make(chan string, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}

type IndexData struct {
	WebSocketHost string
}

func Index(env Env) (status Status, headers Headers, body Body) {
	r := env.Request()
	mangotemplate.ForRender(env, "chat/index", &IndexData{WebSocketHost: r.Host})
	headers = Headers{}
	return
}
