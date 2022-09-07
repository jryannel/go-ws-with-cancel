package ws

import (
	"context"
	"fmt"
	"time"
	"wswc/log"

	"github.com/gorilla/websocket"
)

var nextConnId = MakeIdGenerator("conn")

const (
	reconnectInterval = 1 * time.Second
	// max message size in bytes
	maxMessageSize = 512
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Connection represents a cancelable websocket connection.
type Connection struct {
	id     string
	socket *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

// Dial creates a new client connection.
// It tries repeatedly to connect to the server until successful or context is done.
func Dial(ctx context.Context, url string) (*Connection, error) {
	log.Debugf("dial: %s", url)
	socket, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err == nil {
		conn := NewConnection(ctx, socket)
		return conn, nil
	}
	ticker := time.NewTicker(reconnectInterval)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			socket, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
			if err == nil {
				log.Debugf("connected to: %s\n", url)
				return NewConnection(ctx, socket), nil
			} else {
				log.Debugf("dial: %s", err)
			}
		}
	}
	return nil, fmt.Errorf("dial error: %s", url)
}

// NewConnection handles a new websocket
func NewConnection(ctx context.Context, socket *websocket.Conn) *Connection {
	c := &Connection{
		id:     nextConnId(),
		socket: socket,
	}
	c.ctx, c.cancel = context.WithCancel(ctx)
	go c.run()
	return c
}

// Id returns the connection id.
func (c *Connection) Id() string {
	return c.id
}

// run setups the ping/pong handler and starts the read loop.
func (c *Connection) run() {
	c.socket.SetReadLimit(maxMessageSize)
	c.socket.SetPongHandler(func(string) error {
		deadline := time.Now().Add(pongWait)
		log.Debugf("conn: handle pong %v\n", deadline)
		return c.socket.SetReadDeadline(deadline)
	})
	c.socket.SetCloseHandler(func(code int, text string) error {
		// close connection and let write pump handle it
		c.closeSocket(fmt.Errorf("close %d %s", code, text))
		return nil
	})
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			log.Debugf("conn: write ping\n")
			err := c.socket.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				c.closeSocket(err)
				return
			}
		}
	}
}

// Write writes a message to the socket.
func (c *Connection) Write(data []byte) error {
	err := c.socket.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		c.closeSocket(err)
	}
	return err
}

// WriteJSON writes a JSON object to the socket.
func (c *Connection) WriteJSON(v interface{}) error {
	err := c.socket.WriteJSON(v)
	if err != nil {
		c.closeSocket(err)
	}
	return err
}

// Read blocks and reads a message from the socket.
func (c *Connection) Read() ([]byte, error) {
	for {
		select {
		case <-c.ctx.Done():
			return nil, c.ctx.Err()
		default:
			_, message, err := c.socket.ReadMessage()
			if err != nil {
				log.Debugf("conn: read error %s\n", err)
				c.closeSocket(err)
				return nil, err
			}
			log.Debugf("conn: read %s\n", message)
			return message, nil
		}
	}
}

// ReadJSON blocks and reads a JSON message from the socket.
func (c *Connection) ReadJSON(v interface{}) error {
	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		default:
			err := c.socket.ReadJSON(v)
			if err != nil {
				log.Debugf("conn: read error %s\n", err)
				c.closeSocket(err)
				return err
			}
			log.Debugf("conn: read %s\n", v)
			return nil
		}
	}
}

// Done returns a channel that is closed when the connection is closed.
func (c *Connection) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Stop closes the connection and the underlying socket.
func (c *Connection) Stop() {
	log.Debugf("conn: stop\n")
	c.cancel()
	c.closeSocket(fmt.Errorf("stopped"))
}

// closeSocket closes the underlying socket and calls the OnClosing callback.
func (c *Connection) closeSocket(err error) {
	log.Debugf("conn: close %s\n", err)
	if c.socket != nil {
		// write close message to socket
		c.socket.WriteMessage(websocket.CloseMessage, []byte{})
		// close socket
		c.socket.Close()
		c.socket = nil
	}
}
