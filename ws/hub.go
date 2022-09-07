package ws

import (
	"context"
	"encoding/json"
	"net/http"

	"wswc/log"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Hub struct {
	id          string
	register    chan *Connection
	unregister  chan *Connection
	broadcast   chan []byte
	connections map[*Connection]bool
	ctx         context.Context
	requests    chan *Request
}

var nextHubId = MakeIdGenerator("hub")

func NewHub(ctx context.Context) *Hub {
	return &Hub{
		id:          nextHubId(),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		broadcast:   make(chan []byte),
		connections: make(map[*Connection]bool),
		requests:    make(chan *Request),
		ctx:         ctx,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case <-h.ctx.Done():
			return
		case conn := <-h.register:
			log.Infof("register conn: %s\n", conn.Id())
			h.connections[conn] = true
		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				conn.Stop()
			}
		case data := <-h.broadcast:
			for conn := range h.connections {
				conn.Write(data)
			}
		}
	}
}

func (h *Hub) Requests() <-chan *Request {
	return h.requests
}

// Broadcast sends data to all connections
func (h *Hub) Broadcast(data []byte) {
	h.broadcast <- data
}

// BroadcastJSON encodes v as JSON and broadcasts it to all connections.
func (h *Hub) BroadcastJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	h.Broadcast(data)
	return nil
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Infoln("upgrade:", err)
		return
	}
	log.Infof("new connection: %s\n", socket.RemoteAddr())
	conn := NewConnection(h.ctx, socket)
	go func() {
		// unregister connection when it is done
		<-conn.Done()
		log.Infof("unregister conn: %s\n", conn.Id())
		h.unregister <- conn
	}()
	go func() {
		// read requests from connection and stream them to the hub
		defer func() {
			h.unregister <- conn
		}()
		for {
			select {
			case <-h.ctx.Done():
				return
			default:
				data, err := conn.Read()
				log.Debugf("read: %s\n", data)
				if err != nil {
					log.Debugf("read error: %s\n", err)
					return
				}
				req := &Request{
					data:       data,
					connection: conn,
					err:        err,
				}
				h.requests <- req
			}
		}
	}()
	h.register <- conn
}
