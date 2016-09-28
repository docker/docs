package adminserver

import (
	"bytes"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type webSocketClient struct {
	id      string
	user    *authn.User
	ws      *websocket.Conn
	manager *webSocketManager
	ch      chan interface{}
	readCh  chan interface{}
	doneCh  chan bool
}

type filter func(*webSocketClient, interface{}) bool

type pubStruct struct {
	data   interface{}
	filter filter
}

func (c *webSocketClient) Write(data interface{}) {
	select {
	case c.ch <- data:
	default:
		c.manager.Del(c)
	}
}

func (c *webSocketClient) Done() {
	c.doneCh <- true
}

func (c *webSocketClient) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *webSocketClient) listenWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.manager.Del(c)
		c.ws.Close()
	}()
	for {
		select {
		case message := <-c.ch:
			if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.manager.Err(err)
				return
			}
			if err := c.ws.WriteJSON(message); err != nil {
				c.manager.Err(err)
				return
			}
		case <-ticker.C:
			if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.manager.Err(err)
				return
			}
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.manager.Err(err)
				return
			}
		case <-c.doneCh:
			return
		}
	}
}

func (c *webSocketClient) listenRead() {
	defer func() {
		c.manager.Del(c)
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				c.manager.Err(err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.readCh <- message
	}
}

type webSocketManager struct {
	clients map[string]*webSocketClient
	addCh   chan *webSocketClient
	delCh   chan *webSocketClient
	pubCh   chan pubStruct
	errCh   chan error
}

func NewWebSocketManager() *webSocketManager {
	return &webSocketManager{
		clients: make(map[string]*webSocketClient),
		addCh:   make(chan *webSocketClient),
		delCh:   make(chan *webSocketClient),
		pubCh:   make(chan pubStruct),
		errCh:   make(chan error),
	}
}

// NewWebSocketClient creates a new websocketclient, registers the client, and listens
func (m *webSocketManager) NewWebSocketClient(ws *websocket.Conn, manager *webSocketManager, user *authn.User) {
	ch := make(chan interface{}, maxMessageSize)
	readCh := make(chan interface{}, maxMessageSize)
	doneCh := make(chan bool)

	client := &webSocketClient{
		id:      uuid.NewV4().String(),
		user:    user,
		ws:      ws,
		manager: manager,
		ch:      ch,
		readCh:  readCh,
		doneCh:  doneCh,
	}
	m.Add(client)
	client.Listen()
}

func (m *webSocketManager) Add(c *webSocketClient) {
	m.addCh <- c
}

func (m *webSocketManager) Del(c *webSocketClient) {
	c.Done()
	m.delCh <- c
}

func (m *webSocketManager) Publish(data interface{}) {
	m.FilteredPublish(data, nil)
}

func (m *webSocketManager) FilteredPublish(data interface{}, f filter) {
	m.pubCh <- pubStruct{
		data:   data,
		filter: f,
	}
}

func (m *webSocketManager) Err(err error) {
	m.errCh <- err
}

func (m *webSocketManager) sendAll(pub pubStruct) {
	for _, c := range m.clients {
		if pub.filter == nil || pub.filter(c, pub.data) {
			c.Write(pub.data)
		}
	}
}

func (m *webSocketManager) Listen() {
	for {
		select {
		case c := <-m.addCh:
			m.clients[c.id] = c
		case c := <-m.delCh:
			delete(m.clients, c.id)
		case data := <-m.pubCh:
			m.sendAll(data)
		case err := <-m.errCh:
			log.Errorf("WebSocket error: %+v", err)
		}
	}
}
